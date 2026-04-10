package database

import (
	"fmt"
	"time"

	"wacast/core/models"
	"wacast/core/utils"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// CreateContact creates a new contact
func (d *Database) CreateContact(contact *models.Contact) error {
	query := `
		INSERT INTO contact (id, group_id, name, phone, additional_data, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := d.Exec(query,
		contact.ID, contact.GroupID, contact.Name, contact.Phone,
		contact.AdditionalData, time.Now(), time.Now(),
	)

	if err != nil {
		utils.Error("Failed to create contact", zap.Error(err))
		return err
	}

	return nil
}

// GetContactByID retrieves a contact by ID
func (d *Database) GetContactByID(contactID uuid.UUID) (*models.Contact, error) {
	query := `
		SELECT id, group_id, name, phone, additional_data, created_at, updated_at, deleted_at
		FROM contact
		WHERE id = $1 AND deleted_at IS NULL
	`

	contact := &models.Contact{}
	err := d.QueryRow(query, contactID).Scan(
		&contact.ID, &contact.GroupID, &contact.Name, &contact.Phone,
		&contact.AdditionalData, &contact.CreatedAt, &contact.UpdatedAt, &contact.DeletedAt,
	)

	if err != nil {
		return nil, err
	}

	return contact, nil
}

// GetContactsByGroupID retrieves contacts by group
func (d *Database) GetContactsByGroupID(groupID uuid.UUID, limit, offset int) ([]models.Contact, error) {
	query := `
		SELECT id, group_id, name, phone, additional_data, created_at, updated_at, deleted_at
		FROM contact
		WHERE group_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := d.Query(query, groupID, limit, offset)
	if err != nil {
		utils.Error("Failed to get contacts", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	contacts := []models.Contact{}
	for rows.Next() {
		contact := models.Contact{}
		err := rows.Scan(
			&contact.ID, &contact.GroupID, &contact.Name, &contact.Phone,
			&contact.AdditionalData, &contact.CreatedAt, &contact.UpdatedAt, &contact.DeletedAt,
		)
		if err != nil {
			utils.Error("Failed to scan contact", zap.Error(err))
			continue
		}
		contacts = append(contacts, contact)
	}

	return contacts, nil
}

// UpdateContact updates a contact
func (d *Database) UpdateContact(contactID uuid.UUID, update *models.UpdateContactRequest) error {
	query := `UPDATE contact SET `
	args := []interface{}{}
	argCount := 1

	if update.Name != nil {
		query += fmt.Sprintf("name = $%d", argCount)
		args = append(args, *update.Name)
		argCount++
	}

	if update.Phone != nil {
		if argCount > 1 {
			query += ", "
		}
		query += fmt.Sprintf("phone = $%d", argCount)
		args = append(args, *update.Phone)
		argCount++
	}

	if update.AdditionalData != nil {
		if argCount > 1 {
			query += ", "
		}
		query += fmt.Sprintf("additional_data = $%d", argCount)
		args = append(args, update.AdditionalData)
		argCount++
	}

	if argCount > 1 {
		query += ", "
	}
	query += fmt.Sprintf("updated_at = $%d WHERE id = $%d", argCount, argCount+1)
	args = append(args, time.Now(), contactID)

	_, err := d.Exec(query, args...)
	if err != nil {
		utils.Error("Failed to update contact", zap.Error(err))
		return err
	}

	return nil
}

// DeleteContact soft deletes a contact
func (d *Database) DeleteContact(contactID uuid.UUID) error {
	query := `UPDATE contact SET deleted_at = $1 WHERE id = $2`

	_, err := d.Exec(query, time.Now(), contactID)
	if err != nil {
		utils.Error("Failed to delete contact", zap.Error(err))
		return err
	}

	return nil
}

// CreateGroup creates a new group
func (d *Database) CreateGroup(group *models.Group) error {
	query := `
		INSERT INTO groups (id, user_id, group_name, created_at)
		VALUES ($1, $2, $3, $4)
	`

	_, err := d.Exec(query,
		group.ID, group.UserID, group.GroupName, time.Now(),
	)

	if err != nil {
		utils.Error("Failed to create group", zap.Error(err))
		return err
	}

	return nil
}

// GetGroupByID retrieves a group by ID
func (d *Database) GetGroupByID(groupID uuid.UUID) (*models.Group, error) {
	query := `
		SELECT id, user_id, group_name, created_at, deleted_at
		FROM groups
		WHERE id = $1 AND deleted_at IS NULL
	`

	group := &models.Group{}
	err := d.QueryRow(query, groupID).Scan(
		&group.ID, &group.UserID, &group.GroupName, &group.CreatedAt, &group.DeletedAt,
	)

	if err != nil {
		return nil, err
	}

	return group, nil
}

// GetGroupsByUserID retrieves all groups for a user
func (d *Database) GetGroupsByUserID(userID uuid.UUID) ([]models.Group, error) {
	query := `
		SELECT id, user_id, group_name, created_at, deleted_at
		FROM groups
		WHERE user_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
	`

	rows, err := d.Query(query, userID)
	if err != nil {
		utils.Error("Failed to get groups", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	groups := []models.Group{}
	for rows.Next() {
		group := models.Group{}
		err := rows.Scan(
			&group.ID, &group.UserID, &group.GroupName, &group.CreatedAt, &group.DeletedAt,
		)
		if err != nil {
			utils.Error("Failed to scan group", zap.Error(err))
			continue
		}
		groups = append(groups, group)
	}

	return groups, nil
}

// UpdateGroup updates a group
func (d *Database) UpdateGroup(groupID uuid.UUID, groupName string) error {
	query := `UPDATE groups SET group_name = $1 WHERE id = $2`

	_, err := d.Exec(query, groupName, groupID)
	if err != nil {
		utils.Error("Failed to update group", zap.Error(err))
		return err
	}

	return nil
}

// DeleteGroup soft deletes a group
func (d *Database) DeleteGroup(groupID uuid.UUID) error {
	query := `UPDATE groups SET deleted_at = $1 WHERE id = $2`

	_, err := d.Exec(query, time.Now(), groupID)
	if err != nil {
		utils.Error("Failed to delete group", zap.Error(err))
		return err
	}

	return nil
}

// CountGroupContacts counts contacts in a group
func (d *Database) CountGroupContacts(groupID uuid.UUID) (int64, error) {
	query := `SELECT COUNT(*) FROM contact WHERE group_id = $1 AND deleted_at IS NULL`

	var count int64
	err := d.QueryRow(query, groupID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// GetContactByPhone retrieves a contact by phone number
func (d *Database) GetContactByPhone(phone string) (*models.Contact, error) {
	query := `
		SELECT id, group_id, name, phone, additional_data, created_at, updated_at, deleted_at
		FROM contact
		WHERE phone = $1 AND deleted_at IS NULL
		LIMIT 1
	`

	contact := &models.Contact{}
	err := d.QueryRow(query, phone).Scan(
		&contact.ID, &contact.GroupID, &contact.Name, &contact.Phone,
		&contact.AdditionalData, &contact.CreatedAt, &contact.UpdatedAt, &contact.DeletedAt,
	)

	if err != nil {
		return nil, err
	}

	return contact, nil
}
