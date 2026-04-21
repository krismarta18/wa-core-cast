package models

import (
	"time"

	"github.com/google/uuid"
)

// UsageQuota tracks message & device usage per billing period
type UsageQuota struct {
	ID                  uuid.UUID  `json:"id"`
	UserID              uuid.UUID  `json:"user_id"`
	SubscriptionID      uuid.UUID  `json:"subscription_id"`
	PeriodStart         time.Time  `json:"period_start"`
	PeriodEnd           time.Time  `json:"period_end"`
	MessagesUsed        int        `json:"messages_used"`
	MessagesLimit       int        `json:"messages_limit"`
	DevicesUsed         int        `json:"devices_used"`
	DevicesLimit        int        `json:"devices_limit"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

// TableName returns the table name
func (UsageQuota) TableName() string {
	return "usage_quotas"
}

// Invoice represents a billing invoice
type Invoice struct {
	ID             uuid.UUID  `json:"id"`
	UserID         uuid.UUID  `json:"user_id"`
	SubscriptionID uuid.UUID  `json:"subscription_id"`
	InvoiceNumber  string     `json:"invoice_number"`
	TotalAmount    float64    `json:"total_amount"`
	Currency       string     `json:"currency"`
	Status         string     `json:"status"` // draft/sent/paid/overdue/cancelled
	IssuedAt       time.Time  `json:"issued_at"`
	DueAt          time.Time  `json:"due_at"`
	PaidAt         *time.Time `json:"paid_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// TableName returns the table name
func (Invoice) TableName() string {
	return "invoices"
}

// InvoiceItem represents a line item on an invoice
type InvoiceItem struct {
	ID          uuid.UUID `json:"id"`
	InvoiceID   uuid.UUID `json:"invoice_id"`
	Description string    `json:"description"`
	Quantity    int       `json:"quantity"`
	UnitPrice   float64   `json:"unit_price"`
	TotalPrice  float64   `json:"total_price"`
	CreatedAt   time.Time `json:"created_at"`
}

// TableName returns the table name
func (InvoiceItem) TableName() string {
	return "invoice_items"
}
