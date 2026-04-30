package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/inconshreveable/go-update"
	"go.uber.org/zap"
	"wacast/core/config"
	"wacast/core/database"
	"wacast/core/utils"
	"embed"
	"encoding/base64"

	"github.com/getlantern/systray"
)

//go:embed icon/favicon.png icon/tray_16.png
var iconFS embed.FS

// Valid Windows ICO file (Base64) - 16x16
const trayIcoBase64 = "AAABAAEAEBAAAAEAIADeAgAAFgAAAIlQTkcNChoKAAAADUlIRFIAAAAQAAAAEAgCAAAAkJFoNgAAAqVJREFUeJw0kktvG1UUx8+9d16ecZxm4sgU13amVRtIQwUiSGWHKlCWUFjBmiVb2PULwBfgC7BgjcSy4iEQIAFVaW2SoMRO/IwfOPbYnrmPc9BEQvov/ovfle4552cBA3sv9N7d5td9YEAIYIg0kSEymBVFlBjsL/QPPWzNmPXStbVHr/NqnjReBYRhoEkrjQpJIUkkSZgabC/U14eWe1DhN/IoMzpP3oPia/cKkcV4Zzb6tvlLc9ENhfNqpfT9cTMJXX5vS3gPI37dJ4UFyH0SvefY8vH0xyfxM8+GD248aE/61cCqOG7ecZvDS5hLixBQIUo8eGF/INsrSgLMHc+a59POILj48NZbT9rfaaNvrge/AouROJmMthXf8cs/XfzxfvjOXX5TLaVayafDhqRESEcwVGkaMA6GLNKEEoUSoPVgNnrc+tnjrq9dh9uj5XiexlyBTCQimFSTIZ79J8XlIlnO0wL6Xx19k5PuF3uffX730xKGtGLPu53+eJYkahanYECI+yXY8PRKbYuiEKbx78lv3T9bvbPO5MJ3Cx+98rA+7C1Wsn7amSw1TFOLFEJiuILoWrlChdu16l/Tv+dJHG3WQh5u8OD+1sustPuo/iUhI01WdpfEOCl/dtTo9weLdHlns1oTZT+2K5Wt0M/vvXjr6UmjALmJisEQEx/vQrTuErrMZDMx4HNrx4tGw9Eb0a7NhCvs1rB3eNk7vexDe25BYiA1uTzUIstI4kBVvBNOy6Wg+ObO/vPzf9ZzebCd+qQLhsAgY29X2H7JQQxsbTFCRWuiELnRIlkyAotZBs04np+O+0ppOhwzKHrsYJuFHst2jKAJTCYsIwK8kpcITVZovKLfuwwAYNNjtzdgzQGCjEbKQvB/rl7GKZxlLv0XAAD//yfxpzyhE0lhAAAAAElFTkSuQmCC"

type ControlPanel struct {
	app          *App
	launcherPort int
	guiCmd       *exec.Cmd
}

func NewControlPanel() *ControlPanel {
	return &ControlPanel{
		app:          NewApp(),
		launcherPort: 19991,
	}
}

func (cp *ControlPanel) Run() {
	systray.Run(cp.onReady, cp.onExit)
}

