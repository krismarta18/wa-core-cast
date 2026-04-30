package utils

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// GetDataDir returns the absolute path to the directory where application data should be stored.
// On Windows, it prefers %APPDATA%\WACAST if installed in Program Files.
// Otherwise, it falls back to the current working directory for portability.
func GetDataDir() string {
	// 1. Check for explicit environment override
	if envDir := os.Getenv("WACAST_DATA_DIR"); envDir != "" {
		_ = os.MkdirAll(envDir, 0755)
		return envDir
	}

	// 2. Determine base directory
	exePath, err := os.Executable()
	if err != nil {
		return "."
	}
	exeDir := filepath.Dir(exePath)

	// 3. If running on Windows and inside Program Files, use AppData
	if runtime.GOOS == "windows" {
		programFiles := os.Getenv("ProgramFiles")
		programFilesX86 := os.Getenv("ProgramFiles(x86)")

		isSystemDir := false
		if programFiles != "" && strings.HasPrefix(strings.ToLower(exeDir), strings.ToLower(programFiles)) {
			isSystemDir = true
		} else if programFilesX86 != "" && strings.HasPrefix(strings.ToLower(exeDir), strings.ToLower(programFilesX86)) {
			isSystemDir = true
		}

		if isSystemDir {
			appData := os.Getenv("APPDATA")
			if appData != "" {
				dataDir := filepath.Join(appData, "WACAST")
				_ = os.MkdirAll(dataDir, 0755)
				return dataDir
			}
		}
	}

	// Default fallback to exe directory (portable mode)
	return exeDir
}

// GetDataPath joins the data directory with the provided elements
func GetDataPath(elem ...string) string {
	base := GetDataDir()
	return filepath.Join(append([]string{base}, elem...)...)
}
