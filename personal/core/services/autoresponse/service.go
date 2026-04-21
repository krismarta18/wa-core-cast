package autoresponse

import (
	"regexp"
	"strings"
	"time"

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

// DetectAndReply checks if a message matches an active keyword based on MatchType and optional schedule
// Delimiter '|' can be used in Keyword field to specify time range: "keyword|08:00-17:00"
func (s *Service) DetectAndReply(userID uuid.UUID, deviceID *uuid.UUID, incomingMsg string) string {
	msgContent := strings.ToLower(strings.TrimSpace(incomingMsg))
	if msgContent == "" {
		return ""
	}

	keywords, err := s.store.GetKeywordsByUserID(userID)
	if err != nil {
		utils.Error("Failed to fetch keywords from store", zap.Error(err))
		return ""
	}

	currentTime := time.Now().Format("15:04")

	for _, kw := range keywords {
		if !kw.IsActive {
			continue
		}

		// Device filter
		if kw.DeviceID != nil && deviceID != nil {
			if *kw.DeviceID != *deviceID {
				continue
			}
		}

		// Parse keyword and optional schedule: "keyword|08:00-17:00"
		rawKeyword := kw.Keyword
		var schedule string
		if parts := strings.Split(rawKeyword, "|"); len(parts) > 1 {
			rawKeyword = strings.TrimSpace(parts[0])
			schedule = strings.TrimSpace(parts[1])
		}

		// Time Schedule Check (Optional)
		if schedule != "" {
			timeParts := strings.Split(schedule, "-")
			if len(timeParts) == 2 {
				start := strings.TrimSpace(timeParts[0])
				end := strings.TrimSpace(timeParts[1])
				// Skip if current time is outside range
				if currentTime < start || currentTime > end {
					continue
				}
			}
		}

		// Support comma-separated triggers within the keyword part
		triggers := strings.Split(rawKeyword, ",")
		matched := false
		
		for _, trigger := range triggers {
			cleanTrigger := strings.ToLower(strings.TrimSpace(trigger))
			if cleanTrigger == "" {
				continue
			}

			// Flexible matching based on MatchType
			// MatchType is expected to be "exact", "contains", "regex", "starts_with", "ends_with"
			switch kw.MatchType {
			case "regex":
				if match, _ := regexp.MatchString(cleanTrigger, msgContent); match {
					matched = true
				}
			case "exact":
				if msgContent == cleanTrigger {
					matched = true
				}
			case "starts_with":
				if strings.HasPrefix(msgContent, cleanTrigger) {
					matched = true
				}
			case "ends_with":
				if strings.HasSuffix(msgContent, cleanTrigger) {
					matched = true
				}
			default: // Default to "contains" if not specified or "contains"
				if strings.Contains(msgContent, cleanTrigger) {
					matched = true
				}
			}

			if matched {
				utils.Info("Auto-reply matched",
					zap.String("match_type", kw.MatchType),
					zap.String("trigger", cleanTrigger),
					zap.String("message_content", msgContent),
				)
				break
			}
		}

		if matched {
			return kw.ResponseText
		}
	}

	return ""
}