func (cp *ControlPanel) onReady() {
	// Set Tray Icon using proper ICO data
	iconData, _ := base64.StdEncoding.DecodeString(trayIcoBase64)
	systray.SetIcon(iconData)
	systray.SetTitle("WACAST")
	systray.SetTooltip("WACAST WhatsApp Gateway")

	// Menu items
	mOpenDash := systray.AddMenuItem("Open Dashboard", "Open the main web dashboard")
	mOpenCtrl := systray.AddMenuItem("Open Control Panel", "Open this control window")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Exit WACAST", "Shut down everything")

	// Handle menu clicks
	go func() {
		for {
			select {
			case <-mOpenDash.ClickedCh:
				cp.openDashboard()
			case <-mOpenCtrl.ClickedCh:
				cp.OpenUI()
			case <-mQuit.ClickedCh:
				cp.handleExit(nil, nil)
			}
		}
	}()

	// Start HTTP server for GUI in background
	mux := http.NewServeMux()
	mux.HandleFunc("/", cp.handleUI)
	mux.HandleFunc("/favicon.png", func(w http.ResponseWriter, r *http.Request) {
		data, err := iconFS.ReadFile("icon/favicon.png")
		if err != nil {
			http.Error(w, "Icon not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "image/png")
		w.Write(data)
	})
	mux.HandleFunc("/api/status", cp.handleStatus)
	mux.HandleFunc("/api/start", cp.handleStart)
	mux.HandleFunc("/api/stop", cp.handleStop)
	mux.HandleFunc("/api/config", cp.handleConfig)
	mux.HandleFunc("/api/test-db", cp.handleTestDB)
	mux.HandleFunc("/api/open-dashboard", cp.handleOpenDashboard)
	mux.HandleFunc("/api/logs", cp.handleLogs)
	mux.HandleFunc("/api/exit", cp.handleExit)
	mux.HandleFunc("/api/check-update", cp.handleCheckUpdate)
	mux.HandleFunc("/api/do-update", cp.handleDoUpdate)

	server := &http.Server{
		Addr:    fmt.Sprintf("127.0.0.1:%d", cp.launcherPort),
		Handler: mux,
	}

	go func() {
		utils.Info("Starting GUI server...", zap.String("addr", server.Addr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			utils.Error("GUI Server failed to start!", zap.Error(err))
			os.Exit(1)
		}
	}()

	// Wait a bit for server to bind
	time.Sleep(500 * time.Millisecond)

	utils.Info("WACAST Control Panel is active.")
	
	// Launch UI initially
	go cp.OpenUI()
}

func (cp *ControlPanel) onExit() {
	// Clean up here if needed
}

func (cp *ControlPanel) OpenUI() {
	url := fmt.Sprintf("http://127.0.0.1:%d", cp.launcherPort)
	if runtime.GOOS == "windows" {
		edgePath := "msedge"
		// Try to find the actual path if not in PATH
		if _, err := exec.LookPath(edgePath); err != nil {
			paths := []string{
				`C:\Program Files (x86)\Microsoft\Edge\Application\msedge.exe`,
				`C:\Program Files\Microsoft\Edge\Application\msedge.exe`,
			}
			for _, p := range paths {
				if _, err := os.Stat(p); err == nil {
					edgePath = p
					break
				}
			}
		}
		// Use 'start /wait "" "path"' syntax
		// First, check if the edgePath is actually a valid file
		useFallback := false
		if _, err := os.Stat(edgePath); err != nil && edgePath != "msedge" {
			utils.Warn("Edge not found at path, using fallback", zap.String("path", edgePath))
			useFallback = true
		}

		if !useFallback {
			// Run browser directly with unique profile to allow tracking and force-close
			userDataDir := os.Getenv("TEMP") + "\\wacast_gui_profile"
			cp.guiCmd = exec.Command(edgePath, "--app="+url, "--window-size=380,640", "--user-data-dir="+userDataDir, "--no-first-run")
			_ = cp.guiCmd.Start()
			return
		}
	} else {
		cp.guiCmd = exec.Command("google-chrome", "--app="+url, "--window-size=380,640")
		_ = cp.guiCmd.Start()
		return
	}
	
	// Fallback if no specific browser found
	cp.triggerFallback(url)
}

func (cp *ControlPanel) triggerFallback(url string) {
	// Fallback: Open in default browser (won't block, so we add a manual block)
	fallbackCmd := exec.Command("cmd", "/c", "start", url)
	_ = fallbackCmd.Run()
	
	fmt.Printf("\nWAJIB DIBACA: Gagal membuka jendela khusus. Membuka di browser default...\n")
	fmt.Printf("Aplikasi akan terus berjalan di latar belakang.\n")
	fmt.Printf("TUTUP TERMINAL INI (Ctrl+C) UNTUK MEMATIKAN APLIKASI.\n")
	
	// Manual block since 'start url' doesn't block
	select {}
}


func (cp *ControlPanel) openDashboard() {
	cfg, _ := config.LoadConfig()
	url := fmt.Sprintf("http://localhost:%d", cfg.ServerPort)
	_ = utils.OpenBrowser(url)
}

// --- HTTP Handlers ---

func (cp *ControlPanel) handleUI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, `
	<!DOCTYPE html>
	<html>
	<head>
		<title>WACAST Control Panel</title>
		<link rel="icon" type="image/png" href="/favicon.png">
		<style>
			:root {
				--bg: #0f111a;
				--card: #1a1d2e;
				--accent: #3b82f6;
				--text: #e2e8f0;
				--text-dim: #94a3b8;
				--success: #10b981;
				--danger: #ef4444;
			}
			body { 
				font-family: 'Segoe UI', system-ui, sans-serif; 
				margin: 0; background: var(--bg); color: var(--text); 
				overflow: hidden; user-select: none;
			}
			.header { 
				background: var(--card); padding: 14px 20px; 
				display: flex; justify-content: space-between; align-items: center;
				border-bottom: 1px solid #2d3748;
			}
			.header h1 { margin: 0; font-size: 14px; font-weight: 700; color: var(--text); letter-spacing: 1.5px; }
			
			.status-indicator { display: flex; align-items: center; gap: 8px; font-size: 11px; font-weight: 700; }
			.dot { width: 10px; height: 10px; border-radius: 50%; background: var(--danger); box-shadow: 0 0 10px var(--danger); transition: 0.3s; }
			.dot.active { background: var(--success); box-shadow: 0 0 10px var(--success); }

			.tabs { display: flex; background: var(--card); padding: 0 12px; gap: 8px; border-bottom: 1px solid #2d3748; }
			.tab { 
				padding: 10px 16px; cursor: pointer; font-size: 11px; font-weight: 600; 
				color: var(--text-dim); border-bottom: 2px solid transparent; transition: 0.2s;
			}
			.tab:hover { color: var(--text); }
			.tab.active { color: var(--accent); border-bottom-color: var(--accent); }
			
			.content { height: calc(100vh - 90px); overflow-y: auto; padding: 16px; box-sizing: border-box; }
			
			.control-row { display: grid; grid-template-columns: 1fr 1fr 1fr; gap: 10px; margin-bottom: 16px; }
			
			.btn { 
				padding: 10px; border-radius: 8px; border: none; cursor: pointer; 
				font-weight: 700; font-size: 11px; transition: 0.2s; 
				display: flex; align-items: center; justify-content: center;
			}
			.btn-start { background: var(--success); color: white; }
			.btn-stop { background: var(--danger); color: white; }
			.btn-web { background: #334155; color: white; }
			.btn:disabled { opacity: 0.2; cursor: not-allowed; filter: grayscale(1); }
			
			.console { 
				background: #000; color: #10b981; padding: 12px; border-radius: 8px; 
				font-family: 'Fira Code', 'Consolas', monospace; font-size: 10px; height: 380px; 
				overflow-y: auto; line-height: 1.6; border: 1px solid #1e293b;
			}
			.log-line { margin-bottom: 4px; border-bottom: 1px solid #111; padding-bottom: 2px; }
			.log-time { color: #475569; margin-right: 6px; }
			.log-ERROR { color: #f87171; }
			.log-WARN { color: #fbbf24; }
			
			.card { background: var(--card); border-radius: 10px; padding: 16px; margin-bottom: 16px; border: 1px solid #2d3748; }
			h3 { margin: 0 0 12px 0; font-size: 13px; color: var(--accent); text-transform: uppercase; letter-spacing: 1px; }

			.form-group { margin-bottom: 12px; }
			.form-group label { display: block; margin-bottom: 5px; font-size: 10px; color: var(--text-dim); font-weight: 700; }
			.form-group input { 
				width: 100%; padding: 8px 12px; background: #0f172a; border: 1px solid #334155; 
				border-radius: 6px; color: white; font-size: 12px; outline: none; box-sizing: border-box;
			}
			.form-group input:focus { border-color: var(--accent); }

			.hint { font-size: 10px; color: var(--text-dim); text-align: center; margin-top: 20px; line-height: 1.4; }

			::-webkit-scrollbar { width: 5px; }
			::-webkit-scrollbar-thumb { background: #334155; border-radius: 10px; }
		</style>
	</head>
	<body>
		<div class="header">
			<div>
				<h1>WACAST CORE</h1>
				<div style="font-size: 9px; color: var(--text-dim); margin-top: 2px;">VERSION ` + Version + `</div>
			</div>
			<div class="status-indicator">
				<div id="status-dot" class="dot"></div>
				<span id="status-text">OFFLINE</span>
			</div>
		</div>
		
		<div class="tabs">
			<div id="tab-dash" class="tab active" onclick="showTab('dash')">CONTROL</div>
			<div id="tab-cfg" class="tab" onclick="showTab('cfg')">DATABASE</div>
		</div>
		
		<div id="content-dash" class="content">
			<div class="control-row">
				<button id="btn-start" class="btn btn-start" onclick="startServer()">START</button>
				<button id="btn-stop" class="btn btn-stop" onclick="stopServer()" disabled>STOP</button>
				<button class="btn btn-web" onclick="openDashboard()">DASHBOARD</button>
			</div>
			
			<div style="margin-bottom: 16px;">
				<button class="btn btn-stop" style="width:100%; background: #450a0a; border: 1px solid #991b1b;" onclick="exitApp()">EXIT APPLICATION</button>
			</div>
			
			<div id="update-bar" style="display:none; background: #1e293b; padding: 10px; border-radius: 8px; margin-bottom: 16px; border: 1px solid var(--accent); font-size: 11px;">
				<div style="display:flex; justify-content:space-between; align-items:center;">
					<span id="update-msg">New version available!</span>
					<button class="btn btn-start" style="padding: 4px 10px; font-size: 9px;" id="btn-do-update" onclick="doUpdate()">UPDATE NOW</button>
				</div>
			</div>

			<div class="console" id="console"></div>
			<div style="text-align:center; margin-top:10px">
				<button class="tab" style="padding: 5px 10px; border: 1px solid #334155; border-radius: 4px; font-size: 9px;" onclick="checkUpdate()">Check for Updates</button>
			</div>
			<div class="hint">Window can be closed safely.<br>Reopen from System Tray (near clock).</div>
		</div>
		
		<div id="content-cfg" class="content" style="display:none">
			<div class="card">
				<h3>Connection</h3>
				<div class="form-group"><label>Host</label><input type="text" id="host"></div>
				<div style="display:flex; gap:10px">
					<div class="form-group" style="flex:1"><label>Port</label><input type="number" id="port"></div>
					<div class="form-group" style="flex:2"><label>DB Name</label><input type="text" id="dbname"></div>
				</div>
				<div class="form-group"><label>User</label><input type="text" id="user"></div>
				<div class="form-group"><label>Password</label><input type="password" id="pass"></div>
				<div style="display: flex; gap: 10px; margin-top: 10px;">
					<button class="btn btn-web" style="flex:1" onclick="testDB()">TEST</button>
					<button class="btn btn-start" style="flex:1" onclick="saveDB()">SAVE</button>
				</div>
			</div>
		</div>

		<script>
			let isRunning = false;
			function showTab(id) {
				document.getElementById('content-dash').style.display = id === 'dash' ? 'block' : 'none';
				document.getElementById('content-cfg').style.display = id === 'cfg' ? 'block' : 'none';
				document.getElementById('tab-dash').className = 'tab' + (id === 'dash' ? ' active' : '');
				document.getElementById('tab-cfg').className = 'tab' + (id === 'cfg' ? ' active' : '');
				if(id === 'cfg') loadConfig();
			}
			function log(msg) {
				const con = document.getElementById('console');
				const line = document.createElement('div');
				line.className = 'log-line';
				let cleanMsg = msg;
				if(msg.includes('[ERROR]')) { line.className += ' log-ERROR'; cleanMsg = msg.replace('[ERROR]', ''); }
				if(msg.includes('[WARN]')) { line.className += ' log-WARN'; cleanMsg = msg.replace('[WARN]', ''); }
				line.innerHTML = '<span class="log-time">' + new Date().toLocaleTimeString() + '</span>' + cleanMsg;
				con.appendChild(line);
				con.scrollTop = con.scrollHeight;
				if(con.childNodes.length > 200) con.removeChild(con.firstChild);
			}
			function updateStatus(running) {
				if(isRunning === running) return;
				isRunning = running;
				document.getElementById('status-text').innerText = running ? 'ONLINE' : 'OFFLINE';
				if(running) document.getElementById('status-dot').classList.add('active');
				else document.getElementById('status-dot').classList.remove('active');
				document.getElementById('btn-start').disabled = running;
				document.getElementById('btn-stop').disabled = !running;
			}
			async function startServer() { await fetch('/api/start'); }
			async function stopServer() { await fetch('/api/stop'); }
			async function loadConfig() {
				const c = await fetch('/api/config').then(r => r.json());
				document.getElementById('host').value = c.host; document.getElementById('port').value = c.port;
				document.getElementById('dbname').value = c.dbname; document.getElementById('user').value = c.user;
				document.getElementById('pass').value = c.pass;
			}
			async function saveDB() {
				const body = { host: document.getElementById('host').value, port: parseInt(document.getElementById('port').value), dbname: document.getElementById('dbname').value, user: document.getElementById('user').value, pass: document.getElementById('pass').value };
				const r = await fetch('/api/config', { method: 'POST', body: JSON.stringify(body) }).then(r => r.json());
				alert(r.message);
			}
			async function testDB() {
				const body = { host: document.getElementById('host').value, port: parseInt(document.getElementById('port').value), dbname: document.getElementById('dbname').value, user: document.getElementById('user').value, pass: document.getElementById('pass').value };
				const r = await fetch('/api/test-db', { method: 'POST', body: JSON.stringify(body) }).then(r => r.json());
				alert(r.message);
			}
			function openDashboard() { fetch('/api/open-dashboard'); }
			async function exitApp() { 
				if(confirm('Exit Application?')) {
					await fetch('/api/exit'); 
					window.close();
				}
			}
			
			async function checkUpdate() {
				log('Checking for updates...');
				try {
					const r = await fetch('/api/check-update').then(r => r.json());
					if(r.update_available) {
						document.getElementById('update-bar').style.display = 'block';
						document.getElementById('update-msg').innerText = 'Version ' + r.latest_version + ' is available!';
						log('[WARN] New version found: ' + r.latest_version);
					} else {
						alert('You are on the latest version.');
						log('No updates found.');
					}
				} catch(e) { log('[ERROR] Failed to check for updates'); }
			}

			async function doUpdate() {
				if(!confirm('Application will restart to apply update. Continue?')) return;
				document.getElementById('btn-do-update').disabled = true;
				document.getElementById('btn-do-update').innerText = 'DOWNLOADING...';
				log('Starting update download...');
				try {
					const r = await fetch('/api/do-update').then(r => r.json());
					if(r.success) {
						log('Update applied! Application will restart in 2 seconds...');
					} else {
						alert('Update failed: ' + r.message);
						log('[ERROR] Update failed: ' + r.message);
						document.getElementById('btn-do-update').disabled = false;
						document.getElementById('btn-do-update').innerText = 'UPDATE NOW';
					}
				} catch(e) { log('[ERROR] Update process failed'); }
			}

			setInterval(async () => {
				try { 
					const r = await fetch('/api/status').then(r => r.json()); 
					updateStatus(r.running);
					const l = await fetch('/api/logs').then(r => r.json());
					if(l && l.length > 0) l.forEach(msg => log(msg));
				} catch(e) {}
			}, 800);
		</script>
	</body>
	</html>
	`)
}

func (cp *ControlPanel) handleStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"running": %v}`, cp.app.isRunning)
}

func (cp *ControlPanel) handleLogs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var logs []string
Loop:
	for {
		select {
		case msg := <-utils.LogChannel:
			logs = append(logs, msg)
		default:
			break Loop
		}
	}
	json.NewEncoder(w).Encode(logs)
}

func (cp *ControlPanel) handleStart(w http.ResponseWriter, r *http.Request) {
	go cp.app.Start()
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, `{"success": true}`)
}

func (cp *ControlPanel) handleStop(w http.ResponseWriter, r *http.Request) {
	cp.app.Stop()
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, `{"success": true}`)
}

func (cp *ControlPanel) handleExit(w http.ResponseWriter, r *http.Request) {
	// 1. Force kill any browser window with WACAST title
	if runtime.GOOS == "windows" {
		// Attempt to kill by window title which is very effective for browser apps
		_ = exec.Command("taskkill", "/F", "/FI", "WINDOWTITLE eq WACAST Control Panel*", "/T").Run()
	}

	// 2. Kill the direct GUI process if still alive
	if cp.guiCmd != nil && cp.guiCmd.Process != nil {
		_ = cp.guiCmd.Process.Kill()
	}

	// 3. Stop the core app
	cp.app.Stop()

	if w != nil {
		fmt.Fprint(w, `{"success": true}`)
	}
	time.Sleep(500 * time.Millisecond)
	os.Exit(0)
}


func (cp *ControlPanel) restartApp() {
	// 1. Get the path to the current executable
	self, err := os.Executable()
	if err != nil {
		self = os.Args[0]
	}

	utils.Info("Restarting application...", zap.String("path", self))

	// 2. Start new process using CMD START (Windows specific)
	// This ensures the new process is completely detached from the current shell
	cmd := exec.Command("cmd", "/c", "start", "", self)
	_ = cmd.Run()

	// 3. Exit current process immediately
	cp.app.Stop()
	os.Exit(0)
}

func (cp *ControlPanel) handleConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		cfg, _ := config.LoadConfig()
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"host":"%s", "port":%d, "user":"%s", "pass":"%s", "dbname":"%s"}`,
			cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.DBName)
	} else {
		var req struct {
			Host   string `json:"host"`
			Port   int    `json:"port"`
			User   string `json:"user"`
			Pass   string `json:"pass"`
			DBName string `json:"dbname"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err == nil {
			newCfg := &config.DatabaseConfig{Host: req.Host, Port: req.Port, User: req.User, Password: req.Pass, DBName: req.DBName, SSLMode: "disable"}
			_ = config.SaveDatabaseConfigToEnv(newCfg)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{"success": true, "message": "Config Saved"}`)
		}
	}
}

