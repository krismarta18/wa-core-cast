package models

import (
	"time"

	"github.com/google/uuid"
)

// DailyMessageStats holds aggregated daily message statistics per device
type DailyMessageStats struct {
	ID             uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID         uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	DeviceID       uuid.UUID `json:"device_id" gorm:"type:uuid"`
	StatDate       time.Time `json:"stat_date" gorm:"type:date;not null"`
	SentCount      int       `json:"sent_count" gorm:"default:0"`
	FailedCount    int       `json:"failed_count" gorm:"default:0"`
	DeliveredCount int       `json:"delivered_count" gorm:"default:0"`
	ReceivedCount  int       `json:"received_count" gorm:"default:0"`
	SuccessRate    float64   `json:"success_rate" gorm:"type:numeric(5,2)"`
	CreatedAt      time.Time `json:"created_at" gorm:"default:now()"`
}

// TableName returns the table name
func (DailyMessageStats) TableName() string {
	return "daily_message_stats"
}

// FailureRecord logs a message delivery failure for analysis
type FailureRecord struct {
	ID             uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID         uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	DeviceID       uuid.UUID `json:"device_id" gorm:"type:uuid"`
	MessageID      uuid.UUID `json:"message_id" gorm:"type:uuid"`
	RecipientPhone string    `json:"recipient_phone" gorm:"type:varchar(30);not null"`
	FailureType    string    `json:"failure_type" gorm:"type:varchar(50);not null"` // send_failed, timeout, invalid_number, banned
	FailureReason  string    `json:"failure_reason" gorm:"type:text"`
	OccurredAt     time.Time `json:"occurred_at" gorm:"default:now()"`
	CreatedAt      time.Time `json:"created_at" gorm:"default:now()"`
}

// TableName returns the table name
func (FailureRecord) TableName() string {
	return "failure_records"
}

// ServiceHealthCheck records periodic health check results
type ServiceHealthCheck struct {
	ID           uuid.UUID  `json:"id"`
	ServiceName  string     `json:"service_name"`
	Status       string     `json:"status"` // healthy/degraded/unhealthy
	ResponseTime *int       `json:"response_time_ms,omitempty"`
	Message      *string    `json:"message,omitempty"`
	StartedAt    *time.Time `json:"started_at,omitempty"`
	CheckedAt    time.Time  `json:"checked_at"`
}

// TableName returns the table name
func (ServiceHealthCheck) TableName() string {
	return "service_health_checks"
}

// ResourceUsageMetric holds system resource usage snapshots
type ResourceUsageMetric struct {
	ID          uuid.UUID `json:"id"`
	MetricType  string    `json:"metric_type"` // cpu/memory/disk/connections
	Value       float64   `json:"value"`
	Unit        string    `json:"unit"` // percent/MB/GB/count
	RecordedAt  time.Time `json:"recorded_at"`
}

// TableName returns the table name
func (ResourceUsageMetric) TableName() string {
	return "resource_usage_metrics"
}
