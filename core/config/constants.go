package config

const (
	// Service Info
	ServiceName    = "wacast-core"
	ServiceVersion = "1.0.0"

	// Default Values
	DefaultServerPort        = 8080
	DefaultServerHost        = "0.0.0.0"
	DefaultLogLevel          = "debug"
	DefaultConnectionPool    = 25
	DefaultConnectionMaxAge  = 5 // minutes
	DefaultConnectionTimeout = 5 // seconds

	// WhatsApp
	DefaultSessionTimeout = 300 // seconds

	// Message Status
	MessageStatusPending   = 0
	MessageStatusSent      = 1
	MessageStatusDelivered = 2
	MessageStatusRead      = 3
	MessageStatusFailed    = 4

	// Device Status
	DeviceStatusActive     = 1
	DeviceStatusInactive   = 0
	DeviceStatusDisconnect = 2

	// Message Direction
	MessageDirectionIn  = "IN"
	MessageDirectionOut = "OUT"

	// Pagination
	DefaultPageSize = 20
	MaxPageSize     = 100
)
