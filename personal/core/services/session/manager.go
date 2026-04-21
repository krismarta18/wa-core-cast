package session

import (
	"context"
	"sync"
	"time"

	"wacast/core/utils"

	"go.uber.org/zap"
)

// CleanupInactiveSessionsInterval is the interval to check for inactive sessions
const CleanupInactiveSessionsInterval = 5 * time.Minute

// Manager handles session lifecycle and auto-reconnection
type Manager struct {
	service           *Service
	cleanupTicker     *time.Ticker
	done              chan struct{}
	autoReconnectEnabled bool
	reconnectInterval time.Duration
}

// NewManager creates a new session manager
func NewManager(service *Service, autoReconnectEnabled bool, reconnectInterval time.Duration) *Manager {
	return &Manager{
		service:              service,
		done:                 make(chan struct{}),
		autoReconnectEnabled: autoReconnectEnabled,
		reconnectInterval:    reconnectInterval,
	}
}

// Start begins the manager's background operations
func (m *Manager) Start() {
	m.cleanupTicker = time.NewTicker(CleanupInactiveSessionsInterval)

	go func() {
		for {
			select {
			case <-m.cleanupTicker.C:
				m.cleanupInactiveSessions()
				if m.autoReconnectEnabled {
					m.attemptReconnects()
				}
			case <-m.done:
				return
			}
		}
	}()

	utils.Info("Session manager started",
		zap.Bool("auto_reconnect", m.autoReconnectEnabled),
		zap.Duration("reconnect_interval", m.reconnectInterval),
	)
}

// Stop gracefully stops the manager
func (m *Manager) Stop() {
	close(m.done)
	if m.cleanupTicker != nil {
		m.cleanupTicker.Stop()
	}
	utils.Info("Session manager stopped")
}

// cleanupInactiveSessions removes sessions that have been inactive too long
func (m *Manager) cleanupInactiveSessions() {
	sessions := m.service.GetAllActiveSessions()
	now := time.Now().Unix()

	for _, session := range sessions {
		if m.service.sessionTimeout > 0 {
			inactiveFor := now - session.LastActivity
			if inactiveFor > int64(m.service.sessionTimeout) {
				utils.Warn("Closing inactive session",
					zap.String("device_id", session.ID),
					zap.Duration("inactive_for", time.Duration(inactiveFor)*time.Second),
				)

				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				if err := m.service.StopSession(ctx, session.ID); err != nil {
					utils.Error("Failed to stop inactive session",
						zap.String("device_id", session.ID),
						zap.Error(err),
					)
				}
				cancel()
			}
		}
	}
}

// attemptReconnects tries to reconnect sessions that are disconnected
// but should be active
func (m *Manager) attemptReconnects() {
	// This would be called periodically to:
	// 1. Query database for devices that should be active
	// 2. Check if session exists and is connected
	// 3. If not, try to restore from session_data
	// 4. Log any reconnection attempts
	
	// TODO: Implement full reconnection logic
}

// SessionHealthCheck periodically verifies session health
type SessionHealthCheck struct {
	mu               sync.RWMutex
	lastCheckTime    map[string]time.Time
	healthCheckInterval time.Duration
}

// NewSessionHealthCheck creates a new health checker
func NewSessionHealthCheck(interval time.Duration) *SessionHealthCheck {
	return &SessionHealthCheck{
		lastCheckTime:       make(map[string]time.Time),
		healthCheckInterval: interval,
	}
}

// Check performs a health check on a session
func (shc *SessionHealthCheck) Check(deviceID string, session *WhatsAppSession) (bool, error) {
	shc.mu.Lock()
	lastCheck, exists := shc.lastCheckTime[deviceID]
	shc.mu.Unlock()

	// Skip if recently checked
	if exists && time.Since(lastCheck) < shc.healthCheckInterval {
		return true, nil
	}

	// Verify connection is still alive
	// In whatsmeow, we can check if client is connected
	isConnected := session.Client != nil && session.Client.IsConnected()

	shc.mu.Lock()
	shc.lastCheckTime[deviceID] = time.Now()
	shc.mu.Unlock()

	if !isConnected {
		utils.Warn("Session health check failed - session not connected",
			zap.String("device_id", deviceID),
		)
		return false, nil
	}

	return true, nil
}

// GetHealthStatus returns health status for all sessions
func (shc *SessionHealthCheck) GetHealthStatus(service *Service) map[string]bool {
	status := make(map[string]bool)
	
	for _, session := range service.GetAllActiveSessions() {
		ok, _ := shc.Check(session.ID, session)
		status[session.ID] = ok
	}

	return status
}

// SessionMetrics holds statistics about sessions
type SessionMetrics struct {
	mu                  sync.RWMutex
	TotalSessionsCreated int
	TotalSessionsEnded   int
	CurrentActiveSessions int
	FailedConnections   int
	SuccessfulReconnects int
}

// RecordSessionCreated records a new session creation
func (sm *SessionMetrics) RecordSessionCreated() {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.TotalSessionsCreated++
	sm.CurrentActiveSessions++
}

// RecordSessionEnded records a session termination
func (sm *SessionMetrics) RecordSessionEnded() {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.TotalSessionsEnded++
	if sm.CurrentActiveSessions > 0 {
		sm.CurrentActiveSessions--
	}
}

// RecordFailedConnection records a failed connection attempt
func (sm *SessionMetrics) RecordFailedConnection() {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.FailedConnections++
}

// RecordSuccessfulReconnect records a successful reconnection
func (sm *SessionMetrics) RecordSuccessfulReconnect() {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.SuccessfulReconnects++
}

// GetMetrics returns current metrics
func (sm *SessionMetrics) GetMetrics() map[string]interface{} {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	return map[string]interface{}{
		"total_sessions_created":   sm.TotalSessionsCreated,
		"total_sessions_ended":     sm.TotalSessionsEnded,
		"current_active_sessions": sm.CurrentActiveSessions,
		"failed_connections":       sm.FailedConnections,
		"successful_reconnects":    sm.SuccessfulReconnects,
	}
}
