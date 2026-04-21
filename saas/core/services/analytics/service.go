package analytics

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"wacast/core/models"
)

type Service struct {
	store *Store
}

func NewService(store *Store) *Service {
	return &Service{store: store}
}

type UsageStatsResponse struct {
	TotalSent    int                        `json:"total_sent"`
	TotalFailed  int                        `json:"total_failed"`
	SuccessRate  float64                    `json:"success_rate"`
	Daily        []models.DailyMessageStats `json:"daily"`
	DeviceStats  []DeviceStat               `json:"device_stats"`
}

type DeviceStat struct {
	Name        string  `json:"name"`
	SentCount   int     `json:"sent"`
	SuccessRate float64 `json:"success"`
}

func (s *Service) GetUsageStats(ctx context.Context, userID uuid.UUID) (*UsageStatsResponse, error) {
	daily, err := s.store.GetDailyStats(userID, 7)
	if err != nil {
		return nil, err
	}
	if daily == nil {
		daily = make([]models.DailyMessageStats, 0)
	}

	// Use comprehensive stats for the summary cards to ensure sync with Failure Rate menu
	compSent, compFailed, _ := s.store.GetMessageFailureStats(userID, 7)
	
	totalSent := compSent
	totalFailed := compFailed
	
	// If for some reason compStats are 0 but daily has data, fallback to daily (safety)
	if totalSent == 0 && totalFailed == 0 {
		for _, d := range daily {
			totalSent += d.SentCount
			totalFailed += d.FailedCount
		}
	}

	successRate := 0.0
	if totalSent > 0 {
		successRate = float64(totalSent-totalFailed) / float64(totalSent) * 100
	}

	// For now, device stats will be simple mapping
	rawDeviceStats, _ := s.store.GetDeviceStats(userID)
	deviceStats := make([]DeviceStat, 0)
	for _, r := range rawDeviceStats {
		deviceStats = append(deviceStats, DeviceStat{
			Name:        "Device " + r.DeviceID.String()[:8], // Mocking name prefix
			SentCount:   r.SentCount,
			SuccessRate: r.SuccessRate,
		})
	}

	return &UsageStatsResponse{
		TotalSent:   totalSent,
		TotalFailed: totalFailed,
		SuccessRate: successRate,
		Daily:       daily,
		DeviceStats: deviceStats,
	}, nil
}

type FailureRateResponse struct {
	TotalFailed7Days int              `json:"total_failed"`
	FailureRate      float64          `json:"failure_rate"`
	AvgRetryTime     string           `json:"avg_retry_time"`
	ReasonStats      []ReasonStat     `json:"reason_stats"`
	LatestLogs       []FailureLogItem `json:"latest_logs"`
}

type ReasonStat struct {
	Reason string  `json:"reason"`
	Count  int     `json:"count"`
	Pct    float64 `json:"pct"`
}

type FailureLogItem struct {
	ID        uuid.UUID `json:"id"`
	To        string    `json:"to"`
	Device    string    `json:"device"`
	Reason    string    `json:"reason"`
	Time      string    `json:"time"`
	Type      string    `json:"type"`
}

func (s *Service) GetFailureAnalytics(ctx context.Context, userID uuid.UUID) (*FailureRateResponse, error) {
	// 1. Get comprehensive reason stats (from failure_records + messages + scheduled + broadcast)
	reasons, err := s.store.GetFailureReasonStats(userID, 7)
	if err != nil {
		return nil, err
	}

	totalFailed := 0
	reasonStats := make([]ReasonStat, 0)
	if len(reasons) > 0 {
		reasonStats = make([]ReasonStat, len(reasons))
		for i, r := range reasons {
			totalFailed += r.Count
			reasonStats[i] = ReasonStat{
				Reason: r.Type,
				Count:  r.Count,
			}
		}

		// Calculate percentages for reasons
		for i := range reasonStats {
			if totalFailed > 0 {
				reasonStats[i].Pct = float64(reasonStats[i].Count) / float64(totalFailed) * 100
			}
		}
	}

	// 2. Get comprehensive logs (UNION of all sources)
	logItems := make([]FailureLogItem, 0)
	dbLogs, err := s.store.GetFailedMessagesFromQueue(userID, 20)
	if err == nil && dbLogs != nil {
		logItems = dbLogs
	} else {
		// Fallback to basic failure records if needed
		logs, _ := s.store.GetLatestFailures(userID, 20)
		for _, l := range logs {
			logItems = append(logItems, FailureLogItem{
				ID:     l.ID,
				To:     l.RecipientPhone,
				Device: "Device " + l.DeviceID.String()[:8],
				Reason: l.FailureReason,
				Time:   l.OccurredAt.Format("02 Jan 2006, 15:04"),
				Type:   l.FailureType,
			})
		}
	}

	// 3. Compute real failure rate and totals using the new comprehensive stats method
	totalSentStats, totalFailedStats, avgRetry := s.store.GetMessageFailureStats(userID, 7)

	// Ensure we use the most accurate total failed count
	if totalFailedStats > totalFailed {
		totalFailed = totalFailedStats
	}

	failureRate := 0.0
	if totalSentStats > 0 {
		failureRate = float64(totalFailed) / float64(totalSentStats) * 100
	}

	avgRetryStr := "0s"
	if avgRetry > 0 {
		avgRetryStr = fmt.Sprintf("%.1fs", avgRetry)
	}

	return &FailureRateResponse{
		TotalFailed7Days: totalFailed,
		FailureRate:      failureRate,
		AvgRetryTime:     avgRetryStr,
		ReasonStats:      reasonStats,
		LatestLogs:       logItems,
	}, nil
}

func (s *Service) RecordSent(userID, deviceID uuid.UUID) error {
	return s.store.RecordSuccess(userID, deviceID)
}

func (s *Service) RecordFailure(userID, deviceID, messageID uuid.UUID, phone, failureType, reason string) error {
	record := &models.FailureRecord{
		UserID:         userID,
		DeviceID:       deviceID,
		MessageID:      messageID,
		RecipientPhone: phone,
		FailureType:    failureType,
		FailureReason:  reason,
	}
	return s.store.RecordFailure(record)
}

func (s *Service) RecordDelivery(userID, deviceID uuid.UUID) error {
	return s.store.RecordDelivery(userID, deviceID)
}

func (s *Service) RecordReceived(userID, deviceID uuid.UUID) error {
	return s.store.RecordIncoming(userID, deviceID)
}
