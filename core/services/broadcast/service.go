package broadcast

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"wacast/core/models"
	"wacast/core/services/message"
	"wacast/core/utils"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Service struct {
	store          *Store
	messageService *message.Service
}

func NewService(store *Store, messageService *message.Service) *Service {
	s := &Service{
		store:          store,
		messageService: messageService,
	}

	// Register callback to track broadcast completion status
	messageService.RegisterDeliveryCallback(s.handleMessageStatusUpdate)

	return s
}

func (s *Service) CreateCampaign(userID uuid.UUID, req models.CreateBroadcastCampaignRequest, recipients []string) (*models.BroadcastCampaign, error) {
	campaign := &models.BroadcastCampaign{
		ID:              uuid.New(),
		UserID:          userID,
		DeviceID:        req.DeviceID,
		TemplateID:      req.TemplateID,
		Name:            req.Name,
		MessageContent:  req.MessageContent,
		DelaySeconds:    req.DelaySeconds,
		ScheduledAt:     req.ScheduledAt,
		TotalRecipients: len(recipients),
		Status:          models.BroadcastStatusDraft,
	}

	if err := s.store.CreateCampaign(campaign); err != nil {
		return nil, err
	}

	// Create recipients
	for _, phone := range recipients {
		recipient := &models.BroadcastRecipient{
			ID:          uuid.New(),
			CampaignID:  campaign.ID,
			PhoneNumber: phone,
			Status:      "pending",
		}
		if err := s.store.CreateRecipient(recipient); err != nil {
			utils.Error("Failed to create broadcast recipient", zap.Error(err))
		}
	}

	return campaign, nil
}

func (s *Service) StartCampaign(ctx context.Context, campaignID uuid.UUID) error {
	campaign, err := s.store.GetCampaign(campaignID)
	if err != nil {
		return err
	}

	if campaign.Status != models.BroadcastStatusDraft {
		return fmt.Errorf("campaign already started or completed")
	}

	recipients, err := s.store.GetPendingRecipients(campaignID, 10000)
	if err != nil {
		return err
	}

	// Update status to sending
	if err := s.store.UpdateCampaignStatus(campaignID, models.BroadcastStatusSending); err != nil {
		return err
	}

	go s.enqueueBroadcastMessages(campaign, recipients)

	return nil
}

func (s *Service) enqueueBroadcastMessages(campaign *models.BroadcastCampaign, recipients []models.BroadcastRecipient) {
	baseDelay := 3
	if campaign.DelaySeconds > 0 {
		baseDelay = campaign.DelaySeconds
	}

	currentScheduledTime := time.Now()
	if campaign.ScheduledAt != nil && campaign.ScheduledAt.After(time.Now()) {
		currentScheduledTime = *campaign.ScheduledAt
	}

	broadcastID := campaign.ID.String()
	deviceID := campaign.DeviceID.String()

	for i, r := range recipients {
		// Add jitter to delay (0-2 seconds extra)
		jitter := rand.Intn(2000)
		delay := time.Duration(baseDelay)*time.Second + time.Duration(jitter)*time.Millisecond
		
		if i > 0 {
			currentScheduledTime = currentScheduledTime.Add(delay)
		}

		// Enqueue via message service
		content := ""
		if campaign.MessageContent != nil {
			content = *campaign.MessageContent
		}

		// Use the internal queueing mechanism by calling SendScheduledMessage
		// This ensures analytics and everything are tracked.
		_, err := s.messageService.SendScheduledMessage(
			context.Background(),
			deviceID,
			r.PhoneNumber,
			content,
			currentScheduledTime,
			nil, // mediaUrl
			"text",
			nil, // caption
			&broadcastID,
		)

		if err != nil {
			utils.Error("Failed to enqueue broadcast message", 
				zap.String("campaign_id", broadcastID),
				zap.String("recipient", r.PhoneNumber),
				zap.Error(err),
			)
			s.store.UpdateRecipientStatus(r.ID, "failed", utils.StringPtr(err.Error()))
		} else {
			// Link the queued message to the broadcast in DB
			// Actually, SendScheduledMessage creates a record in 'messages'.
			// We need to inject the BroadcastID. Currently SendScheduledMessage doesn't take BroadcastID.
			// I'll need to update MessageService to allow passing BroadcastID or do a manually Enqueue.
			
			// For now, I'll update the record manually right after.
			// TODO: Add BroadcastID param to SendScheduledMessage for better efficiency.
		}
	}
}

func (s *Service) handleMessageStatusUpdate(update *message.MessageStatusUpdate) {
	// Logic to update campaign counts when a message with BroadcastID is updated
	// Needs database support to find campaign by message
}

func (s *Service) ListCampaigns(userID uuid.UUID) ([]models.BroadcastCampaign, error) {
	return s.store.ListCampaigns(userID, 50, 0)
}

func (s *Service) GetCampaign(id uuid.UUID) (*models.BroadcastCampaign, error) {
	return s.store.GetCampaign(id)
}
