package broadcast

import (
	"wacast/core/database"
	"wacast/core/models"

	"github.com/google/uuid"
)

type Store struct {
	db *database.Database
}

func NewStore(db *database.Database) *Store {
	return &Store{db: db}
}

func (s *Store) CreateCampaign(campaign *models.BroadcastCampaign) error {
	return s.db.CreateBroadcastCampaign(campaign)
}

func (s *Store) GetCampaign(id uuid.UUID) (*models.BroadcastCampaign, error) {
	return s.db.GetBroadcastCampaignByID(id)
}

func (s *Store) ListCampaigns(userID uuid.UUID, limit, offset int) ([]models.BroadcastCampaign, error) {
	return s.db.GetBroadcastCampaignsByUserID(userID, limit, offset)
}

func (s *Store) UpdateCampaignStatus(id uuid.UUID, status string) error {
	return s.db.UpdateBroadcastCampaignStatus(id, status)
}

func (s *Store) UpdateCampaignProgress(id uuid.UUID, success, failed int) error {
	return s.db.UpdateBroadcastCampaignProgress(id, success, failed)
}

func (s *Store) CreateRecipient(recipient *models.BroadcastRecipient) error {
	return s.db.CreateBroadcastRecipient(recipient)
}

func (s *Store) ListRecipients(campaignID uuid.UUID, limit, offset int) ([]models.BroadcastRecipient, error) {
	return s.db.GetBroadcastRecipientsByCampaignID(campaignID, limit, offset)
}

func (s *Store) UpdateRecipientStatus(id uuid.UUID, status string, errMsg *string) error {
	return s.db.UpdateBroadcastRecipientStatus(id, status, errMsg)
}

func (s *Store) GetPendingRecipients(campaignID uuid.UUID, limit int) ([]models.BroadcastRecipient, error) {
	return s.db.GetPendingBroadcastRecipients(campaignID, limit)
}

func (s *Store) DeleteCampaign(id uuid.UUID) error {
	return s.db.DeleteBroadcastCampaign(id)
}
