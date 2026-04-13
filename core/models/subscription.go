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
	RenewalDate *time.Time `json:"renewal_date,omitempty"`
	AutoRenew   bool       `json:"auto_renew"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
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
