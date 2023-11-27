package commands

import (
	"log/slog"
	"os/exec"
)

// run wg-quick up <confPath>
func WgQuickUp(confPath string) {
	cmd := exec.Command("wg-quick", "up", confPath)

	_, err := cmd.CombinedOutput()
	if err != nil {
		slog.Debug("Failed to run wg-quick up "+ confPath)
	}
}

// run wg-quick down <confPath>
func WgQuickDown(confPath string) {
	cmd := exec.Command("wg-quick", "down", confPath)

	_, err := cmd.CombinedOutput()
	if err != nil {
		slog.Debug("Failed to run wg-quick down "+ confPath)
	}
}