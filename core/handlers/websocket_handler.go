package handlers

import (
	"fmt"
	"net/http"
	"sync"

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

	// Upgrade connection to WebSocket
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		utils.Error("Failed to upgrade WebSocket connection",
			zap.String("device_id", deviceID),
			zap.Error(err),
		)
		return
	}
	defer conn.Close()

	// Register client
	h.mu.Lock()
	if h.clients[deviceID] == nil {
		h.clients[deviceID] = make(map[*websocket.Conn]bool)
	}
	h.clients[deviceID][conn] = true
	h.mu.Unlock()

	utils.Info("WebSocket client connected",
		zap.String("device_id", deviceID),
		zap.Int("total_clients", len(h.clients[deviceID])),
	)

	// Send initial QR code if it exists
	qrString := h.sessionService.GetQRCode(deviceID)
	status := h.sessionService.GetSessionStatus(deviceID)

	if qrString != "" {
		// QR code exists, send it immediately
		qrPreview := qrString
		if len(qrString) > 50 {
			qrPreview = qrString[:50] + "..."
		}

		update := &QRCodeUpdate{
			DeviceID:  deviceID,
			QRCode:    qrString,
			Status:    int(status),
			Timestamp: 0,
			ImageURL:  fmt.Sprintf("/api/v1/sessions/%s/qr?format=png", deviceID),
		}
		err := conn.WriteJSON(update)
		if err != nil {
			utils.Error("Failed to send initial QR update",
				zap.String("device_id", deviceID),
				zap.Error(err),
			)
		} else {
			utils.Info("Initial QR update sent",
				zap.String("device_id", deviceID),
				zap.String("qr_preview", qrPreview),
			)
		}
	} else {
		// No QR yet, inform client to wait
		conn.WriteJSON(map[string]interface{}{
			"device_id": deviceID,
			"message":   "Waiting for QR code from WhatsApp...",
			"status":    int(status),
		})
		utils.Info("No QR code yet, client waiting",
			zap.String("device_id", deviceID),
		)
	}

	// Read messages from client (keep connection alive)
	go func() {
		defer func() {
			// Deregister client
			h.mu.Lock()
			delete(h.clients[deviceID], conn)
			h.mu.Unlock()

			utils.Info("WebSocket client disconnected",
				zap.String("device_id", deviceID),
				zap.Int("total_clients", len(h.clients[deviceID])),
			)
		}()

		for {
			var msg map[string]interface{}
			err := conn.ReadJSON(&msg)
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					utils.Warn("WebSocket error",
						zap.String("device_id", deviceID),
						zap.Error(err),
					)
				}
				return
			}

			// Handle ping/pong for keep-alive
			if msgType, ok := msg["type"].(string); ok && msgType == "ping" {
				conn.WriteJSON(map[string]string{"type": "pong"})
			}
		}
	}()
}

