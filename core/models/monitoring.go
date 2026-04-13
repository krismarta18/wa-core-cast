package models

import (
	"time"

	"github.com/google/uuid"
)

// DailyMessageStats holds aggregated daily message statistics per device
type DailyMessageStats struct {
	ID             uuid.UUID `json:"id"`
	UserID         uuid.UUID `json:"user_id"`
	DeviceID       uuid.UUID `json:"device_id"`
	StatDate       time.Time `json:"stat_date"`
	TotalSent      int       `json:"total_sent"`
	TotalReceived  int       `json:"total_received"`
	TotalFailed    int       `json:"total_failed"`
	TotalDelivered int       `json:"total_delivered"`
	TotalRead      int       `json:"total_read"`
	CreatedAt      time.Time `json:"created_at"`
}

// TableName returns the table name
func (DailyMessageStats) TableName() string {
	return "daily_message_stats"
}

// FailureRecord logs a message delivery failure for analysis
type FailureRecord struct {
	ID           uuid.UUID `json:"id"`
	MessageID    uuid.UUID `json:"message_id"`
	DeviceID     uuid.UUID `json:"device_id"`
	FailureType  string    `json:"failure_type"` // network/timeout/invalid_number/etc
	ErrorMessage string    `json:"error_message"`
	RetryCount   int       `json:"retry_count"`
	CreatedAt    time.Time `json:"created_at"`
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
