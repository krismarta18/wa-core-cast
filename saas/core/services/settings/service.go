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
	rows, err := db.QueryContext(ctx, "SELECT key, value FROM system_settings")
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
	_, err := db.ExecContext(ctx, "UPDATE system_settings SET value = $1, updated_at = NOW() WHERE key = $2", value, key)
	return err
}
