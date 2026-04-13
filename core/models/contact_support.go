package models

import (
	"time"

	"github.com/google/uuid"
)

// ContactLabel represents a label/tag that can be applied to contacts
type ContactLabel struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Name      string    `json:"name"`
	Color     string    `json:"color"` // hex color code
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName returns the table name
func (ContactLabel) TableName() string {
	return "contact_labels"
}

// ContactGroupMember is the join table between contacts and groups
type ContactGroupMember struct {
	ContactID uuid.UUID `json:"contact_id"`
	GroupID   uuid.UUID `json:"group_id"`
	AddedAt   time.Time `json:"added_at"`
}

// TableName returns the table name
func (ContactGroupMember) TableName() string {
	return "contact_group_members"
}

// Blacklist holds phone numbers that are blocked from messaging
type Blacklist struct {
	ID          uuid.UUID  `json:"id"`
	UserID      uuid.UUID  `json:"user_id"`
	PhoneNumber string     `json:"phone_number"`
	Reason      *string    `json:"reason,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

// TableName returns the table name
func (Blacklist) TableName() string {
	return "blacklists"
}

// CreateContactLabelRequest is the request struct for creating a label
type CreateContactLabelRequest struct {
	Name  string `json:"name" binding:"required"`
	Color string `json:"color" binding:"required"`
}
