package contact

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"

	"wacast/core/database"
	"wacast/core/models"
)

// Store handles database operations for contacts
type Store struct {
	db *database.Database
}

// NewStore creates a new contact store
func NewStore(db *database.Database) *Store {
	return &Store{db: db}
}

// --- Contact Management ---

// CreateContact creates a new contact
func (s *Store) CreateContact(c *models.Contact) error {
	query := `
		INSERT INTO contacts (id, user_id, label_id, name, phone_number, additional_data, note, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	now := time.Now()
	c.CreatedAt = now
	c.UpdatedAt = now

	_, err := s.db.Exec(query, c.ID, c.UserID, c.LabelID, c.Name, c.PhoneNumber, c.AdditionalData, c.Note, c.CreatedAt, c.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create contact: %w", err)
	}
	return nil
}

// GetContact gets a contact by ID
func (s *Store) GetContact(id uuid.UUID) (*models.Contact, error) {
	query := `SELECT id, user_id, label_id, name, phone_number, additional_data, note, created_at, updated_at FROM contacts WHERE id = $1 AND deleted_at IS NULL`
	var c models.Contact
	err := s.db.QueryRow(query, id).Scan(&c.ID, &c.UserID, &c.LabelID, &c.Name, &c.PhoneNumber, &c.AdditionalData, &c.Note, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get contact: %w", err)
	}
	return &c, nil
}

// ListContacts lists contacts for a user
func (s *Store) ListContacts(userID uuid.UUID) ([]*models.Contact, error) {
	query := `SELECT id, user_id, label_id, name, phone_number, additional_data, note, created_at, updated_at FROM contacts WHERE user_id = $1 AND deleted_at IS NULL ORDER BY name ASC`
	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list contacts: %w", err)
	}
	defer rows.Close()

	var contacts []*models.Contact
	for rows.Next() {
		var c models.Contact
		if err := rows.Scan(&c.ID, &c.UserID, &c.LabelID, &c.Name, &c.PhoneNumber, &c.AdditionalData, &c.Note, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		contacts = append(contacts, &c)
	}
	return contacts, nil
}

// UpdateContact updates a contact
func (s *Store) UpdateContact(c *models.Contact) error {
	query := `
		UPDATE contacts 
		SET label_id = $1, name = $2, phone_number = $3, additional_data = $4, note = $5, updated_at = $6
		WHERE id = $7 AND user_id = $8
	`
	c.UpdatedAt = time.Now()
	_, err := s.db.Exec(query, c.LabelID, c.Name, c.PhoneNumber, c.AdditionalData, c.Note, c.UpdatedAt, c.ID, c.UserID)
	if err != nil {
		return fmt.Errorf("failed to update contact: %w", err)
	}
	return nil
}

// DeleteContact soft deletes a contact
func (s *Store) DeleteContact(id uuid.UUID, userID uuid.UUID) error {
	query := `UPDATE contacts SET deleted_at = $1 WHERE id = $2 AND user_id = $3`
	_, err := s.db.Exec(query, time.Now(), id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete contact: %w", err)
	}
	return nil
}

// --- Group Management ---

// CreateGroup creates a new contact group
func (s *Store) CreateGroup(g *models.ContactGroup) error {
	query := `
		INSERT INTO contact_groups (id, user_id, name, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	if g.ID == uuid.Nil {
		g.ID = uuid.New()
	}
	now := time.Now()
	g.CreatedAt = now
	g.UpdatedAt = now

	_, err := s.db.Exec(query, g.ID, g.UserID, g.Name, g.Description, g.CreatedAt, g.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create group: %w", err)
	}
	return nil
}

// ListGroups lists groups for a user
func (s *Store) ListGroups(userID uuid.UUID) ([]*models.ContactGroup, error) {
	query := `SELECT id, user_id, name, description, created_at, updated_at FROM contact_groups WHERE user_id = $1 AND deleted_at IS NULL ORDER BY name ASC`
	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list groups: %w", err)
	}
	defer rows.Close()

	var groups []*models.ContactGroup
	for rows.Next() {
		var g models.ContactGroup
		if err := rows.Scan(&g.ID, &g.UserID, &g.Name, &g.Description, &g.CreatedAt, &g.UpdatedAt); err != nil {
			return nil, err
		}
		groups = append(groups, &g)
	}
	return groups, nil
}

