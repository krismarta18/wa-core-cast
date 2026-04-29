# Automation Script for Wacast Build & Deploy
# Usage: ./build.ps1

$ErrorActionPreference = "Stop"

Write-Host "`n[INFO] Starting Automation Build & Deploy Process...`n" -ForegroundColor Green

# 1. Build Dashboard
Write-Host "[STEP 1] Building Dashboard (Next.js)..." -ForegroundColor Cyan
Push-Location dashboard
npm run build
Pop-Location

# 2. Update Core AppServer Static Files
Write-Host "[STEP 2] Syncing Dashboard to Core AppServer..." -ForegroundColor Cyan
$sourceDir = "dashboard/dashboard_out"
$destDir = "core/appserver/dashboard_out"

if (Test-Path $destDir) {
    Write-Host "   - Cleaning old dashboard files..." -ForegroundColor Gray
    Remove-Item -Path $destDir -Recurse -Force
}
# Move-Item might fail if destination exists, but we just removed it.
# However, Move-Item dashboard/dashboard_out to core/appserver/ will create core/appserver/dashboard_out
Move-Item -Path $sourceDir -Destination "core/appserver/"

# 3. Build Go Application
Write-Host "[STEP 3] Building Go Core (main.exe)..." -ForegroundColor Cyan
Push-Location core
go build -o main.exe .
Pop-Location

# 4. Prepare Production Folder
Write-Host "[STEP 4] Cleaning Production Folder (preserving licenses)..." -ForegroundColor Cyan
$prodDir = "Production"

if (!(Test-Path $prodDir)) {
    New-Item -ItemType Directory -Path $prodDir
} else {
    # Get all items in Production
    $items = Get-ChildItem -Path $prodDir
    foreach ($item in $items) {
        # Preserve files with .lic extension or containing 'license' in name
        if ($item.Extension -ne ".lic" -and $item.Name -notlike "*license*") {
            try {
                Remove-Item -Path $item.FullName -Recurse -Force -ErrorAction SilentlyContinue
            } catch {
                Write-Host "   - Could not remove $($item.Name) (might be in use)" -ForegroundColor Red
            }
        } else {
            Write-Host "   - Preserving: $($item.Name)" -ForegroundColor Yellow
        }
    }
}

# 5. Move Binary to Production
Write-Host "[STEP 5] Deploying main.exe to Production..." -ForegroundColor Cyan
if (Test-Path "Production/main.exe") {
    # Try to kill main.exe if it's running
    Write-Host "   - Stopping running main.exe..." -ForegroundColor Gray
    taskkill /F /IM main.exe /T 2>$null
    Start-Sleep -Seconds 2
}
Move-Item -Path "core/main.exe" -Destination "$prodDir/main.exe" -Force

Write-Host "`n[DONE] Build and Deploy Successful." -ForegroundColor Green
Write-Host "Binary is ready in: Production\main.exe`n"
