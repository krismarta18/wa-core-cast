package billing

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"wacast/core/database"
	"wacast/core/models"
)

type Service struct {
	db *database.Database
}

var (
	ErrBillingPlanNotFound    = errors.New("billing plan not found")
	ErrBillingPlanInactive    = errors.New("billing plan is inactive")
	ErrNoActiveSubscription   = errors.New("no active subscription found")
	ErrDeviceLimitReached     = errors.New("device limit reached for your current plan")
	ErrMessageLimitReached    = errors.New("daily message limit reached for your current plan")
)

func NewService(db *database.Database) *Service {
	return &Service{db: db}
}

// CheckDeviceLimit verifies if a user can add/connect another device
func (s *Service) CheckDeviceLimit(ctx context.Context, userID uuid.UUID) error {
	// Personal version: No device limits
	return nil
}

// CheckMessageLimit verifies if a user can send more messages today
func (s *Service) CheckMessageLimit(ctx context.Context, userID uuid.UUID) error {
	// Personal version: No message limits
	return nil
}

func (s *Service) CheckoutDummy(ctx context.Context, userID string, planID string) (*models.BillingCheckoutResponse, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}

	pid, err := uuid.Parse(planID)
	if err != nil {
		return nil, fmt.Errorf("invalid plan id: %w", err)
	}

	plan, err := s.db.GetBillingPlanByID(pid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrBillingPlanNotFound
		}
		return nil, fmt.Errorf("get billing plan: %w", err)
	}
	if !plan.IsActive {
		return nil, ErrBillingPlanInactive
	}

	tx, err := s.db.GetConnection().BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	now := time.Now()
	endDate := calculateSubscriptionEnd(now, plan.BillingCycle)
	renewalDate := endDate

	if _, err = tx.ExecContext(ctx, s.db.Translate(`
		UPDATE subscriptions
		SET status = 'inactive', updated_at = $2
		WHERE user_id = $1 AND status = 'active'
	`), uid, now); err != nil {
		return nil, fmt.Errorf("deactivate active subscriptions: %w", err)
	}

	subscriptionID := uuid.New()
	if _, err = tx.ExecContext(ctx, s.db.Translate(`
		INSERT INTO subscriptions (
			id, user_id, plan_id, status, start_date, end_date, renewal_date, auto_renew, max_devices, max_messages_per_day, created_at, updated_at
		)
		VALUES ($1, $2, $3, 'active', $4, $5, $6, true, $7, $8, $4, $4)
	`), subscriptionID, uid, pid, now, endDate, renewalDate, plan.MaxDevices, plan.MaxMessagesPerDay); err != nil {
		return nil, fmt.Errorf("create subscription: %w", err)
	}

	InvoiceID := uuid.New()
	if _, err = tx.ExecContext(ctx, s.db.Translate(`INSERT INTO invoices(id, user_id, subscription_id,invoice_number, issue_date,due_date,
	paid_at, amount, currency, status, payment_method, created_at) VALUES 
	($1, $2,$3,$4,$5::date,$5::date,$5::timestamptz,$6,$7,$8,$9,$5::timestamptz)`), InvoiceID, uid, subscriptionID,
		fmt.Sprintf("INV-%s-%s", now.Format("2006-01"), subscriptionID.String()[:8]),
		now, plan.Price, "Rupiah", "paid", "dummy"); err != nil {
		return nil, fmt.Errorf("create invoice: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	response := &models.BillingCheckoutResponse{
		Subscription: &models.BillingCurrentPlanResponse{
			SubscriptionID: subscriptionID,
			PlanID:         plan.ID,
			Name:           plan.Name,
			Price:          plan.Price,
			BillingCycle:   normalizeBillingCycle(plan.BillingCycle),
			RenewalDate:    &renewalDate,
			QuotaUsed:      0,
			QuotaLimit:     plan.MaxMessagesPerDay,
			DeviceUsed:     0,
			DeviceMax:      plan.MaxDevices,
			AutoRenew:      true,
			Status:         models.SubscriptionStatusActive,
			Features:       plan.Features,
		},
		Invoice: models.BillingInvoiceSummary{
			ID:             fmt.Sprintf("INV-%s-%s", now.Format("2006-01"), subscriptionID.String()[:8]),
			SubscriptionID: subscriptionID,
			Date:           now,
			PlanName:       plan.Name,
			Amount:         plan.Price,
			Status:         models.SubscriptionStatusActive,
		},
		PaymentStatus: "paid",
		PaymentMethod: "dummy",
	}

	return response, nil
}

func (s *Service) GetOverview(ctx context.Context, userID string) (*models.BillingOverviewResponse, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}

	usageHistory, err := s.getUsageHistory(ctx, uid, 7)
	if err != nil {
		return nil, fmt.Errorf("get usage history: %w", err)
	}

	deviceUsed, err := s.db.CountUserDevices(uid)
	if err != nil {
		return nil, fmt.Errorf("count user devices: %w", err)
	}

	response := &models.BillingOverviewResponse{
		UsageHistory: usageHistory,
		Invoices:     []models.BillingInvoiceSummary{}, // Empty invoices for personal
	}

	// Fill with a mock Personal Lifetime plan
	response.CurrentPlan = &models.BillingCurrentPlanResponse{
		Name:         "Personal Lifetime",
		Price:        0,
		BillingCycle: "lifetime",
		QuotaUsed:    sumSentMessages(usageHistory),
		QuotaLimit:   0, // 0 = unlimited in UI usually
		DeviceUsed:   deviceUsed,
		DeviceMax:    0, // 0 = unlimited
		Status:       "active",
		Features:     json.RawMessage(`["Unlimited Devices", "Unlimited Messages", "Priority Support"]`),
	}

	response.Plans = []models.BillingPlanSummary{
		{
			Name:       "Personal Lifetime",
			Price:      0,
			QuotaLimit: 0,
			DeviceMax:  0,
			Current:    true,
			IsActive:   true,
		},
	}

	return response, nil
}

