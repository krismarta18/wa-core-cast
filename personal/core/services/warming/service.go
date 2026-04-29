package warming

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"wacast/core/database"
	"wacast/core/services/message"
	"wacast/core/services/session"
	"wacast/core/utils"
)

// Service handles WhatsApp account warming logic
type Service struct {
	db             *database.Database
	sessionService *session.Service
	messageService *message.Service
	mu             sync.RWMutex
	activeSessions map[string]*WarmingSession
	messagePool    []string
}

// WarmingSession represents an ongoing warming session for a user
type WarmingSession struct {
	UserID    uuid.UUID
	DeviceIDs []uuid.UUID
	StartTime time.Time
	EndTime   time.Time
	Done      chan struct{}
}

// NewService creates a new warming service
func NewService(db *database.Database, sessionService *session.Service, messageService *message.Service) *Service {
	return &Service{
		db:             db,
		sessionService: sessionService,
		messageService: messageService,
		activeSessions: make(map[string]*WarmingSession),
		messagePool: []string{
			"Halo, apa kabar?",
			"Iya, nanti saya kabari lagi ya.",
			"Oke, siap bos!",
			"Lagi dimana sekarang?",
			"Mantap, lanjut!",
			"P",
			"Lagi sibuk nggak?",
			"Aman ya?",
			"Sip, terima kasih infonya.",
			"Boleh minta tolong?",
			"Jangan lupa ya hari ini.",
			"Sudah makan belum?",
			"Nanti kita bahas lagi.",
			"Ok gampang itu.",
			"Siap, laksanakan!",
			"Test koneksi bentar.",
			"Aman bro.",
			"Waduh, kok gitu?",
			"Bisa kirim filenya?",
			"Otw nih.",
		},
	}
}

// StartWarming starts a warming session for the given devices
func (s *Service) StartWarming(ctx context.Context, userID uuid.UUID, deviceIDs []uuid.UUID, durationMinutes int) error {
	if len(deviceIDs) < 2 {
		return fmt.Errorf("minimal butuh 2 nomor untuk fitur warming")
	}

	// 1. Check if all devices belong to user and are connected
	for _, dID := range deviceIDs {
		dev, err := s.db.GetDeviceByID(dID)
		if err != nil || dev == nil {
			return fmt.Errorf("device %s tidak ditemukan", dID)
		}
		if dev.UserID != userID {
			return fmt.Errorf("device %s bukan milik Anda", dID)
		}
		if !s.sessionService.IsSessionActive(dID.String()) {
			return fmt.Errorf("device %s sedang tidak aktif/terkoneksi", dev.PhoneNumber)
		}
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// 2. Check if user already has an active session
	if _, exists := s.activeSessions[userID.String()]; exists {
		return fmt.Errorf("sesi warming Anda sedang berjalan")
	}

	endTime := time.Now().Add(time.Duration(durationMinutes) * time.Minute)

	// 3. Update database status (Lockdown)
	for _, dID := range deviceIDs {
		if err := s.db.UpdateDeviceWarmingStatus(dID, true, &endTime); err != nil {
			utils.Error("Failed to update device warming status", zap.Error(err), zap.String("device_id", dID.String()))
		}
	}

	// 4. Create session
	sess := &WarmingSession{
		UserID:    userID,
		DeviceIDs: deviceIDs,
		StartTime: time.Now(),
		EndTime:   endTime,
		Done:      make(chan struct{}),
	}
	s.activeSessions[userID.String()] = sess

	// 5. Start the worker
	go s.runWarmingWorker(sess)

	utils.Info("Warming session started", zap.String("user_id", userID.String()), zap.Int("duration", durationMinutes))
	return nil
}

// StopWarming stops an active warming session
func (s *Service) StopWarming(userID uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	sess, exists := s.activeSessions[userID.String()]
	if !exists {
		return fmt.Errorf("tidak ada sesi warming yang aktif")
	}

	close(sess.Done)
	delete(s.activeSessions, userID.String())

	// Unlock devices in DB
	for _, dID := range sess.DeviceIDs {
		_ = s.db.UpdateDeviceWarmingStatus(dID, false, nil)
	}

	utils.Info("Warming session stopped manually", zap.String("user_id", userID.String()))
	return nil
}

func (s *Service) runWarmingWorker(sess *WarmingSession) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	// Initial small delay
	time.Sleep(2 * time.Second)

	for {
		select {
		case <-sess.Done:
			return
		case <-ticker.C:
			if time.Now().After(sess.EndTime) {
				s.StopWarming(sess.UserID)
				return
			}

			// Random chance to send a message (e.g., 80% chance)
			if rand.Float32() > 0.8 {
				continue
			}

			// Pick two random different devices from the session
			idx1 := rand.Intn(len(sess.DeviceIDs))
			idx2 := rand.Intn(len(sess.DeviceIDs))
			for idx1 == idx2 {
				idx2 = rand.Intn(len(sess.DeviceIDs))
			}

			devA := sess.DeviceIDs[idx1]
			devB := sess.DeviceIDs[idx2]

			go s.performPingPong(devA, devB, sess.Done)
		}
	}
}

func (s *Service) performPingPong(fromID, toID uuid.UUID, done chan struct{}) {
	// 1. Get target phone number
	toDev, err := s.db.GetDeviceByID(toID)
	if err != nil {
		return
	}

	// 2. Simulate Typing on Sender
	s.sessionService.SetChatPresence(fromID.String(), toDev.PhoneNumber, "composing")
	time.Sleep(time.Duration(rand.Intn(4)+2) * time.Second)

	// 3. Send Message
	content := s.messagePool[rand.Intn(len(s.messagePool))]
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err = s.messageService.SendInternalMessage(ctx, fromID.String(), toDev.PhoneNumber, content)
	if err != nil {
		utils.Error("Warming send failed", zap.Error(err), zap.String("from", fromID.String()), zap.String("to", toDev.PhoneNumber))
		return
	}

	// 4. Wait for receiver to "read" (Simulate reading time)
	time.Sleep(time.Duration(rand.Intn(10)+5) * time.Second)
	
	// Mark as read is handled by the receiver's session usually, but we can force it
	// Actually, whatsmeow handles mark read if enabled. 
	// In our session service, we have global message callbacks.
}

// GetStatus returns the current status of warming for a user
func (s *Service) GetStatus(userID uuid.UUID) map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sess, exists := s.activeSessions[userID.String()]
	if !exists {
		return map[string]interface{}{
			"is_active": false,
		}
	}

	remainingSecs := int(time.Until(sess.EndTime).Seconds())
	if remainingSecs < 0 {
		remainingSecs = 0
	}

	totalSecs := int(sess.EndTime.Sub(sess.StartTime).Seconds())

	return map[string]interface{}{
		"is_active":              true,
		"remaining_seconds":      remainingSecs,
		"total_duration_seconds": totalSecs,
		"device_count":           len(sess.DeviceIDs),
		"start_time":             sess.StartTime,
	}
}
