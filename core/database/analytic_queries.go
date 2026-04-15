package database

import (
	"time"

	"github.com/google/uuid"
	"wacast/core/models"
)

// GetDailyStats retrieves usage statistics for a user over a period of time
func (d *Database) GetDailyStats(userID uuid.UUID, days int) ([]models.DailyMessageStats, error) {
	query := `
		SELECT id, user_id, device_id, stat_date, sent_count, failed_count, 
		       delivered_count, received_count, success_rate, created_at
		FROM daily_message_stats
		WHERE user_id = $1 AND stat_date >= $2
		ORDER BY stat_date ASC
	`
	startDate := time.Now().AddDate(0, 0, -days)
	rows, err := d.Query(query, userID, startDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []models.DailyMessageStats
	for rows.Next() {
		var s models.DailyMessageStats
		err := rows.Scan(
			&s.ID, &s.UserID, &s.DeviceID, &s.StatDate, &s.SentCount, &s.FailedCount,
			&s.DeliveredCount, &s.ReceivedCount, &s.SuccessRate, &s.CreatedAt,
		)
		if err != nil {
			continue
		}
		stats = append(stats, s)
	}
	return stats, nil
}

// GetDeviceUsageStats aggregates usage stats per device for a user
func (d *Database) GetDeviceUsageStats(userID uuid.UUID) ([]models.DailyMessageStats, error) {
	query := `
		SELECT device_id, SUM(sent_count) as sent_count, AVG(success_rate) as success_rate
		FROM daily_message_stats
		WHERE user_id = $1
		GROUP BY device_id
	`
	rows, err := d.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []models.DailyMessageStats
	for rows.Next() {
		var s models.DailyMessageStats
		err := rows.Scan(&s.DeviceID, &s.SentCount, &s.SuccessRate)
		if err != nil {
			continue
		}
		stats = append(stats, s)
	}
	return stats, nil
}

// GetLatestFailureRecords retrieves recent message failures for a user
func (d *Database) GetLatestFailureRecords(userID uuid.UUID, limit int) ([]models.FailureRecord, error) {
	query := `
		SELECT id, user_id, device_id, message_id, recipient_phone, 
		       failure_type, failure_reason, occurred_at, created_at
		FROM failure_records
		WHERE user_id = $1
		ORDER BY occurred_at DESC
		LIMIT $2
	`
	rows, err := d.Query(query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []models.FailureRecord
	for rows.Next() {
		var r models.FailureRecord
		err := rows.Scan(
			&r.ID, &r.UserID, &r.DeviceID, &r.MessageID, &r.RecipientPhone,
			&r.FailureType, &r.FailureReason, &r.OccurredAt, &r.CreatedAt,
		)
		if err != nil {
			continue
		}
		records = append(records, r)
	}
	return records, nil
}

// GetFailureTypeSummary counts failures by type for a user
func (d *Database) GetFailureTypeSummary(userID uuid.UUID, days int) ([]struct {
	Type  string
	Count int
}, error) {
	query := `
		SELECT reason, COUNT(*) as count FROM (
			SELECT failure_type as reason FROM failure_records WHERE user_id = $1 AND occurred_at >= $2
			UNION ALL
			SELECT 'send_failed' as reason FROM messages m JOIN devices dev ON dev.id = m.device_id WHERE dev.user_id = $1 AND m.status_message = 4 AND m.created_at >= $2
			UNION ALL
			SELECT 'init_failed' as reason FROM scheduled_messages WHERE user_id = $1 AND status = 'failed' AND updated_at >= $2
			UNION ALL
			SELECT 'broadcast_failed' as reason FROM broadcast_recipients br JOIN broadcast_campaigns bc ON bc.id = br.campaign_id WHERE bc.user_id = $1 AND br.status = 'failed' AND br.failed_at >= $2
		) combined
		GROUP BY reason
		ORDER BY count DESC
	`
	startDate := time.Now().AddDate(0, 0, -days)
	rows, err := d.Query(query, userID, startDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []struct {
		Type  string
		Count int
	}
	for rows.Next() {
		var r struct {
			Type  string
			Count int
		}
		err := rows.Scan(&r.Type, &r.Count)
		if err != nil {
			continue
		}
		results = append(results, r)
	}
	return results, nil
}

// UpdateDailyStat increments a specific counter in the daily_message_stats table
func (d *Database) UpdateDailyStat(userID, deviceID uuid.UUID, date time.Time, sent, failed, delivered, received int) error {
	// Normalize date to 00:00:00
	statDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	query := `
		INSERT INTO daily_message_stats (
			id, user_id, device_id, stat_date, 
			sent_count, failed_count, delivered_count, received_count, 
			success_rate, created_at
		) VALUES (
			$1, $2, $3, $4, 
			$5, $6, $7, $8, 
			CASE WHEN $5 > 0 THEN (($5 - $6)::float / $5::float) * 100 ELSE 100 END, 
			NOW()
		)
		ON CONFLICT (user_id, device_id, stat_date) DO UPDATE SET
			sent_count = daily_message_stats.sent_count + EXCLUDED.sent_count,
			failed_count = daily_message_stats.failed_count + EXCLUDED.failed_count,
			delivered_count = daily_message_stats.delivered_count + EXCLUDED.delivered_count,
			received_count = daily_message_stats.received_count + EXCLUDED.received_count,
			success_rate = CASE 
				WHEN (daily_message_stats.sent_count + EXCLUDED.sent_count) > 0 
				THEN ((daily_message_stats.sent_count + EXCLUDED.sent_count - daily_message_stats.failed_count - EXCLUDED.failed_count)::float / (daily_message_stats.sent_count + EXCLUDED.sent_count)::float) * 100
				ELSE 100
			END,
			created_at = daily_message_stats.created_at
	`
	_, err := d.Exec(query, uuid.New(), userID, deviceID, statDate, sent, failed, delivered, received)
	return err
}

// RecordFailure inserts a new message failure record into the database
func (d *Database) RecordFailure(record *models.FailureRecord) error {
	query := `
		INSERT INTO failure_records (
			id, user_id, device_id, message_id, recipient_phone, 
			failure_type, failure_reason, occurred_at, created_at
		) VALUES (
			$1, $2, $3, $4, $5, 
			$6, $7, $8, NOW()
		)
	`
	if record.ID == uuid.Nil {
		record.ID = uuid.New()
	}
	if record.OccurredAt.IsZero() {
		record.OccurredAt = time.Now()
	}

	_, err := d.Exec(query, 
		record.ID, record.UserID, record.DeviceID, record.MessageID, record.RecipientPhone,
		record.FailureType, record.FailureReason, record.OccurredAt,
	)
	return err
}

// GetFailedMessagesAsLogs returns failed messages from the messages table as log items
// This captures failures from ALL sources: direct, scheduled, broadcast
func (d *Database) GetFailedMessagesAsLogs(userID uuid.UUID, limit int) ([]struct {
	ID     uuid.UUID
	To     string
	Device string
	Reason string
	Time   string
	Type   string
}, error) {
	query := `
		-- Direct and queued messages
		(SELECT m.id, m.receipt_number as target, m.device_id as device, 
		       COALESCE(m.error_log, 'Unknown error') as reason, m.created_at, 'message' as src_type
		FROM messages m
		JOIN devices dev ON dev.id = m.device_id
		WHERE dev.user_id = $1 
		  AND m.status_message = 4
		  AND m.created_at >= NOW() - INTERVAL '7 days')
		
		UNION ALL

		-- Scheduled messages failures (that didn't reach messages table or are marked failed)
		(SELECT sm.id, sm.recipient_payload->>'phone' as target, sm.device_id::text as device,
		       'Scheduled dispatch failed' as reason, sm.updated_at as created_at, 'scheduled' as src_type
		FROM scheduled_messages sm
		WHERE sm.user_id = $1
		  AND sm.status = 'failed'
		  AND sm.updated_at >= NOW() - INTERVAL '7 days')

		UNION ALL

		-- Broadcast recipient failures
		(SELECT br.id, br.phone_number as target, camp.device_id::text as device,
		       COALESCE(br.error_message, 'Broadcast dispatch failed') as reason, br.failed_at as created_at, 'broadcast' as src_type
		FROM broadcast_recipients br
		JOIN broadcast_campaigns camp ON camp.id = br.campaign_id
		WHERE camp.user_id = $1
		  AND br.status = 'failed'
		  AND br.failed_at >= NOW() - INTERVAL '7 days')

		ORDER BY created_at DESC
		LIMIT $2
	`
	rows, err := d.Query(query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []struct {
		ID     uuid.UUID
		To     string
		Device string
		Reason string
		Time   string
		Type   string
	}
	for rows.Next() {
		var id uuid.UUID
		var to, deviceID, reason, srcType string
		var createdAt time.Time
		if err := rows.Scan(&id, &to, &deviceID, &reason, &createdAt, &srcType); err != nil {
			continue
		}
		deviceLabel := "Device " + deviceID
		if len(deviceID) > 8 {
			deviceLabel = "Device " + deviceID[:8]
		}
		results = append(results, struct {
			ID     uuid.UUID
			To     string
			Device string
			Reason string
			Time   string
			Type   string
		}{
			ID:     id,
			To:     to,
			Device: deviceLabel,
			Reason: reason,
			Time:   createdAt.Format("02 Jan 2006, 15:04"),
			Type:   srcType,
		})
	}
	return results, nil
}

// GetMessageFailureStatsFromQueue returns total sent, total failed, and avg retry time
// by querying the messages table directly (covers direct, scheduled, broadcast)
func (d *Database) GetMessageFailureStatsFromQueue(userID uuid.UUID, days int) (totalSent int, totalFailed int, avgRetrySeconds float64) {
	startDate := time.Now().AddDate(0, 0, -days)

	// Total sent attempts across all tables
	sentQuery := `
		SELECT 
			(SELECT COUNT(*) FROM messages m JOIN devices dev ON dev.id = m.device_id WHERE dev.user_id = $1 AND m.direction = 'OUT' AND m.created_at >= $2) +
			(SELECT COUNT(*) FROM scheduled_messages WHERE user_id = $1 AND created_at >= $2 AND status IN ('sent', 'failed')) +
			(SELECT COUNT(*) FROM broadcast_recipients br JOIN broadcast_campaigns bc ON bc.id = br.campaign_id WHERE bc.user_id = $1 AND br.created_at >= $2 AND br.status IN ('sent', 'failed'))
	`
	_ = d.QueryRow(sentQuery, userID, startDate).Scan(&totalSent)

	// Total failed attempts
	failedQuery := `
		SELECT 
			(SELECT COUNT(*) FROM messages m JOIN devices dev ON dev.id = m.device_id WHERE dev.user_id = $1 AND m.status_message = 4 AND m.created_at >= $2) +
			(SELECT COUNT(*) FROM scheduled_messages WHERE user_id = $1 AND created_at >= $2 AND status = 'failed') +
			(SELECT COUNT(*) FROM broadcast_recipients br JOIN broadcast_campaigns bc ON bc.id = br.campaign_id WHERE bc.user_id = $1 AND br.failed_at >= $2 AND br.status = 'failed')
	`
	_ = d.QueryRow(failedQuery, userID, startDate).Scan(&totalFailed)

	// Avg retry: currently mostly tracked in messages table
	retryQuery := `
		SELECT COALESCE(AVG(retry_count * 5.0), 0) FROM messages m
		JOIN devices dev ON dev.id = m.device_id
		WHERE dev.user_id = $1 AND m.retry_count > 0 AND m.created_at >= $2
	`
	_ = d.QueryRow(retryQuery, userID, startDate).Scan(&avgRetrySeconds)

	return
}