func (s *Service) getUsageHistory(ctx context.Context, userID uuid.UUID, days int) ([]models.BillingUsagePoint, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT DATE(created_at) AS usage_date,
		       COUNT(*) FILTER (WHERE direction = 'outbound' AND status IN ('sent', 'delivered', 'read')) AS sent_count,
		       COUNT(*) FILTER (WHERE direction = 'outbound' AND status = 'failed') AS failed_count
		FROM messages
		WHERE user_id = $1 AND created_at >= CURRENT_DATE - ($2::int - 1) * INTERVAL '1 day'
		GROUP BY DATE(created_at)
		ORDER BY usage_date ASC
	`, userID, days)
	if err != nil {
		if !isMissingColumnError(err, "status") && !isMissingColumnError(err, "user_id") {
			return nil, err
		}

		rows, err = s.db.QueryContext(ctx, `
			SELECT DATE(m.created_at) AS usage_date,
			       COUNT(*) FILTER (WHERE UPPER(m.direction) IN ('OUT', 'OUTBOUND') AND m.status_message IN (1, 2, 3)) AS sent_count,
			       COUNT(*) FILTER (WHERE UPPER(m.direction) IN ('OUT', 'OUTBOUND') AND m.status_message = 4) AS failed_count
			FROM messages m
			INNER JOIN devices d ON d.id = m.device_id
			WHERE d.user_id = $1 AND m.created_at >= CURRENT_DATE - ($2::int - 1) * INTERVAL '1 day'
			GROUP BY DATE(m.created_at)
			ORDER BY usage_date ASC
		`, userID, days)
		if err != nil {
			return nil, err
		}
	}
	defer rows.Close()

	type dailyUsage struct {
		Date   time.Time
		Sent   int64
		Failed int64
	}

	usageMap := make(map[string]dailyUsage, days)
	for rows.Next() {
		var entry dailyUsage
		if scanErr := rows.Scan(&entry.Date, &entry.Sent, &entry.Failed); scanErr != nil {
			return nil, scanErr
		}
		usageMap[entry.Date.Format("2006-01-02")] = entry
	}

	result := make([]models.BillingUsagePoint, 0, days)
	today := time.Now()
	start := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location()).AddDate(0, 0, -(days - 1))
	for i := 0; i < days; i++ {
		day := start.AddDate(0, 0, i)
		key := day.Format("2006-01-02")
		entry, ok := usageMap[key]
		if !ok {
			result = append(result, models.BillingUsagePoint{Date: day.Format("02 Jan"), Sent: 0, Failed: 0})
			continue
		}
		result = append(result, models.BillingUsagePoint{Date: day.Format("02 Jan"), Sent: entry.Sent, Failed: entry.Failed})
	}

	return result, nil
}

func (s *Service) getInvoiceHistory(ctx context.Context, userID uuid.UUID) ([]models.BillingInvoiceSummary, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT s.id,
		       COALESCE(s.start_date, s.created_at) AS invoice_date,
		       COALESCE(p.name, 'Unknown Plan') AS plan_name,
		       COALESCE(p.price, 0) AS amount,
		       COALESCE(s.status, 'inactive') AS status
		FROM subscriptions s
		LEFT JOIN billing_plans p ON p.id = s.plan_id
		WHERE s.user_id = $1
		ORDER BY COALESCE(s.start_date, s.created_at) DESC
		LIMIT 12
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	invoices := make([]models.BillingInvoiceSummary, 0)
	for rows.Next() {
		var invoice models.BillingInvoiceSummary
		if scanErr := rows.Scan(&invoice.SubscriptionID, &invoice.Date, &invoice.PlanName, &invoice.Amount, &invoice.Status); scanErr != nil {
			return nil, scanErr
		}
		invoice.ID = fmt.Sprintf("INV-%s-%s", invoice.Date.Format("2006-01"), invoice.SubscriptionID.String()[:8])
		invoices = append(invoices, invoice)
	}

	return invoices, nil
}

func normalizeBillingCycle(cycle string) string {
	if cycle == "" {
		return "monthly"
	}
	return cycle
}

func sumSentMessages(history []models.BillingUsagePoint) int64 {
	var total int64
	for _, point := range history {
		total += point.Sent
	}
	return total
}

func isMissingColumnError(err error, column string) bool {
	if err == nil {
		return false
	}

	message := strings.ToLower(err.Error())
	return strings.Contains(message, "column \""+strings.ToLower(column)+"\" does not exist")
}

func calculateSubscriptionEnd(start time.Time, cycle string) time.Time {
	switch strings.ToLower(cycle) {
	case "yearly", "annual":
		return start.AddDate(1, 0, 0)
	default:
		return start.AddDate(0, 1, 0)
	}
}
