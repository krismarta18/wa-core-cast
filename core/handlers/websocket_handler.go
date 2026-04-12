package handlers

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"

	"wacast/core/services/session"
	"wacast/core/utils"
)

// WebSocketHandler handles real-time WebSocket connections
type WebSocketHandler struct {
	sessionService *session.Service
	clients        map[string]map[*websocket.Conn]bool // deviceID -> connections
	broadcast      chan *QRCodeUpdate
	mu             sync.RWMutex
	upgrader       websocket.Upgrader
}

// QRCodeUpdate represents a QR code update event
type QRCodeUpdate struct {
	DeviceID  string `json:"device_id"`
	QRCode    string `json:"qr_code"`
	Status    int    `json:"status"`
	Timestamp int64  `json:"timestamp"`
	ImageURL  string `json:"image_url"`
}

// NewWebSocketHandler creates a new WebSocket handler
func NewWebSocketHandler(svc *session.Service) *WebSocketHandler {
	handler := &WebSocketHandler{
		sessionService: svc,
		clients:        make(map[string]map[*websocket.Conn]bool),
		broadcast:      make(chan *QRCodeUpdate, 100),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// Allow all origins for development
				// In production, restrict to your domain
				return true
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	}

	// Start broadcast goroutine
	go handler.broadcastLoop()

	return handler
}

// ConnectQR connects a client to real-time QR code updates via WebSocket
// GET /ws/sessions/:device_id/qr
func (h *WebSocketHandler) ConnectQR(c *gin.Context) {
	deviceID := c.Param("device_id")

	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		utils.Error("Failed to upgrade WebSocket connection",
			zap.String("device_id", deviceID),
			zap.Error(err),
		)
		return
	}

	// Send initial state BEFORE registering as broadcast target
	// to avoid concurrent write race with broadcastLoop.
	qrString := h.sessionService.GetQRCode(deviceID)
	status := h.sessionService.GetSessionStatus(deviceID)
	if qrString != "" {
		conn.WriteJSON(&QRCodeUpdate{
			DeviceID:  deviceID,
			QRCode:    qrString,
			Status:    int(status),
			Timestamp: time.Now().UnixMilli(),
			ImageURL:  fmt.Sprintf("/api/v1/sessions/%s/qr?format=png", deviceID),
		})
	} else {
		conn.WriteJSON(map[string]interface{}{
			"device_id": deviceID,
			"message":   "Waiting for QR code from WhatsApp...",
			"status":    int(status),
		})
	}

	// Register for broadcast updates
	h.mu.Lock()
	if h.clients[deviceID] == nil {
		h.clients[deviceID] = make(map[*websocket.Conn]bool)
	}
	h.clients[deviceID][conn] = true
	h.mu.Unlock()

	utils.Info("WebSocket client connected", zap.String("device_id", deviceID))

	// Deregister and close when the read loop exits
	defer func() {
		h.mu.Lock()
		if h.clients[deviceID] != nil {
			delete(h.clients[deviceID], conn)
		}
		h.mu.Unlock()
		conn.Close()
		utils.Info("WebSocket client disconnected", zap.String("device_id", deviceID))
	}()

	// Read loop — blocks the goroutine, keeping the connection alive.
	// All writes happen exclusively in broadcastLoop to avoid concurrent write panics.
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				utils.Warn("WebSocket error", zap.String("device_id", deviceID), zap.Error(err))
			}
			return
		}
	}
}

// NotifyQRUpdate sends QR code update to all connected clients for a device
func (h *WebSocketHandler) NotifyQRUpdate(deviceID, qrCode string, status int) {
	update := &QRCodeUpdate{
		DeviceID:  deviceID,
		QRCode:    qrCode,
		Status:    status,
		Timestamp: time.Now().UnixMilli(),
		ImageURL:  fmt.Sprintf("/api/v1/sessions/%s/qr?format=png", deviceID),
	}

	select {
	case h.broadcast <- update:
	default:
		utils.Warn("Broadcast channel full, discarding QR update",
			zap.String("device_id", deviceID),
		)
	}
}

// broadcastLoop broadcasts updates to all connected clients
func (h *WebSocketHandler) broadcastLoop() {
	for update := range h.broadcast {
		h.mu.RLock()
		clients := h.clients[update.DeviceID]
		h.mu.RUnlock()

		if len(clients) == 0 {
			continue
		}

		// Send update to all clients for this device
		h.mu.RLock()
		for client := range clients {
			err := client.WriteJSON(update)
			if err != nil {
				utils.Warn("Failed to send WebSocket message",
					zap.String("device_id", update.DeviceID),
					zap.Error(err),
				)
				client.Close()
				delete(clients, client)
			}
		}
		h.mu.RUnlock()
	}
}

