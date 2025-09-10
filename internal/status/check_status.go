package status

import (
	"MiniProject/internal/types"
	"log/slog"
	"time"

	"github.com/prometheus-community/pro-bing"
)

func CheckStatus(device *types.Device) {

	pinger, err := probing.NewPinger(device.IpAddress)

	if err != nil {

		slog.Error("Error creating pinger: ", err)

		device.Status = "InActive"
	}

	pinger.SetPrivileged(true)

	pinger.Count = 5

	pinger.Timeout = time.Second * 5

	slog.Info("Pinging ", pinger.Addr())

	err = pinger.Run()

	if err != nil {
		slog.Error("Error pinging ", pinger.Addr(), ": ", err)
		device.Status = "InActive"
	}

	stats := pinger.Statistics()

	slog.Info("Statistics: ", stats)

	successRate := float64(stats.PacketsRecv) / float64(stats.PacketsSent)

	slog.Info("Success rate: ", successRate)

	if successRate > 0.33 {
		device.Status = "Active"
	} else {
		device.Status = "Inactive"
	}
}
