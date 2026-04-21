package settings

import (
	"context"
	"fmt"

	"wacast/core/database"
)

type Service struct {
	db *database.Database
}

func NewService(db *database.Database) *Service {
	return &Service{db: db}
}

func (s *Service) GetSettings(ctx context.Context) (map[string]string, error) {
	db := s.db.GetConnection()
	rows, err := db.QueryContext(ctx, "SELECT key, value FROM settings")
	if err != nil {
		return nil, fmt.Errorf("get-settings: %w", err)
	}
	defer rows.Close()

	settings := make(map[string]string)
	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err != nil {
			continue
		}
		settings[key] = value
	}

	return settings, nil
}

func (s *Service) UpdateSetting(ctx context.Context, key, value string) error {
	db := s.db.GetConnection()
	query := `
		INSERT INTO settings (key, value, updated_at) 
		VALUES ($1, $2, NOW()) 
		ON CONFLICT (key) DO UPDATE SET 
			value = EXCLUDED.value, 
			updated_at = EXCLUDED.updated_at
	`
	_, err := db.ExecContext(ctx, query, key, value)
	return err
}
