package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Subscription status constants
const (
	SubscriptionStatusActive    = "active"
	SubscriptionStatusInactive  = "inactive"
	SubscriptionStatusExpired   = "expired"
	SubscriptionStatusCancelled = "cancelled"
)

// BillingPlan represents a subscription plan
type BillingPlan struct {
	ID                 uuid.UUID       `json:"id"`
	Name               string          `json:"name"`
	Price              float64         `json:"price"`
	BillingCycle       string          `json:"billing_cycle"` // monthly / yearly
	MaxDevices         int             `json:"max_devices"`
	MaxMessagesPerDay  int             `json:"max_messages_per_day"`
	Features           json.RawMessage `json:"features"`
	IsActive           bool            `json:"is_active"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at"`
}

// TableName returns the table name
func (BillingPlan) TableName() string {
	return "billing_plans"
}

// Subscription represents a user subscription
type Subscription struct {
	ID          uuid.UUID  `json:"id"`
	UserID      uuid.UUID  `json:"user_id"`
	PlanID      uuid.UUID  `json:"plan_id"`
	Status      string     `json:"status"` // active/inactive/expired/cancelled
	StartDate   *time.Time `json:"start_date,omitempty"`
	EndDate     *time.Time `json:"end_date,omitempty"`
	RenewalDate       *time.Time `json:"renewal_date,omitempty"`
	AutoRenew         bool       `json:"auto_renew"`
	MaxDevices        int        `json:"max_devices"`
	MaxMessagesPerDay int        `json:"max_messages_per_day"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

// TableName returns the table name
func (Subscription) TableName() string {
	return "subscriptions"
}

// IsActive checks if subscription is active
func (s *Subscription) IsActive() bool {
	return s.Status == SubscriptionStatusActive
}

// CreateSubscriptionRequest is the request struct for creating a subscription
type CreateSubscriptionRequest struct {
	UserID    uuid.UUID `json:"user_id" binding:"required"`
	PlanID    uuid.UUID `json:"plan_id" binding:"required"`
	AutoRenew bool      `json:"auto_renew"`
}

// UpdateSubscriptionRequest is the request struct for updating a subscription
type UpdateSubscriptionRequest struct {
	Status    *string    `json:"status"`
	PlanID    *uuid.UUID `json:"plan_id"`
	AutoRenew *bool      `json:"auto_renew"`
}

// SubscriptionResponse is the response struct for subscription endpoints
type SubscriptionResponse struct {
	ID          uuid.UUID  `json:"id"`
	UserID      uuid.UUID  `json:"user_id"`
	PlanID      uuid.UUID  `json:"plan_id"`
	Status      string     `json:"status"`
	StartDate   *time.Time `json:"start_date,omitempty"`
	EndDate     *time.Time `json:"end_date,omitempty"`
	RenewalDate *time.Time `json:"renewal_date,omitempty"`
	AutoRenew   bool       `json:"auto_renew"`
	CreatedAt   time.Time  `json:"created_at"`
}

// ToResponse converts Subscription to SubscriptionResponse
func (s *Subscription) ToResponse() *SubscriptionResponse {
	return &SubscriptionResponse{
		ID:          s.ID,
		UserID:      s.UserID,
		PlanID:      s.PlanID,
		Status:      s.Status,
		StartDate:   s.StartDate,
		EndDate:     s.EndDate,
		RenewalDate: s.RenewalDate,
		AutoRenew:   s.AutoRenew,
		CreatedAt:   s.CreatedAt,
	}
}

// BillingOverviewResponse is the dashboard response for billing & quota.
type BillingOverviewResponse struct {
	CurrentPlan  *BillingCurrentPlanResponse `json:"current_plan,omitempty"`
	UsageHistory []BillingUsagePoint         `json:"usage_history"`
	Plans        []BillingPlanSummary        `json:"plans"`
	Invoices     []BillingInvoiceSummary     `json:"invoices"`
}

type BillingCheckoutRequest struct {
	PlanID string `json:"plan_id" binding:"required"`
}

type BillingCheckoutResponse struct {
	Subscription *BillingCurrentPlanResponse `json:"subscription"`
	Invoice      BillingInvoiceSummary       `json:"invoice"`
	PaymentStatus string                     `json:"payment_status"`
	PaymentMethod string                     `json:"payment_method"`
}

type BillingCurrentPlanResponse struct {
	SubscriptionID uuid.UUID   `json:"subscription_id"`
	PlanID         uuid.UUID   `json:"plan_id"`
	Name           string      `json:"name"`
	Price          float64     `json:"price"`
	BillingCycle   string      `json:"billing_cycle"`
	RenewalDate    *time.Time  `json:"renewal_date,omitempty"`
	QuotaUsed      int64       `json:"quota_used"`
	QuotaLimit     int         `json:"quota_limit"`
	DeviceUsed     int         `json:"device_used"`
	DeviceMax      int         `json:"device_max"`
	AutoRenew      bool        `json:"auto_renew"`
	Status         string      `json:"status"`
	Features       json.RawMessage `json:"features,omitempty"`
}

type BillingUsagePoint struct {
	Date   string `json:"date"`
	Sent   int64  `json:"sent"`
	Failed int64  `json:"failed"`
}

type BillingPlanSummary struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Price      float64   `json:"price"`
	QuotaLimit int       `json:"quota_limit"`
	DeviceMax  int       `json:"device_max"`
	Current    bool      `json:"current"`
	IsActive   bool      `json:"is_active"`
}

type BillingInvoiceSummary struct {
	ID          string     `json:"id"`
	SubscriptionID uuid.UUID `json:"subscription_id"`
	Date        time.Time  `json:"date"`
	PlanName    string     `json:"plan_name"`
	Amount      float64    `json:"amount"`
	Status      string     `json:"status"`
}