// DeleteGroup soft deletes a group
func (s *Store) DeleteGroup(id uuid.UUID, userID uuid.UUID) error {
	query := `UPDATE contact_groups SET deleted_at = $1 WHERE id = $2 AND user_id = $3`
	_, err := s.db.Exec(query, time.Now(), id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete group: %w", err)
	}
	return nil
}

// AddMemberToGroup adds a contact to a group
func (s *Store) AddMemberToGroup(groupID uuid.UUID, contactID uuid.UUID) error {
	query := `INSERT INTO contact_group_members (id, group_id, contact_id, created_at) VALUES ($1, $2, $3, $4)`
	_, err := s.db.Exec(query, uuid.New(), groupID, contactID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to add member to group: %w", err)
	}
	return nil
}

// RemoveMemberFromGroup removes a contact from a group
func (s *Store) RemoveMemberFromGroup(groupID uuid.UUID, contactID uuid.UUID) error {
	query := `DELETE FROM contact_group_members WHERE group_id = $1 AND contact_id = $2`
	_, err := s.db.Exec(query, groupID, contactID)
	if err != nil {
		return fmt.Errorf("failed to remove member from group: %w", err)
	}
	return nil
}

// GetGroupMembers lists members of a group
func (s *Store) GetGroupMembers(groupID uuid.UUID) ([]*models.Contact, error) {
	query := `
		SELECT c.id, c.user_id, c.label_id, c.name, c.phone_number, c.additional_data, c.note, c.created_at, c.updated_at 
		FROM contacts c
		JOIN contact_group_members cgm ON c.id = cgm.contact_id
		WHERE cgm.group_id = $1 AND c.deleted_at IS NULL
	`
	rows, err := s.db.Query(query, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to list group members: %w", err)
	}
	defer rows.Close()

	var contacts []*models.Contact
	for rows.Next() {
		var c models.Contact
		if err := rows.Scan(&c.ID, &c.UserID, &c.LabelID, &c.Name, &c.PhoneNumber, &c.AdditionalData, &c.Note, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		contacts = append(contacts, &c)
	}
	return contacts, nil
}

// --- Blacklist Management ---

// BlacklistNumber blocks a number
func (s *Store) BlacklistNumber(userID uuid.UUID, phone string, reason string) error {
	query := `INSERT INTO blacklists (id, user_id, phone_number, reason, blocked_at, created_at) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := s.db.Exec(query, uuid.New(), userID, phone, reason, time.Now(), time.Now())
	if err != nil {
		return fmt.Errorf("failed to blacklist number: %w", err)
	}
	return nil
}

// UnblacklistNumber unblocks a number
func (s *Store) UnblacklistNumber(id uuid.UUID, userID uuid.UUID) error {
	query := `DELETE FROM blacklists WHERE id = $1 AND user_id = $2`
	_, err := s.db.Exec(query, id, userID)
	if err != nil {
		return fmt.Errorf("failed to unblacklist number: %w", err)
	}
	return nil
}

// ListBlacklist lists blacklisted numbers for a user
func (s *Store) ListBlacklist(userID uuid.UUID) ([]map[string]interface{}, error) {
	query := `SELECT id, phone_number, reason, blocked_at FROM blacklists WHERE user_id = $1 AND unblocked_at IS NULL ORDER BY blocked_at DESC`
	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list blacklist: %w", err)
	}
	defer rows.Close()

	var list []map[string]interface{}
	for rows.Next() {
		var id uuid.UUID
		var phone, reason string
		var blockedAt time.Time
		if err := rows.Scan(&id, &phone, &reason, &blockedAt); err != nil {
			return nil, err
		}
		list = append(list, map[string]interface{}{
			"id":           id,
			"phone_number": phone,
			"reason":       reason,
			"blocked_at":   blockedAt,
		})
	}
	return list, nil
}

// IsBlacklisted checks if a number is blacklisted for a user
func (s *Store) IsBlacklisted(userID uuid.UUID, phone string) (bool, error) {
	query := `SELECT COUNT(*) FROM blacklists WHERE user_id = $1 AND phone_number = $2 AND unblocked_at IS NULL`
	var count int
	err := s.db.QueryRow(query, userID, phone).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
