package autoresponse

import (
	"strings"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"wacast/core/utils"
)

type Service struct {
	store *Store
}

func NewService(store *Store) *Service {
	return &Service{store: store}
}

// -- Keywords --

func (s *Service) CreateKeyword(userID uuid.UUID, req *CreateKeywordRequest) (*Keyword, error) {
	kw := &Keyword{
		UserID:       userID,
		Keyword:      strings.ToLower(strings.TrimSpace(req.Keyword)),
		ResponseText: req.ResponseText,
		IsActive:     true,
	}

	if req.DeviceID != nil && *req.DeviceID != "" {
		did, err := uuid.Parse(*req.DeviceID)
		if err == nil {
			kw.DeviceID = &did
		}
	}

	if err := s.store.CreateKeyword(kw); err != nil {
		return nil, err
	}
	return kw, nil
}

func (s *Service) GetKeywords(userID uuid.UUID) ([]*Keyword, error) {
	return s.store.GetKeywordsByUserID(userID)
}

func (s *Service) UpdateKeyword(id uuid.UUID, userID uuid.UUID, req *UpdateKeywordRequest) (*Keyword, error) {
	kw, err := s.store.GetKeywordByID(id, userID)
	if err != nil {
		return nil, err
	}

	if req.Keyword != nil {
		kw.Keyword = strings.ToLower(strings.TrimSpace(*req.Keyword))
	}
	if req.ResponseText != nil {
		kw.ResponseText = *req.ResponseText
	}
	if req.IsActive != nil {
		kw.IsActive = *req.IsActive
	}

	if err := s.store.UpdateKeyword(kw); err != nil {
		return nil, err
	}
	return kw, nil
}

func (s *Service) DeleteKeyword(id uuid.UUID, userID uuid.UUID) error {
	return s.store.DeleteKeyword(id, userID)
}

func (s *Service) ToggleKeyword(id uuid.UUID, userID uuid.UUID) (*Keyword, error) {
	kw, err := s.store.GetKeywordByID(id, userID)
	if err != nil {
		return nil, err
	}

	kw.IsActive = !kw.IsActive
	if err := s.store.UpdateKeyword(kw); err != nil {
		return nil, err
	}
	return kw, nil
}

// -- Templates --

func (s *Service) CreateTemplate(userID uuid.UUID, req *CreateTemplateRequest) (*MessageTemplate, error) {
	tmpl := &MessageTemplate{
		UserID:   userID,
		Name:     req.Name,
		Category: req.Category,
		Content:  req.Content,
	}

	if err := s.store.CreateTemplate(tmpl); err != nil {
		return nil, err
	}
	return tmpl, nil
}

func (s *Service) GetTemplates(userID uuid.UUID) ([]*MessageTemplate, error) {
	return s.store.GetTemplatesByUserID(userID)
}

func (s *Service) UpdateTemplate(id uuid.UUID, userID uuid.UUID, req *UpdateTemplateRequest) (*MessageTemplate, error) {
	tmpl, err := s.store.GetTemplateByID(id, userID)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		tmpl.Name = *req.Name
	}
	if req.Category != nil {
		tmpl.Category = *req.Category
	}
	if req.Content != nil {
		tmpl.Content = *req.Content
	}

	if err := s.store.UpdateTemplate(tmpl); err != nil {
		return nil, err
	}
	return tmpl, nil
}

func (s *Service) DeleteTemplate(id uuid.UUID, userID uuid.UUID) error {
	return s.store.DeleteTemplate(id, userID)
}

// -- Processing --

// DetectAndReply checks if a message exactly matches an active keyword for the given user/device
// Returns the response text if matched, or an empty string if no match
func (s *Service) DetectAndReply(userID uuid.UUID, deviceID *uuid.UUID, incomingMsg string) string {
	msgContent := strings.ToLower(strings.TrimSpace(incomingMsg))
	if msgContent == "" {
		return ""
	}

	utils.Debug("DetectAndReply details", 
		zap.String("user_id", userID.String()),
		zap.String("content", msgContent),
	)

	keywords, err := s.store.GetKeywordsByUserID(userID)
	if err != nil {
		utils.Error("Failed to fetch keywords from store", zap.Error(err))
		return ""
	}

	utils.Debug("Keywords found for user", zap.Int("count", len(keywords)), zap.String("user_id", userID.String()))

	for _, kw := range keywords {
		utils.Debug("Checking keyword", 
			zap.String("kw", kw.Keyword), 
			zap.Bool("active", kw.IsActive),
			zap.Any("kw_device_id", kw.DeviceID),
		)
		
		if !kw.IsActive {
			continue
		}
		
		// If keyword requires a specific device match
		if kw.DeviceID != nil && deviceID != nil {
			if *kw.DeviceID != *deviceID {
				utils.Debug("Device ID mismatch", zap.String("expected", kw.DeviceID.String()), zap.String("actual", deviceID.String()))
				continue
			}
		}

		// Support comma-separated keywords
		triggers := strings.Split(kw.Keyword, ",")
		matched := false
		
		for _, trigger := range triggers {
			cleanTrigger := strings.ToLower(strings.TrimSpace(trigger))
			if cleanTrigger == "" {
				continue
			}

			// Check if message contains the trigger
			if strings.Contains(msgContent, cleanTrigger) {
				utils.Info("Auto-reply matched",
					zap.String("trigger", cleanTrigger),
					zap.String("keyword_record", kw.Keyword),
					zap.String("message_content", msgContent),
				)
				matched = true
				break
			}
		}

		if matched {
			return kw.ResponseText
		}
	}

	return ""
}
