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
		INSERT INTO contacts (id, user_id, label_id, name, phone_number, note, additional_data, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	now := time.Now()
	_, err := d.Exec(query,
		contact.ID, contact.UserID, contact.LabelID, contact.Name, contact.PhoneNumber,
		contact.Note, contact.AdditionalData, now, now,
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
		SELECT id, user_id, label_id, name, phone_number, note, additional_data, created_at, updated_at, deleted_at
		FROM contacts
		WHERE id = $1 AND deleted_at IS NULL
	`

	contact := &models.Contact{}
	err := d.QueryRow(query, contactID).Scan(
		&contact.ID, &contact.UserID, &contact.LabelID, &contact.Name, &contact.PhoneNumber,
		&contact.Note, &contact.AdditionalData, &contact.CreatedAt, &contact.UpdatedAt, &contact.DeletedAt,
	)

	if err != nil {
		return nil, err
	}

	return contact, nil
}

// GetContactsByUserID retrieves all contacts for a user
func (d *Database) GetContactsByUserID(userID uuid.UUID, limit, offset int) ([]models.Contact, error) {
	query := `
		SELECT id, user_id, label_id, name, phone_number, note, additional_data, created_at, updated_at, deleted_at
		FROM contacts
		WHERE user_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := d.Query(query, userID, limit, offset)
	if err != nil {
		utils.Error("Failed to get contacts", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	contacts := []models.Contact{}
	for rows.Next() {
		contact := models.Contact{}
		err := rows.Scan(
			&contact.ID, &contact.UserID, &contact.LabelID, &contact.Name, &contact.PhoneNumber,
			&contact.Note, &contact.AdditionalData, &contact.CreatedAt, &contact.UpdatedAt, &contact.DeletedAt,
		)
		if err != nil {
			utils.Error("Failed to scan contact", zap.Error(err))
			continue
		}
		contacts = append(contacts, contact)
	}

	return contacts, nil
}

// GetContactsByGroupID retrieves contacts belonging to a group via the join table
func (d *Database) GetContactsByGroupID(groupID uuid.UUID, limit, offset int) ([]models.Contact, error) {
	query := `
		SELECT c.id, c.user_id, c.label_id, c.name, c.phone_number, c.note,
		       c.additional_data, c.created_at, c.updated_at, c.deleted_at
		FROM contacts c
		INNER JOIN contact_group_members cgm ON cgm.contact_id = c.id
		WHERE cgm.group_id = $1 AND c.deleted_at IS NULL
		ORDER BY c.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := d.Query(query, groupID, limit, offset)
	if err != nil {
		utils.Error("Failed to get contacts by group", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	contacts := []models.Contact{}
	for rows.Next() {
		contact := models.Contact{}
		err := rows.Scan(
			&contact.ID, &contact.UserID, &contact.LabelID, &contact.Name, &contact.PhoneNumber,
			&contact.Note, &contact.AdditionalData, &contact.CreatedAt, &contact.UpdatedAt, &contact.DeletedAt,
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
	query := `UPDATE contacts SET `
	args := []interface{}{}
	argCount := 1

	if update.Name != nil {
		query += fmt.Sprintf("name = $%d", argCount)
		args = append(args, *update.Name)
		argCount++
	}

	if update.PhoneNumber != nil {
		if argCount > 1 {
			query += ", "
		}
		query += fmt.Sprintf("phone_number = $%d", argCount)
		args = append(args, *update.PhoneNumber)
		argCount++
	}

	if update.LabelID != nil {
		if argCount > 1 {
			query += ", "
		}
		query += fmt.Sprintf("label_id = $%d", argCount)
		args = append(args, *update.LabelID)
		argCount++
	}

	if update.Note != nil {
		if argCount > 1 {
			query += ", "
		}
		query += fmt.Sprintf("note = $%d", argCount)
		args = append(args, *update.Note)
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
	query := `UPDATE contacts SET deleted_at = $1 WHERE id = $2`

	_, err := d.Exec(query, time.Now(), contactID)
	if err != nil {
		utils.Error("Failed to delete contact", zap.Error(err))
		return err
	}

	return nil
}

// CreateContactGroup creates a new contact group
func (d *Database) CreateContactGroup(group *models.ContactGroup) error {
	query := `
		INSERT INTO contact_groups (id, user_id, name, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	now := time.Now()
	_, err := d.Exec(query,
		group.ID, group.UserID, group.Name, group.Description, now, now,
	)

	if err != nil {
		utils.Error("Failed to create contact group", zap.Error(err))
		return err
	}

	return nil
}

// GetContactGroupByID retrieves a contact group by ID
func (d *Database) GetContactGroupByID(groupID uuid.UUID) (*models.ContactGroup, error) {
	query := `
		SELECT id, user_id, name, description, created_at, updated_at, deleted_at
		FROM contact_groups
		WHERE id = $1 AND deleted_at IS NULL
	`

	group := &models.ContactGroup{}
	err := d.QueryRow(query, groupID).Scan(
		&group.ID, &group.UserID, &group.Name, &group.Description,
		&group.CreatedAt, &group.UpdatedAt, &group.DeletedAt,
	)

	if err != nil {
		return nil, err
	}

	return group, nil
}

// GetContactGroupsByUserID retrieves all contact groups for a user
func (d *Database) GetContactGroupsByUserID(userID uuid.UUID) ([]models.ContactGroup, error) {
	query := `
		SELECT id, user_id, name, description, created_at, updated_at, deleted_at
		FROM contact_groups
		WHERE user_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
	`

	rows, err := d.Query(query, userID)
	if err != nil {
		utils.Error("Failed to get contact groups", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	groups := []models.ContactGroup{}
	for rows.Next() {
		group := models.ContactGroup{}
		err := rows.Scan(
			&group.ID, &group.UserID, &group.Name, &group.Description,
			&group.CreatedAt, &group.UpdatedAt, &group.DeletedAt,
		)
		if err != nil {
			utils.Error("Failed to scan group", zap.Error(err))
			continue
		}
		groups = append(groups, group)
	}

	return groups, nil
}