// ServeQRPage serves an HTML page with real-time QR code display and auto-refresh countdown.
// GET /qr/:device_id
func (h *WebSocketHandler) ServeQRPage(c *gin.Context) {
	deviceID := c.Param("device_id")

	html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>WhatsApp QR Code - %s</title>
    <style>
        *{margin:0;padding:0;box-sizing:border-box}
        body{font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,sans-serif;background:linear-gradient(135deg,#128C7E 0%%,#075E54 100%%);min-height:100vh;display:flex;align-items:center;justify-content:center;padding:20px}
        .card{background:#fff;border-radius:20px;box-shadow:0 20px 60px rgba(0,0,0,.3);padding:36px 32px;max-width:420px;width:100%%;text-align:center}
        h1{color:#075E54;font-size:22px;margin-bottom:4px}
        .sub{color:#888;font-size:13px;margin-bottom:20px}
        .badge{display:inline-block;padding:5px 14px;border-radius:20px;font-size:11px;font-weight:700;letter-spacing:.5px;text-transform:uppercase;margin-bottom:18px}
        .badge.pending{background:#fff8e1;color:#f57f17}
        .badge.active{background:#e8f5e9;color:#2e7d32}
        .badge.waiting{background:#f3f3f3;color:#888}
        .qr-wrap{position:relative;background:#fafafa;border-radius:14px;padding:20px;border:2px dashed #ddd;margin-bottom:16px;min-height:200px;display:flex;align-items:center;justify-content:center}
        #qrImage{max-width:100%%;border-radius:8px;display:none;transition:opacity .3s}
        #qrImage.expired{opacity:.25;filter:blur(2px)}
        .overlay{position:absolute;inset:0;display:none;flex-direction:column;align-items:center;justify-content:center;border-radius:14px;background:rgba(255,255,255,.85)}
        .overlay.show{display:flex}
        .spinner{width:28px;height:28px;border:3px solid #eee;border-top:3px solid #128C7E;border-radius:50%%;animation:spin 1s linear infinite;margin-bottom:8px}
        @keyframes spin{to{transform:rotate(360deg)}}
        .overlay p{color:#555;font-size:13px;font-weight:600}
        #loading{color:#aaa;font-size:13px}
        .countdown-bar{height:4px;background:#eee;border-radius:2px;margin-bottom:16px;overflow:hidden}
        .countdown-fill{height:100%%;background:#128C7E;border-radius:2px;transition:width 1s linear}
        .countdown-fill.urgent{background:#e53935}
        .countdown-text{font-size:12px;color:#999;margin-bottom:16px}
        .countdown-text span{font-weight:700}
        .steps{background:#e7f8f5;border-left:3px solid #128C7E;padding:12px 14px;border-radius:6px;text-align:left;font-size:12px;color:#1a5c52;line-height:1.8;margin-bottom:16px}
        .steps b{display:block;margin-bottom:4px;font-size:13px}
        .ws-dot{display:inline-block;width:8px;height:8px;border-radius:50%%;background:#ccc;margin-right:6px;transition:background .3s}
        .ws-dot.ok{background:#4caf50;animation:pulse 2s infinite}
        @keyframes pulse{0%%,100%%{opacity:1}50%%{opacity:.4}}
        .ws-label{font-size:11px;color:#aaa}
    </style>
</head>
<body>
<div class="card">
    <h1>📱 WhatsApp Link</h1>
    <p class="sub">Device: <code>%s</code></p>

    <div class="badge waiting" id="badge">Menghubungkan...</div>

    <div class="qr-wrap">
        <p id="loading">Menunggu QR code...</p>
        <img id="qrImage" alt="QR Code WhatsApp">
        <div class="overlay" id="overlay">
            <div class="spinner"></div>
            <p id="overlayMsg">Memperbarui QR...</p>
        </div>
    </div>

    <!-- countdown bar only shown when QR is visible -->
    <div id="countdownSection" style="display:none">
        <div class="countdown-bar"><div class="countdown-fill" id="countdownFill" style="width:100%%"></div></div>
        <p class="countdown-text">QR kedaluwarsa dalam <span id="countdownNum">20</span> detik</p>
    </div>

    <div class="steps">
        <b>Cara menghubungkan:</b>
        1. Buka WhatsApp di HP &rarr; <b>Perangkat Tertaut</b><br>
        2. Ketuk <b>"Tautkan Perangkat"</b><br>
        3. Arahkan kamera ke QR code di atas<br>
        4. Tunggu konfirmasi otomatis ✅
    </div>

    <div>
        <span class="ws-dot" id="wsDot"></span>
        <span class="ws-label" id="wsLabel">Menghubungkan WebSocket...</span>
    </div>
</div>

<script>
const deviceId = '%s';
const QR_TTL = 20; // WhatsApp QR expires ~20 seconds

const badge        = document.getElementById('badge');
const loadingEl    = document.getElementById('loading');
const qrImage      = document.getElementById('qrImage');
const overlay      = document.getElementById('overlay');
const overlayMsg   = document.getElementById('overlayMsg');
const countdownSec = document.getElementById('countdownSection');
const countdownFill= document.getElementById('countdownFill');
const countdownNum = document.getElementById('countdownNum');
const wsDot        = document.getElementById('wsDot');
const wsLabel      = document.getElementById('wsLabel');

let ws = null;
let countdownTimer = null;
let countdownVal = QR_TTL;
let reconnectAttempts = 0;
let fallbackTimer = null;

// ─── QR display ──────────────────────────────────────────────────
function showQR(imageUrl) {
    clearCountdown();
    const url = imageUrl + (imageUrl.includes('?') ? '&' : '?') + 't=' + Date.now();
    qrImage.src = url;
    qrImage.classList.remove('expired');
    qrImage.style.display = 'block';
    loadingEl.style.display = 'none';
    overlay.classList.remove('show');
    countdownSec.style.display = 'block';
    setBadge('pending', '⏳ Scan QR Code');
    startCountdown();
}

function showConnected() {
    clearCountdown();
    qrImage.style.display = 'none';
    loadingEl.style.display = 'none';
    countdownSec.style.display = 'none';
    overlay.classList.remove('show');
    setBadge('active', '✅ Terhubung!');
    loadingEl.style.display = 'block';
    loadingEl.textContent = '✅ WhatsApp berhasil terhubung!';
}

function showRefreshing() {
    qrImage.classList.add('expired');
    overlay.classList.add('show');
    overlayMsg.textContent = 'Memperbarui QR code...';
    countdownSec.style.display = 'none';
}

function setBadge(cls, text) {
    badge.className = 'badge ' + cls;
    badge.textContent = text;
}

// ─── Countdown ───────────────────────────────────────────────────
function startCountdown() {
    countdownVal = QR_TTL;
    updateCountdownUI();
    countdownTimer = setInterval(() => {
        countdownVal--;
        updateCountdownUI();
        if (countdownVal <= 0) {
            clearInterval(countdownTimer);
            // QR expired — show refreshing state and fallback-poll
            showRefreshing();
            fetchQRFallback();
        }
    }, 1000);
}

function clearCountdown() {
    if (countdownTimer) { clearInterval(countdownTimer); countdownTimer = null; }
}

function updateCountdownUI() {
    const pct = (countdownVal / QR_TTL) * 100;
    countdownFill.style.width = pct + '%%';
    countdownNum.textContent = countdownVal;
    if (countdownVal <= 5) {
        countdownFill.classList.add('urgent');
        countdownNum.style.color = '#e53935';
    } else {
        countdownFill.classList.remove('urgent');
        countdownNum.style.color = '';
    }
}

// ─── WebSocket ────────────────────────────────────────────────────
function connectWS() {
    if (ws && ws.readyState === WebSocket.OPEN) return;
    const proto = location.protocol === 'https:' ? 'wss:' : 'ws:';
    ws = new WebSocket(proto + '//' + location.host + '/ws/sessions/' + deviceId + '/qr');

    ws.onopen = () => {
        wsDot.className = 'ws-dot ok';
        wsLabel.textContent = 'Real-time terhubung';
        reconnectAttempts = 0;
        clearFallbackPoll();
    };

    ws.onmessage = (evt) => {
        try {
            const d = JSON.parse(evt.data);
            handleUpdate(d);
        } catch(e) { console.error('WS parse error', e); }
    };

    ws.onerror = () => {
        wsDot.className = 'ws-dot';
        wsLabel.textContent = 'WebSocket error';
    };

    ws.onclose = () => {
        wsDot.className = 'ws-dot';
        wsLabel.textContent = 'Terputus, mencoba ulang...';
        startFallbackPoll(); // poll REST while WS is down
        if (reconnectAttempts < 10) {
            reconnectAttempts++;
            setTimeout(connectWS, Math.min(3000 * reconnectAttempts, 15000));
        }
    };
}

function handleUpdate(d) {
    if (!d) return;
    // Connected success
    if (d.status === 1) { showConnected(); return; }
    // Has QR image URL
    if (d.image_url) {
        clearFallbackPoll();
        showQR(d.image_url);
        return;
    }
    // No QR yet
    if (d.message) {
        loadingEl.textContent = d.message;
    }
}

// ─── Fallback REST poll (when WS down or QR expired) ─────────────
function fetchQRFallback() {
    fetch('/api/v1/sessions/' + deviceId + '/qr')
        .then(r => r.json())
        .then(d => {
            if (d.qr_code_image && d.qr_code_image.direct_url) {
                showQR(d.qr_code_image.direct_url);
            }
        })
        .catch(() => {});
}

function startFallbackPoll() {
    clearFallbackPoll();
    fallbackTimer = setInterval(fetchQRFallback, 5000);
}

function clearFallbackPoll() {
    if (fallbackTimer) { clearInterval(fallbackTimer); fallbackTimer = null; }
}

// ─── Init ────────────────────────────────────────────────────────
window.addEventListener('load', () => {
    fetchQRFallback(); // immediate check
    connectWS();
});

window.addEventListener('focus', () => {
    if (!ws || ws.readyState !== WebSocket.OPEN) connectWS();
});
</script>
</body>
</html>
`, deviceID, deviceID, deviceID)

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(200, html)
}
