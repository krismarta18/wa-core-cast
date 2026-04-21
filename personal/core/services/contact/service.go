package contact

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"wacast/core/models"
)

// Service handles contact-related business logic
type Service struct {
	store *Store
}

// NewService creates a new contact service
func NewService(store *Store) *Service {
	return &Service{store: store}
}

// CreateContact handles contact creation
func (s *Service) CreateContact(userID uuid.UUID, req models.CreateContactRequest) (*models.Contact, error) {
	additionalData, _ := json.Marshal(req.AdditionalData)
	
	c := &models.Contact{
		UserID:         userID,
		LabelID:        req.LabelID,
		Name:           req.Name,
		PhoneNumber:    req.PhoneNumber,
		Note:           req.Note,
		AdditionalData: additionalData,
	}

	if err := s.store.CreateContact(c); err != nil {
		return nil, err
	}
	return c, nil
}

// ListContacts returns all contacts for a user
func (s *Service) ListContacts(userID uuid.UUID) ([]*models.Contact, error) {
	return s.store.ListContacts(userID)
}

// UpdateContact handles contact updates
func (s *Service) UpdateContact(userID uuid.UUID, contactID uuid.UUID, req models.UpdateContactRequest) (*models.Contact, error) {
	c, err := s.store.GetContact(contactID)
	if err != nil {
		return nil, err
	}
	if c == nil || c.UserID != userID {
		return nil, fmt.Errorf("contact not found")
	}

	if req.Name != nil {
		c.Name = *req.Name
	}
	if req.PhoneNumber != nil {
		c.PhoneNumber = *req.PhoneNumber
	}
	if req.LabelID != nil {
		c.LabelID = req.LabelID
	}
	if req.Note != nil {
		c.Note = req.Note
	}
	if req.AdditionalData != nil {
		additionalData, _ := json.Marshal(req.AdditionalData)
		c.AdditionalData = additionalData
	}

	if err := s.store.UpdateContact(c); err != nil {
		return nil, err
	}
	return c, nil
}

// DeleteContact removes a contact
func (s *Service) DeleteContact(userID uuid.UUID, contactID uuid.UUID) error {
	return s.store.DeleteContact(contactID, userID)
}

// --- Group Management ---

// ListGroups returns all groups for a user
func (s *Service) ListGroups(userID uuid.UUID) ([]*models.ContactGroup, error) {
	return s.store.ListGroups(userID)
}

// CreateGroup handles group creation
func (s *Service) CreateGroup(userID uuid.UUID, req models.CreateContactGroupRequest) (*models.ContactGroup, error) {
	g := &models.ContactGroup{
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
	}
	if err := s.store.CreateGroup(g); err != nil {
		return nil, err
	}
	return g, nil
}

// DeleteGroup removes a group
func (s *Service) DeleteGroup(userID uuid.UUID, groupID uuid.UUID) error {
	return s.store.DeleteGroup(groupID, userID)
}

// AddMemberToGroup adds a contact to group
func (s *Service) AddMemberToGroup(userID uuid.UUID, groupID uuid.UUID, contactID uuid.UUID) error {
	// Verify group belongs to user
	groups, err := s.store.ListGroups(userID)
	if err != nil {
		return err
	}
	owned := false
	for _, g := range groups {
		if g.ID == groupID {
			owned = true
			break
		}
	}
	if !owned {
		return fmt.Errorf("forbidden")
	}

	return s.store.AddMemberToGroup(groupID, contactID)
}

// RemoveMemberFromGroup removes member
func (s *Service) RemoveMemberFromGroup(userID uuid.UUID, groupID uuid.UUID, contactID uuid.UUID) error {
	// Verify group ownership
	groups, err := s.store.ListGroups(userID)
	if err != nil {
		return err
	}
	owned := false
	for _, g := range groups {
		if g.ID == groupID {
			owned = true
			break
		}
	}
	if !owned {
		return fmt.Errorf("forbidden")
	}

	return s.store.RemoveMemberFromGroup(groupID, contactID)
}

// GetGroupMembers returns members
func (s *Service) GetGroupMembers(userID uuid.UUID, groupID uuid.UUID) ([]*models.Contact, error) {
	return s.store.GetGroupMembers(groupID)
}

// --- Blacklist Management ---

// ListBlacklist returns blacklist
func (s *Service) ListBlacklist(userID uuid.UUID) ([]map[string]interface{}, error) {
	return s.store.ListBlacklist(userID)
}

// BlacklistNumber blocks number
func (s *Service) BlacklistNumber(userID uuid.UUID, phone string, reason string) error {
	return s.store.BlacklistNumber(userID, phone, reason)
}

// UnblacklistNumber unblocks number
func (s *Service) UnblacklistNumber(userID uuid.UUID, id uuid.UUID) error {
	return s.store.UnblacklistNumber(id, userID)
}

// IsBlacklisted checks status
func (s *Service) IsBlacklisted(userID uuid.UUID, phone string) (bool, error) {
	return s.store.IsBlacklisted(userID, phone)
}
