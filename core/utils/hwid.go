package utils

import (
	"crypto/sha256"
	"fmt"
	"os/exec"
	"strings"
)

// GetHWID returns a unique hardware fingerprint for the current Windows machine.
// It combines the Motherboard Serial Number and the CPU Processor ID.
func GetHWID() (string, error) {
	// 1. Get Motherboard Serial Number
	mbSerial, err := getWmicInfo("baseboard", "get", "serialnumber")
	if err != nil {
		return "", fmt.Errorf("failed to get motherboard serial: %w", err)
	}

	// 2. Get CPU Processor ID
	cpuID, err := getWmicInfo("cpu", "get", "processorid")
	if err != nil {
		return "", fmt.Errorf("failed to get cpu id: %w", err)
	}

	// Combine and hash for a clean fingerprint
	raw := fmt.Sprintf("WACAST-HWID-%s-%s", mbSerial, cpuID)
	hash := sha256.Sum256([]byte(raw))
	
	// Convert to hex and take first 16 characters for a manageable ID
	fullHex := fmt.Sprintf("%x", hash)
	displayID := strings.ToUpper(fullHex[:16])

	return displayID, nil
}

func getWmicInfo(category, action, property string) (string, error) {
	cmd := exec.Command("wmic", category, action, property)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(out), "\n")
	if len(lines) < 2 {
		return "", fmt.Errorf("unexpected wmic output format")
	}

	// The second line usually contains the actual value
	result := strings.TrimSpace(lines[1])
	if result == "" || strings.EqualFold(result, "To be filled by O.E.M.") {
		// Fallback for some machines where serial is not populated
		return "GENERIC-HARDWARE-ID", nil
	}

	return result, nil
}
