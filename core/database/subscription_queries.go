package database

import (
	"wacast/core/models"
	"wacast/core/utils"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// CreateSubscription creates a new subscription
func (d *Database) CreateSubscription(subscription *models.Subscription) error {
	query := `
		INSERT INTO subscriptions (id, user_id, plan_id, status, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := d.Exec(query,
		subscription.ID, subscription.UserID, subscription.PlanID, 
		subscription.Status, subscription.CreatedAt,
	)

	if err != nil {
		utils.Error("Failed to create subscription", zap.Error(err))
		return err
	}

	utils.Debug("Subscription created", zap.String("subscription_id", subscription.ID.String()))
	return nil
}

// GetSubscriptionByID retrieves a subscription by ID
func (d *Database) GetSubscriptionByID(subID uuid.UUID) (*models.Subscription, error) {
	query := `
		SELECT id, user_id, plan_id, status, created_at
		FROM subscriptions
		WHERE id = $1
	`

	subscription := &models.Subscription{}
	err := d.QueryRow(query, subID).Scan(
		&subscription.ID, &subscription.UserID, &subscription.PlanID, 
		&subscription.Status, &subscription.CreatedAt,
	)

	if err != nil {
		utils.Debug("Subscription not found", zap.String("subscription_id", subID.String()))
		return nil, err
	}

	return subscription, nil
}

// GetSubscriptionByUserID retrieves subscription for a user
func (d *Database) GetSubscriptionByUserID(userID uuid.UUID) (*models.Subscription, error) {
	query := `
		SELECT id, user_id, plan_id, status, created_at
		FROM subscriptions
		WHERE user_id = $1 AND status = 1
		ORDER BY created_at DESC
		LIMIT 1
	`

	subscription := &models.Subscription{}
	err := d.QueryRow(query, userID).Scan(
		&subscription.ID, &subscription.UserID, &subscription.PlanID, 
		&subscription.Status, &subscription.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return subscription, nil
}

// UpdateSubscription updates subscription
func (d *Database) UpdateSubscription(subID uuid.UUID, update *models.UpdateSubscriptionRequest) error {
	query := `UPDATE subscriptions SET `
	args := []interface{}{}
	updated := false

	if update.Status != nil {
		query += "status = $1"
		args = append(args, *update.Status)
		updated = true
	}

	if update.PlanID != nil {
		if updated {
			query += ", "
		}
		query += "plan_id = $2"
		args = append(args, *update.PlanID)
		updated = true
	}

	if !updated {
		return nil
	}

	args = append(args, subID)
	query += " WHERE id = $" + string(rune(len(args)))

	_, err := d.Exec(query, args...)
	if err != nil {
		utils.Error("Failed to update subscription", zap.Error(err))
		return err
	}

	return nil
}

// GetBillingPlanByID retrieves a billing plan by ID
func (d *Database) GetBillingPlanByID(planID uuid.UUID) (*models.BillingPlan, error) {
	query := `
		SELECT id, name, price, max_device, max_messages_day, features
		FROM billing_plans
		WHERE id = $1
	`

	plan := &models.BillingPlan{}
	err := d.QueryRow(query, planID).Scan(
		&plan.ID, &plan.Name, &plan.Price, &plan.MaxDevice, 
		&plan.MaxMessagesDay, &plan.Features,
	)

	if err != nil {
		utils.Debug("Billing plan not found", zap.String("plan_id", planID.String()))
		return nil, err
	}

	return plan, nil
}

// GetAllBillingPlans retrieves all billing plans
func (d *Database) GetAllBillingPlans() ([]models.BillingPlan, error) {
	query := `
		SELECT id, name, price, max_device, max_messages_day, features
		FROM billing_plans
		ORDER BY price ASC
	`

	rows, err := d.Query(query)
	if err != nil {
		utils.Error("Failed to get billing plans", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	plans := []models.BillingPlan{}
	for rows.Next() {
		plan := models.BillingPlan{}
		err := rows.Scan(
			&plan.ID, &plan.Name, &plan.Price, &plan.MaxDevice, 
			&plan.MaxMessagesDay, &plan.Features,
		)
		if err != nil {
			utils.Error("Failed to scan billing plan", zap.Error(err))
			continue
		}
		plans = append(plans, plan)
	}

	return plans, nil
}

// CreateBillingPlan creates a new billing plan
func (d *Database) CreateBillingPlan(plan *models.BillingPlan) error {
	query := `
		INSERT INTO billing_plans (id, name, price, max_device, max_messages_day, features)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := d.Exec(query,
		plan.ID, plan.Name, plan.Price, plan.MaxDevice, 
		plan.MaxMessagesDay, plan.Features,
	)

	if err != nil {
		utils.Error("Failed to create billing plan", zap.Error(err))
		return err
	}

	return nil
}

// GetSubscriptionCount returns total active subscriptions
func (d *Database) GetSubscriptionCount() (int64, error) {
	query := `SELECT COUNT(*) FROM subscriptions WHERE status = 1`

	var count int64
	err := d.QueryRow(query).Scan(&count)
	if err != nil {
		utils.Error("Failed to get subscription count", zap.Error(err))
		return 0, err
	}

	return count, nil
}

// DeactivateSubscription deactivates a subscription
func (d *Database) DeactivateSubscription(subID uuid.UUID) error {
	query := `UPDATE subscriptions SET status = 0 WHERE id = $1`

	_, err := d.Exec(query, subID)
	if err != nil {
		utils.Error("Failed to deactivate subscription", zap.Error(err))
		return err
	}

	return nil
}

// ActivateSubscription activates a subscription
func (d *Database) ActivateSubscription(subID uuid.UUID) error {
	query := `UPDATE subscriptions SET status = 1 WHERE id = $1`

	_, err := d.Exec(query, subID)
	if err != nil {
		utils.Error("Failed to activate subscription", zap.Error(err))
		return err
	}

	return nil
}
