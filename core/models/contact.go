package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Contact represents a contact
type Contact struct {
	ID             uuid.UUID       `json:"id"`
	UserID         uuid.UUID       `json:"user_id"`
	LabelID        *uuid.UUID      `json:"label_id,omitempty"`
	Name           string          `json:"name"`
	PhoneNumber    string          `json:"phone_number"`
	Note           *string         `json:"note,omitempty"`
	AdditionalData json.RawMessage `json:"additional_data"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
	DeletedAt      *time.Time      `json:"deleted_at,omitempty"`
}

// TableName returns the table name
func (Contact) TableName() string {
	return "contacts"
}

// ContactGroup represents a contact group
type ContactGroup struct {
	ID          uuid.UUID  `json:"id"`
	UserID      uuid.UUID  `json:"user_id"`
	Name        string     `json:"name"`
	Description *string    `json:"description,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

// TableName returns the table name
func (ContactGroup) TableName() string {
	return "contact_groups"
}

// CreateContactRequest is the request struct for creating a contact
type CreateContactRequest struct {
	LabelID        *uuid.UUID             `json:"label_id"`
	Name           string                 `json:"name" binding:"required"`
	PhoneNumber    string                 `json:"phone_number" binding:"required"`
	Note           *string                `json:"note"`
	AdditionalData map[string]interface{} `json:"additional_data"`
}

// UpdateContactRequest is the request struct for updating a contact
type UpdateContactRequest struct {
	LabelID        *uuid.UUID             `json:"label_id"`
	Name           *string                `json:"name"`
	PhoneNumber    *string                `json:"phone_number"`
	Note           *string                `json:"note"`
	AdditionalData map[string]interface{} `json:"additional_data"`
}

// CreateContactGroupRequest is the request struct for creating a group
type CreateContactGroupRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`
}

// ContactResponse is the response struct
type ContactResponse struct {
	ID          uuid.UUID  `json:"id"`
	UserID      uuid.UUID  `json:"user_id"`
	LabelID     *uuid.UUID `json:"label_id,omitempty"`
	Name        string     `json:"name"`
	PhoneNumber string     `json:"phone_number"`
	Note        *string    `json:"note,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

// ToResponse converts Contact to response
func (c *Contact) ToResponse() *ContactResponse {
	return &ContactResponse{
		ID:          c.ID,
		UserID:      c.UserID,
		LabelID:     c.LabelID,
		Name:        c.Name,
		PhoneNumber: c.PhoneNumber,
		Note:        c.Note,
		CreatedAt:   c.CreatedAt,
	}
}

// ContactGroupResponse is the response struct
type ContactGroupResponse struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Description *string    `json:"description,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

// ToResponse converts ContactGroup to response
func (g *ContactGroup) ToResponse() *ContactGroupResponse {
	return &ContactGroupResponse{
		ID:          g.ID,
		Name:        g.Name,
		Description: g.Description,
		CreatedAt:   g.CreatedAt,
	}
}