func (cp *ControlPanel) handleTestDB(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Host   string `json:"host"`
		Port   int    `json:"port"`
		User   string `json:"user"`
		Pass   string `json:"pass"`
		DBName string `json:"dbname"`
	}
	success := false
	msg := "Failed"
	if err := json.NewDecoder(r.Body).Decode(&req); err == nil {
		testCfg := &config.DatabaseConfig{Host: req.Host, Port: req.Port, User: req.User, Password: req.Pass, DBName: req.DBName, SSLMode: "disable"}
		db, err := database.InitDatabase(testCfg)
		if err == nil {
			success = true
			msg = "Connected!"
			db.Close()
		} else {
			msg = fmt.Sprintf("Error: %v", err)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"success": %v, "message": "%s"}`, success, msg)
}

func (cp *ControlPanel) handleOpenDashboard(w http.ResponseWriter, r *http.Request) {
	cp.openDashboard()
}

func (cp *ControlPanel) handleCheckUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// Example: In production, fetch this from your website
	// resp, err := http.Get("https://your-website.com/version.json")
	
	// MOCK: For now, we simulate that version 1.1.0 is available
	currentVersion := Version
	latestVersion := "1.1.0"
	downloadURL := "https://github.com/krismarta18/wa-core-cast/releases/download/v1.1.0/wacast.exe"

	if latestVersion != currentVersion {
		fmt.Fprintf(w, `{"update_available": true, "latest_version": "%s", "url": "%s"}`, latestVersion, downloadURL)
	} else {
		fmt.Fprint(w, `{"update_available": false}`)
	}
}

func (cp *ControlPanel) handleDoUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// 1. Get the download URL (normally from the check-update response)
	downloadURL := "https://github.com/krismarta18/wa-core-cast/releases/download/v1.1.0/wacast.exe"

	// 2. Download the new binary
	resp, err := http.Get(downloadURL)
	if err != nil {
		fmt.Fprintf(w, `{"success": false, "message": "Download failed: %v"}`, err)
		return
	}
	defer resp.Body.Close()

	// 3. Apply the update
	err = update.Apply(resp.Body, update.Options{})
	if err != nil {
		fmt.Fprintf(w, `{"success": false, "message": "Apply failed: %v"}`, err)
		return
	}

	fmt.Fprint(w, `{"success": true}`)
	
	// 4. Trigger restart after a short delay
	go func() {
		time.Sleep(2 * time.Second)
		cp.restartApp()
	}()
}
