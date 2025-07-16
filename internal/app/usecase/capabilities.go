package usecase

import (
	"fmt"
	"log/slog"

	"github.com/mrhinton101/fluyt/domain/cue"
	"github.com/mrhinton101/fluyt/internal/app/core/logger"
	"github.com/mrhinton101/fluyt/internal/app/ports"
)

type CapabilitiesResult struct {
	Target    string
	Encodings []string
	Models    []string
	Versions  []string
}

func CollectCapabilities(devices *cue.DeviceSubsList, clientFactory func(string) ports.GNMIClient) []CapabilitiesResult {
	results := []CapabilitiesResult{}

	for _, dev := range devices.Devices {
		client := clientFactory(dev.Address)

		if err := client.Connect(); err != nil {
			logger.SLogger(logger.LogEntry{
				Level:     slog.LevelError,
				Err:       err,
				Component: "gNMI",
				Action:    "connect",
				Msg:       "failed to connect",
				Target:    dev.Address,
			})
			continue
		}
		defer client.Close()

		caps, err := client.Capabilities()
		if err != nil {
			logger.SLogger(logger.LogEntry{
				Level:     slog.LevelError,
				Err:       err,
				Component: "gNMI",
				Action:    "capabilities",
				Msg:       "failed to fetch capabilities",
				Target:    dev.Address,
			})
			continue
		}
		fmt.Println(caps)

		results = append(results, CapabilitiesResult{
			Target:    dev.Address,
			Encodings: caps["encodings"],
			Models:    caps["models"],
			Versions:  caps["versions"],
		})
	}

	return results
}
