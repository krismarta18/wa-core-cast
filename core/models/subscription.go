package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// BillingPlan represents a subscription plan
type BillingPlan struct {
	ID              uuid.UUID       `json:"id"`
	Name            string          `json:"name"`
	Price           float64         `json:"price"`
	MaxDevice       int32           `json:"max_device"`
	MaxMessagesDay  int32           `json:"max_messages_day"`
	Features        json.RawMessage `json:"features"`
	CreatedAt       *time.Time      `json:"created_at,omitempty"`
}

// TableName returns the table name
func (BillingPlan) TableName() string {
	return "billing_plans"
}

// Subscription represents a user subscription
type Subscription struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	PlanID    uuid.UUID `json:"plan_id"`
	Status    int32     `json:"status"` // 0: inactive, 1: active
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

// TableName returns the table name
func (Subscription) TableName() string {
	return "subscriptions"
}

// IsActive checks if subscription is active
func (s *Subscription) IsActive() bool {
	return s.Status == 1
}

// CreateSubscriptionRequest is the request struct for creating a subscription
type CreateSubscriptionRequest struct {
	UserID uuid.UUID `json:"user_id" binding:"required"`
	PlanID uuid.UUID `json:"plan_id" binding:"required"`
}

// UpdateSubscriptionRequest is the request struct for updating a subscription
type UpdateSubscriptionRequest struct {
	Status *int32 `json:"status"`
	PlanID *uuid.UUID `json:"plan_id"`
}

// SubscriptionResponse is the response struct for subscription endpoints
type SubscriptionResponse struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	PlanID    uuid.UUID `json:"plan_id"`
	Status    int32     `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// ToResponse converts Subscription to SubscriptionResponse
func (s *Subscription) ToResponse() *SubscriptionResponse {
	return &SubscriptionResponse{
		ID:        s.ID,
		UserID:    s.UserID,
		PlanID:    s.PlanID,
		Status:    s.Status,
		CreatedAt: s.CreatedAt,
	}
}
