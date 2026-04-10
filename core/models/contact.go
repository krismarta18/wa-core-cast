package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Contact represents a contact
type Contact struct {
	ID               uuid.UUID       `json:"id"`
	GroupID          uuid.UUID       `json:"group_id"`
	Name             string          `json:"name"`
	Phone            string          `json:"phone"`
	AdditionalData   json.RawMessage `json:"additional_data"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        *time.Time      `json:"updated_at,omitempty"`
	DeletedAt        *time.Time      `json:"deleted_at,omitempty"`
}

// TableName returns the table name
func (Contact) TableName() string {
	return "contact"
}

// Group represents a contact group
type Group struct {
	ID        uuid.UUID  `json:"id"`
	UserID    uuid.UUID  `json:"user_id"`
	GroupName string     `json:"group_name"`
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// TableName returns the table name
func (Group) TableName() string {
	return "groups"
}

// CreateContactRequest is the request struct for creating a contact
type CreateContactRequest struct {
	GroupID          uuid.UUID              `json:"group_id" binding:"required"`
	Name             string                 `json:"name" binding:"required"`
	Phone            string                 `json:"phone" binding:"required"`
	AdditionalData   map[string]interface{} `json:"additional_data"`
}

// UpdateContactRequest is the request struct for updating a contact
type UpdateContactRequest struct {
	Name           *string                 `json:"name"`
	Phone          *string                 `json:"phone"`
	AdditionalData map[string]interface{}  `json:"additional_data"`
}

// CreateGroupRequest is the request struct for creating a group
type CreateGroupRequest struct {
	GroupName string `json:"group_name" binding:"required"`
}

// ContactResponse is the response struct
type ContactResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
}

// ToResponse converts Contact to response
func (c *Contact) ToResponse() *ContactResponse {
	return &ContactResponse{
		ID:        c.ID,
		Name:      c.Name,
		Phone:     c.Phone,
		CreatedAt: c.CreatedAt,
	}
}

// GroupResponse is the response struct
type GroupResponse struct {
	ID        uuid.UUID `json:"id"`
	GroupName string    `json:"group_name"`
	CreatedAt time.Time `json:"created_at"`
}

// ToResponse converts Group to response
func (g *Group) ToResponse() *GroupResponse {
	return &GroupResponse{
		ID:        g.ID,
		GroupName: g.GroupName,
		CreatedAt: g.CreatedAt,
	}
}
