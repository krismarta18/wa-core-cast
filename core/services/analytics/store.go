package analytics

import (
	"time"

	"github.com/google/uuid"
	"wacast/core/database"
	"wacast/core/models"
)

type Store struct {
	db *database.Database
}

func NewStore(db *database.Database) *Store {
	return &Store{db: db}
}

// GetDailyStats retrieves message stats for the last N days
func (s *Store) GetDailyStats(userID uuid.UUID, days int) ([]models.DailyMessageStats, error) {
	return s.db.GetDailyStats(userID, days)
}

// GetDeviceStats aggregates stats per device for a user
func (s *Store) GetDeviceStats(userID uuid.UUID) ([]models.DailyMessageStats, error) {
	return s.db.GetDeviceUsageStats(userID)
}

// GetLatestFailures retrieves the most recent failure records
func (s *Store) GetLatestFailures(userID uuid.UUID, limit int) ([]models.FailureRecord, error) {
	return s.db.GetLatestFailureRecords(userID, limit)
}

// GetFailureReasonStats counts failures by type
func (s *Store) GetFailureReasonStats(userID uuid.UUID, days int) ([]struct {
	Type  string `json:"reason"`
	Count int    `json:"count"`
}, error) {
	raw, err := s.db.GetFailureTypeSummary(userID, days)
	if err != nil {
		return nil, err
	}

	results := make([]struct {
		Type  string `json:"reason"`
		Count int    `json:"count"`
	}, len(raw))

	for i, r := range raw {
		results[i].Type = r.Type
		results[i].Count = r.Count
	}

	return results, nil
}

func (s *Store) RecordSuccess(userID, deviceID uuid.UUID) error {
	return s.db.UpdateDailyStat(userID, deviceID, time.Now(), 1, 0, 0, 0)
}

func (s *Store) RecordFailure(record *models.FailureRecord) error {
	// Increment failure count in daily stats
	if err := s.db.UpdateDailyStat(record.UserID, record.DeviceID, record.OccurredAt, 1, 1, 0, 0); err != nil {
		return err
	}
	// Save detailed failure log
	return s.db.RecordFailure(record)
}

func (s *Store) RecordDelivery(userID, deviceID uuid.UUID) error {
	return s.db.UpdateDailyStat(userID, deviceID, time.Now(), 0, 0, 1, 0)
}

func (s *Store) RecordIncoming(userID, deviceID uuid.UUID) error {
	return s.db.UpdateDailyStat(userID, deviceID, time.Now(), 0, 0, 0, 1)
}

// GetFailedMessagesFromQueue queries the messages table for failed messages
func (s *Store) GetFailedMessagesFromQueue(userID uuid.UUID, limit int) ([]FailureLogItem, error) {
	raw, err := s.db.GetFailedMessagesAsLogs(userID, limit)
	if err != nil {
		return nil, err
	}
	items := make([]FailureLogItem, len(raw))
	for i, r := range raw {
		items[i] = FailureLogItem{
			ID:     r.ID,
			To:     r.To,
			Device: r.Device,
			Reason: r.Reason,
			Time:   r.Time,
			Type:   r.Type,
		}
	}
	return items, nil
}

// GetMessageFailureStats computes total sent, total failed, avg retry from messages table
func (s *Store) GetMessageFailureStats(userID uuid.UUID, days int) (totalSent int, totalFailed int, avgRetrySeconds float64) {
	return s.db.GetMessageFailureStatsFromQueue(userID, days)
}
