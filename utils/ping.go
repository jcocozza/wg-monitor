package utils

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/go-ping/ping"
)

// Check if a given address in up
func IsReachable(addr string) bool {
	slog.Debug(fmt.Sprintf("pinging %s", addr))

	pinger, err := ping.NewPinger(addr)
	if err != nil {
		slog.Error("Error creating pinger:", err)
		return false
	}

	pinger.Count = 1 // Number of ICMP packets to send
	pinger.Timeout = 1 * time.Second // Timeout for each packet

	err = pinger.Run()
    if err != nil {
        slog.Error("Failed to ping:", err)
        return false
    }
	stats := pinger.Statistics()

	if stats.PacketsSent == stats.PacketsRecv {
		slog.Debug(fmt.Sprintf("%s is up...", addr))
	} else {
		slog.Debug(fmt.Sprintf("%s is down...", addr))
	}
	return stats.PacketsSent == stats.PacketsRecv
}