package integration

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"wacast/core/database"
	"wacast/core/models"
)

type Service struct {
	db *database.Database
}

func NewService(db *database.Database) *Service {
	return &Service{db: db}
}

// ─── API Keys ────────────────────────────────────────────────────────────────

func (s *Service) ListAPIKeys(userID uuid.UUID) ([]models.APIKey, error) {
	rows, err := s.db.Query("SELECT id, user_id, name, prefix, last_used_at, created_at FROM api_keys WHERE user_id = $1 AND deleted_at IS NULL", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keys []models.APIKey
	for rows.Next() {
		var k models.APIKey
		if err := rows.Scan(&k.ID, &k.UserID, &k.Name, &k.KeyPrefix, &k.LastUsedAt, &k.CreatedAt); err != nil {
			return nil, err
		}
		keys = append(keys, k)
	}
	return keys, nil
}

func (s *Service) CreateAPIKey(userID uuid.UUID, name string) (*models.APIKeyResponse, error) {
	rawKey := fmt.Sprintf("wck_live_%s", uuid.New().String())
	hash := s.hashKey(rawKey)
	id := uuid.New()
	prefix := rawKey[0:12]
	now := time.Now()
	
	_, err := s.db.Exec("INSERT INTO api_keys (id, user_id, name, key_hash, prefix, created_at) VALUES ($1, $2, $3, $4, $5, $6)", 
		id, userID, name, hash, prefix, now)
	if err != nil {
		return nil, err
	}

	return &models.APIKeyResponse{
		ID:        id,
		Name:      name,
		Key:       rawKey,
		Prefix:    prefix,
		CreatedAt: now,
	}, nil
}

func (s *Service) DeleteAPIKey(userID uuid.UUID, id uuid.UUID) error {
	_, err := s.db.Exec("UPDATE api_keys SET deleted_at = $1 WHERE id = $2 AND user_id = $3", time.Now(), id, userID)
	return err
}

func (s *Service) ValidateAPIKey(key string) (*uuid.UUID, error) {
	hash := s.hashKey(key)
	var userID uuid.UUID
	var id uuid.UUID
	
	err := s.db.QueryRow("SELECT id, user_id FROM api_keys WHERE key_hash = $1 AND deleted_at IS NULL", hash).Scan(&id, &userID)
	if err != nil {
		return nil, err
	}

	// Update last used
	go s.db.Exec("UPDATE api_keys SET last_used_at = $1 WHERE id = $2", time.Now(), id)

	return &userID, nil
}

func (s *Service) hashKey(key string) string {
	h := sha256.New()
	h.Write([]byte(key))
	return hex.EncodeToString(h.Sum(nil))
}

// ─── Webhooks ────────────────────────────────────────────────────────────────

func (s *Service) GetWebhookSettings(userID uuid.UUID) (*models.WebhookSettings, error) {
	var settings models.WebhookSettings
	var enabledEvents []byte
	
	err := s.db.QueryRow("SELECT user_id, url, secret, is_active, enabled_events, updated_at FROM webhook_settings WHERE user_id = $1", userID).
		Scan(&settings.UserID, &settings.URL, &settings.Secret, &settings.IsActive, &enabledEvents, &settings.UpdatedAt)
	
	if err != nil {
		// Return default if not found
		return &models.WebhookSettings{
			UserID:        userID,
			IsActive:      false,
			EnabledEvents: json.RawMessage("[]"),
		}, nil
	}
	settings.EnabledEvents = json.RawMessage(enabledEvents)
	return &settings, nil
}

func (s *Service) UpdateWebhookSettings(userID uuid.UUID, settings *models.WebhookSettings) error {
	_, err := s.db.Exec(`
		INSERT INTO webhook_settings (user_id, url, secret, is_active, enabled_events, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (user_id) DO UPDATE SET 
			url = EXCLUDED.url, 
			secret = EXCLUDED.secret, 
			is_active = EXCLUDED.is_active, 
			enabled_events = EXCLUDED.enabled_events, 
			updated_at = EXCLUDED.updated_at`,
		userID, settings.URL, settings.Secret, settings.IsActive, settings.EnabledEvents, time.Now())
	return err
}

func (s *Service) TriggerWebhook(userID uuid.UUID, event string, payload interface{}) {
	settings, err := s.GetWebhookSettings(userID)
	if err != nil || settings.URL == "" || !settings.IsActive {
		return
	}

	// Check if event is enabled
	var enabled []string
	json.Unmarshal(settings.EnabledEvents, &enabled)
	isEnabled := false
	for _, e := range enabled {
		if e == event {
			isEnabled = true
			break
		}
	}
	if !isEnabled {
		return
	}

	go s.deliverWebhook(settings.URL, settings.Secret, event, payload)
}

func (s *Service) deliverWebhook(url, secret, event string, payload interface{}) {
	body, _ := json.Marshal(map[string]interface{}{
		"event":     event,
		"timestamp": time.Now().Unix(),
		"data":      payload,
	})

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "WACAST-Webhook/1.0")

	if secret != "" {
		signer := hmac.New(sha256.New, []byte(secret))
		signer.Write(body)
		signature := hex.EncodeToString(signer.Sum(nil))
		req.Header.Set("X-WACAST-Signature", signature)
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err == nil {
		defer resp.Body.Close()
	}
}
