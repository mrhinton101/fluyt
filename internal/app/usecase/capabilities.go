package usecase

import (
	"log/slog"

	"github.com/mrhinton101/fluyt/domain/cue"
	"github.com/mrhinton101/fluyt/internal/app/core/logger"
	"github.com/mrhinton101/fluyt/internal/app/ports"
)

type CapabilitiesResult struct {
	Target    string
	Encodings []string
	Models    []string
	Versions  string
}

func CollectCapabilities(devices *cue.DeviceSubsList, clientFactory func(cue.DeviceSubPaths) ports.GNMIClient) []CapabilitiesResult {
	results := []CapabilitiesResult{}

	for _, device := range devices.Devices {
		client := clientFactory(device)

		caps, err := client.Capabilities()
		if err != nil {
			logger.SLogger(logger.LogEntry{
				Level:     slog.LevelError,
				Err:       err,
				Component: "gNMI",
				Action:    "capabilities",
				Msg:       "failed to fetch capabilities",
				Target:    device.Address,
			})
			continue
		}

		results = append(results, CapabilitiesResult{
			Target:    device.Address,
			Encodings: caps["encodings"].([]string),
			Models:    caps["models"].([]string),
			Versions:  caps["versions"].(string),
		})
	}

	return results
}
