package autoresponse

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

// -- Keywords --

func (s *Store) CreateKeyword(kw *Keyword) error {
	query := `
		INSERT INTO auto_response_keywords (id, user_id, device_id, keyword, response_text, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	if kw.ID == uuid.Nil {
		kw.ID = uuid.New()
	}
	now := time.Now()
	if kw.CreatedAt.IsZero() {
		kw.CreatedAt = now
	}
	kw.UpdatedAt = now

	_, err := s.db.Exec(query, kw.ID, kw.UserID, kw.DeviceID, kw.Keyword, kw.ResponseText, kw.IsActive, kw.CreatedAt, kw.UpdatedAt)
	return err
}

func (s *Store) GetKeywordsByUserID(userID uuid.UUID) ([]*Keyword, error) {
	query := `
		SELECT id, user_id, device_id, keyword, response_text, is_active, created_at, updated_at
		FROM auto_response_keywords
		WHERE user_id = $1
		ORDER BY created_at DESC
	`
	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keywords []*Keyword
	for rows.Next() {
		var k Keyword
		if err := rows.Scan(&k.ID, &k.UserID, &k.DeviceID, &k.Keyword, &k.ResponseText, &k.IsActive, &k.CreatedAt, &k.UpdatedAt); err != nil {
			return nil, err
		}
		keywords = append(keywords, &k)
	}
	return keywords, nil
}

func (s *Store) UpdateKeyword(kw *Keyword) error {
	query := `
		UPDATE auto_response_keywords
		SET keyword = $1, response_text = $2, is_active = $3, updated_at = $4
		WHERE id = $5 AND user_id = $6
	`
	kw.UpdatedAt = time.Now()
	_, err := s.db.Exec(query, kw.Keyword, kw.ResponseText, kw.IsActive, kw.UpdatedAt, kw.ID, kw.UserID)
	return err
}

func (s *Store) DeleteKeyword(id uuid.UUID, userID uuid.UUID) error {
	query := `DELETE FROM auto_response_keywords WHERE id = $1 AND user_id = $2`
	_, err := s.db.Exec(query, id, userID)
	return err
}

func (s *Store) GetKeywordByID(id uuid.UUID, userID uuid.UUID) (*Keyword, error) {
	query := `
		SELECT id, user_id, device_id, keyword, response_text, is_active, created_at, updated_at
		FROM auto_response_keywords
		WHERE id = $1 AND user_id = $2
	`
	var k Keyword
	err := s.db.QueryRow(query, id, userID).Scan(&k.ID, &k.UserID, &k.DeviceID, &k.Keyword, &k.ResponseText, &k.IsActive, &k.CreatedAt, &k.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &k, nil
}

// -- Templates --

func (s *Store) CreateTemplate(t *MessageTemplate) error {
	query := `
		INSERT INTO message_templates (id, user_id, name, category, content, used_count, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	now := time.Now()
	if t.CreatedAt.IsZero() {
		t.CreatedAt = now
	}
	t.UpdatedAt = now

	_, err := s.db.Exec(query, t.ID, t.UserID, t.Name, t.Category, t.Content, t.UsedCount, t.CreatedAt, t.UpdatedAt)
	return err
}

func (s *Store) GetTemplatesByUserID(userID uuid.UUID) ([]*MessageTemplate, error) {
	query := `
		SELECT id, user_id, name, category, content, used_count, created_at, updated_at
		FROM message_templates
		WHERE user_id = $1
		ORDER BY created_at DESC
	`
	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var templates []*MessageTemplate
	for rows.Next() {
		var t MessageTemplate
		if err := rows.Scan(&t.ID, &t.UserID, &t.Name, &t.Category, &t.Content, &t.UsedCount, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		templates = append(templates, &t)
	}
	return templates, nil
}

func (s *Store) UpdateTemplate(t *MessageTemplate) error {
	query := `
		UPDATE message_templates
		SET name = $1, category = $2, content = $3, updated_at = $4
		WHERE id = $5 AND user_id = $6
	`
	t.UpdatedAt = time.Now()
	_, err := s.db.Exec(query, t.Name, t.Category, t.Content, t.UpdatedAt, t.ID, t.UserID)
	return err
}

func (s *Store) DeleteTemplate(id uuid.UUID, userID uuid.UUID) error {
	query := `DELETE FROM message_templates WHERE id = $1 AND user_id = $2`
	_, err := s.db.Exec(query, id, userID)
	return err
}

func (s *Store) GetTemplateByID(id uuid.UUID, userID uuid.UUID) (*MessageTemplate, error) {
	query := `
		SELECT id, user_id, name, category, content, used_count, created_at, updated_at
		FROM message_templates
		WHERE id = $1 AND user_id = $2
	`
	var t MessageTemplate
	err := s.db.QueryRow(query, id, userID).Scan(&t.ID, &t.UserID, &t.Name, &t.Category, &t.Content, &t.UsedCount, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