// NotifyQRUpdate sends QR code update to all connected clients for a device
func (h *WebSocketHandler) NotifyQRUpdate(deviceID, qrCode string, status int) {
	update := &QRCodeUpdate{
		DeviceID:  deviceID,
		QRCode:    qrCode,
		Status:    status,
		Timestamp: 0,
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

// ServeQRPage serves an HTML page with real-time QR code display
// GET /qr/:device_id
func (h *WebSocketHandler) ServeQRPage(c *gin.Context) {
	deviceID := c.Param("device_id")

	html := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>WhatsApp QR Code - %s</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
            background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
            min-height: 100vh;
            display: flex;
            align-items: center;
            justify-content: center;
            padding: 20px;
        }

        .container {
            background: white;
            border-radius: 20px;
            box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
            padding: 40px;
            max-width: 500px;
            text-align: center;
        }

        .header {
            margin-bottom: 30px;
        }

        .header h1 {
            color: #333;
            margin-bottom: 10px;
            font-size: 24px;
        }

        .header p {
            color: #666;
            font-size: 14px;
        }

        .status {
            display: inline-block;
            padding: 8px 16px;
            border-radius: 20px;
            font-size: 12px;
            font-weight: 600;
            margin-bottom: 20px;
            text-transform: uppercase;
            letter-spacing: 0.5px;
        }

        .status.pending {
            background: #fff3cd;
            color: #856404;
        }

        .status.active {
            background: #d4edda;
            color: #155724;
        }

        .qr-container {
            position: relative;
            margin: 30px 0;
            background: #f8f9fa;
            border-radius: 15px;
            padding: 25px;
            border: 2px dashed #ddd;
        }

        #qrImage {
            max-width: 100%%;
            height: auto;
            border-radius: 10px;
            transition: opacity 0.3s ease;
        }

        .instructions {
            background: #e7f3ff;
            border-left: 4px solid #2196F3;
            padding: 15px;
            border-radius: 5px;
            margin: 20px 0;
            text-align: left;
            font-size: 13px;
            color: #1565c0;
            line-height: 1.6;
        }

        .instructions ol {
            margin-left: 20px;
        }

        .instructions li {
            margin-bottom: 8px;
        }

        .info {
            background: #f5f5f5;
            padding: 15px;
            border-radius: 10px;
            margin: 20px 0;
            font-size: 12px;
            color: #666;
        }

        .connection-status {
            display: inline-flex;
            align-items: center;
            gap: 8px;
            font-size: 12px;
            margin-top: 20px;
            padding: 10px;
            background: #f0f0f0;
            border-radius: 20px;
        }

        .status-dot {
            width: 8px;
            height: 8px;
            border-radius: 50%%;
            background: #ccc;
            transition: background 0.3s ease;
        }

        .status-dot.connected {
            background: #4caf50;
            animation: pulse 2s infinite;
        }

        @keyframes pulse {
            0%% { opacity: 1; }
            50%% { opacity: 0.5; }
            100%% { opacity: 1; }
        }

        .loading {
            text-align: center;
            color: #999;
            margin: 20px 0;
        }

        .spinner {
            display: inline-block;
            width: 20px;
            height: 20px;
            border: 3px solid #f3f3f3;
            border-top: 3px solid #667eea;
            border-radius: 50%%;
            animation: spin 1s linear infinite;
            margin-right: 10px;
            vertical-align: middle;
        }

        @keyframes spin {
            0%% { transform: rotate(0deg); }
            100%% { transform: rotate(360deg); }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>📱 WhatsApp QR Code</h1>
            <p id="deviceId">Device: %s</p>
        </div>

        <div class="status pending" id="status">Connecting...</div>

        <div class="qr-container">
            <div class="loading" id="loading">
                <span class="spinner"></span>
                Waiting for QR code... (Initialize session first via API)
            </div>
            <img id="qrImage" style="display: none;" alt="WhatsApp QR Code">
        </div>

        <div class="instructions">
            <strong>How to use:</strong>
            <ol>
                <li><strong>Step 1:</strong> Initialize session via API: <code style="background: #f0f0f0; padding: 2px 6px; border-radius: 3px;">POST /api/v1/sessions/initiate</code></li>
                <li><strong>Step 2:</strong> Open this page with the <code style="background: #f0f0f0; padding: 2px 6px; border-radius: 3px;">device_id</code></li>
                <li><strong>Step 3:</strong> QR code will appear below once WhatsApp generates it</li>
                <li><strong>Step 4:</strong> Open WhatsApp on your phone → Settings > Linked Devices → Link a Device</li>
                <li><strong>Step 5:</strong> Point your camera at the QR code below</li>
                <li><strong>Step 6:</strong> Wait for confirmation on your phone</li>
            </ol>
        </div>

        <div class="info">
            ℹ️ This page connects to real-time updates. Once a session is initialized, the QR code will appear automatically when WhatsApp generates it.
        </div>

        <div class="connection-status">
            <span class="status-dot" id="connectionDot"></span>
            <span id="connectionText">Connecting...</span>
        </div>
    </div>

    <script>
        const deviceId = '%s';
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const wsUrl = protocol + '//' + window.location.host + '/ws/sessions/' + deviceId + '/qr';

        console.log('WebSocket URL:', wsUrl);
        console.log('Device ID:', deviceId);

        let ws = null;
        let reconnectAttempts = 0;
        const maxReconnectAttempts = 5;
        const reconnectDelay = 3000;

        function connectWebSocket() {
            try {
                ws = new WebSocket(wsUrl);

                ws.onopen = function() {
                    console.log('WebSocket connected');
                    document.getElementById('connectionDot').classList.add('connected');
                    document.getElementById('connectionText').textContent = 'Real-time connected';
                    reconnectAttempts = 0;

                    // Send ping every 30 seconds to keep connection alive
                    setInterval(function() {
                        if (ws && ws.readyState === WebSocket.OPEN) {
                            ws.send(JSON.stringify({ type: 'ping' }));
                        }
                    }, 30000);
                };

                ws.onmessage = function(event) {
                    try {
                        const data = JSON.parse(event.data);
                        updateQRCode(data);
                    } catch (e) {
                        console.error('Failed to parse WebSocket message:', e);
                    }
                };

                ws.onerror = function(error) {
                    console.error('WebSocket error:', error);
                    document.getElementById('connectionDot').classList.remove('connected');
                    document.getElementById('connectionText').textContent = 'Connection error';
                };

                ws.onclose = function() {
                    console.log('WebSocket disconnected');
                    document.getElementById('connectionDot').classList.remove('connected');
                    document.getElementById('connectionText').textContent = 'Disconnected';

                    // Attempt to reconnect
                    if (reconnectAttempts < maxReconnectAttempts) {
                        reconnectAttempts++;
                        console.log('Reconnecting in ' + (reconnectDelay / 1000) + 's...');
                        setTimeout(connectWebSocket, reconnectDelay);
                    }
                };
            } catch (e) {
                console.error('Failed to connect WebSocket:', e);
            }
        }

        function updateQRCode(data) {
            console.log('Received update:', JSON.stringify(data));

            if (!data) {
                console.warn('No data received');
                return;
            }

            // Check if this is a QR code update
            if (!data.image_url && !data.qr_code) {
                console.log('Waiting message received:', data.message);
                document.getElementById('loading').textContent = data.message || 'Waiting for QR code...';
                return;
            }

            const loading = document.getElementById('loading');
            const qrImage = document.getElementById('qrImage');
            const status = document.getElementById('status');

            if (!loading || !qrImage || !status) {
                console.error('Required DOM elements not found');
                return;
            }

            // Hide loading spinner
            loading.style.display = 'none';

            // Use image_url if provided, otherwise construct it
            const imageUrl = data.image_url || '/api/v1/sessions/' + deviceId + '/qr?format=png';
            const cacheBustingUrl = imageUrl + (imageUrl.includes('?') ? '&' : '?') + 't=' + Date.now();
            
            console.log('Loading QR image from:', cacheBustingUrl);
            
            qrImage.src = cacheBustingUrl;
            qrImage.style.display = 'block';
            
            qrImage.onload = function() {
                console.log('✅ QR image loaded successfully');
            };
            
            qrImage.onerror = function() {
                console.error('❌ Failed to load QR image from:', cacheBustingUrl);
                loading.style.display = 'block';
                loading.innerHTML = '<span class="spinner"></span> Failed to load QR image';
            };

            // Update status badge
            if (data.status === 1) {
                status.className = 'status active';
                status.textContent = '✅ Connected';
            } else if (data.status === 2) {
                status.className = 'status pending';
                status.textContent = '⏳ Pending - Scan QR Code';
            } else {
                status.className = 'status pending';
                status.textContent = 'Status: ' + data.status;
            }
            
            console.log('QR code updated successfully');
        }

        // Fetch QR from API as fallback
        function fetchQRFromAPI() {
            console.log('Fetching QR from API...');
            fetch('/api/v1/sessions/' + deviceId + '/qr')
                .then(response => response.json())
                .then(data => {
                    if (data.qr_code_string) {
                        console.log('✅ Got QR from API');
                        updateQRCode({
                            device_id: deviceId,
                            qr_code: data.qr_code_string,
                            status: data.status,
                            image_url: '/api/v1/sessions/' + deviceId + '/qr?format=png'
                        });
                    } else {
                        console.log('No QR available yet from API');
                    }
                })
                .catch(err => console.error('Failed to fetch QR from API:', err));
        }

        // Connect to WebSocket on page load
        window.addEventListener('load', function() {
            console.log('Page loaded, connecting to real-time QR updates...');
            
            // Try to fetch QR from API first (in case session already initialized)
            fetchQRFromAPI();
            
            // Connect to WebSocket for real-time updates
            connectWebSocket();
        });

        // Reconnect on page focus
        window.addEventListener('focus', function() {
            if (ws === null || ws.readyState === WebSocket.CLOSED) {
                console.log('Reconnecting WebSocket on focus...');
                connectWebSocket();
            }
        });
        
        // Also try to refresh QR every 5 seconds as fallback
        setInterval(function() {
            if (!document.getElementById('qrImage') || document.getElementById('qrImage').style.display === 'none') {
                fetchQRFromAPI();
            }
        }, 5000);
    </script>
</body>
</html>
`, deviceID, deviceID, deviceID)

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(200, html)
}
