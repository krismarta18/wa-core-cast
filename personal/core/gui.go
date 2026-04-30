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
)

//go:embed icon/favicon.png
var iconFS embed.FS

// Minimalist Blue "W" Icon (32x32 PNG Base64) - Kept for reference or future use
const trayIconBase64 = "iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAYAAABzenr0AAAACXBIWXMAAAsTAAALEwEAmpwYAAADG0lEQVR4nO2Xv2sbVxTHP+fdnS0pUp2Sxk0pZAsGAt0MtYfSh9AsHTrX/oB0y9Chf0C6du9m6NChY+0f0K106BToYujWEEpD6ZCQptYf9Z0OnS2rS9I797pD8pAs2ZZkJ8W29AsH7vE+730+vO97T8QYw7O2uW7yYhV4BvA68AzgDeAx4FfgY+CDZdnffz98HkS6i6vAc8BrwAunYv9vWfYPvwvgecDzLh7A8yB46HkXz+vAsy6e14HnXbz/9p548TzYxT908X/86XF7X9z03i6ee66Lp6fN7+X7/6UvXj676+LpcXtf3Mvd630ZunhuuXh6XU9f97R5en947rmf/2Nf/FfX9fS06Uv7+b467nU9u66np82ve9o8fXC0+PzX88S79oDAtp19X8f/9X48p81/O0+8v7S3m7tX9fS0+T0v7ul/7In70p8en7S9p80XezWvO6+nv9Xz9X7N6+nPj9ve0+ax87qXz964Wz1tXp94Xf16+rvz7InXndfP98S9Wp/m9R/O09e99t78p67r6elr85fHbe9p85dzf9uA6vXNf6vn9fS0ef6UePe996+nr+vpa/OXp70T97QeHLe9p81fzn/5370A+H/8E/A/A9zH8GvA4/Fv5tL93x8T8H7AnwB+Bn6L8a/AvxHjZ8A/Y/wD8PcY/wT89Y8A/hHjn4B/ZPrX8v5/f0zA7zH+Bfgn4K+Zfne57v8B/9v0p8t1/w/4r6Y/Xa77H8Z/Nf3p5br/Dfy76U+X6/4Z+P9m+tPlev8Y/x/0+vSPl/N9v6u8R/6mXp/+fXW5679WfX+5679m9f3V6T8u5/sPqvfI39XfV6f/Wp098nfV95e7/tuq769O/305339f9R75p67vL3f991XfX53+93K+/1H1HvnXre8vd/3PtT797+V8/1PVe+TfVfX95a7/ver7q9P/Xc73P68vPZ7m5Xp62vzf6vXpX1TviXvS+m6Z9/X0tfm36vXpX1fviXva+v5y3f9Svd+X6unpe9r6vnLdf1f19O/r6evr6+vpe/p6evr6enr6evr6evr6evr6evr6+voX/S80I7N7VqO21AAAAABJRU5ErkJggg=="

type ControlPanel struct {
	app          *App
	launcherPort int
}

func NewControlPanel() *ControlPanel {
	return &ControlPanel{
		app:          NewApp(),
		launcherPort: 19991,
	}
}

func (cp *ControlPanel) Run() {
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
	time.Sleep(1 * time.Second)

	utils.Info("WACAST Control Panel is active.")
	utils.Info("Closing this window will stop the application.")

	// Launch UI
	go cp.OpenUI()
	
	// Block here so the app stays alive until /api/exit or Ctrl+C
	utils.Info("Control Panel is running. Use the 'Exit' button in the UI or Ctrl+C to stop.")
	
	// Keep the main thread alive
	select {}
}

func (cp *ControlPanel) OpenUI() {
	url := fmt.Sprintf("http://127.0.0.1:%d", cp.launcherPort)
	var cmd *exec.Cmd
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
			cmd = exec.Command("cmd", "/c", "start", "", edgePath, "--app="+url, "--window-size=380,640")
			_ = cmd.Run()
			return
		}
	} else {
		cmd = exec.Command("google-chrome", "--app="+url, "--window-size=380,640")
		_ = cmd.Run()
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
				<button class="btn btn-web" onclick="openDashboard()">WEB UI</button>
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
			<button class="btn btn-stop" style="width:100%" onclick="exitApp()">EXIT APPLICATION</button>
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
			function exitApp() { if(confirm('Exit Application?')) fetch('/api/exit'); }
			
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
