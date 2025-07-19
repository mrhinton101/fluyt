package usecase

import (
	"fmt"
	"log/slog"
	"sync"

	"github.com/mrhinton101/fluyt/domain/cue"
	"github.com/mrhinton101/fluyt/internal/app/core/logger"
	"github.com/mrhinton101/fluyt/internal/app/ports"
)

type CapabilitiesResultList struct {
	Results []CapabilitiesResult
}

type CapabilitiesResult struct {
	Target    string
	Encodings []string
	Models    []string
	Versions  string
}

func CollectCapabilities(devices *cue.DeviceList, clientFactory func(cue.Device) ports.GNMIClient) []CapabilitiesResult {
	results := []CapabilitiesResult{}
	resultChan := make(chan CapabilitiesResult)
	var wg sync.WaitGroup
	for _, device := range devices.Devices {
		wg.Add(1)
		client := clientFactory(device)
		// result := CapabilitiesResult{}
		go func() {
			defer wg.Done()
			fmt.Printf("starting goroutine: %v", client.GetAddress())
			caps, err := client.Capabilities()
			if err != nil {
				logger.SLogger(logger.LogEntry{
					Level:     slog.LevelError,
					Err:       err,
					Component: "gNMI",
					Action:    "capabilities",
					Msg:       "failed to fetch capabilities",
					Target:    client.GetAddress(),
				})
			}

			result := CapabilitiesResult{
				Target:    client.GetAddress(),
				Encodings: caps["encodings"].([]string),
				Models:    caps["models"].([]string),
				Versions:  caps["versions"].(string),
			}
			resultChan <- result
			fmt.Println("goroutine written to resultchan")
		}()
	}
	go func() {
		wg.Wait()
		close(resultChan)
	}()
	for result := range resultChan {
		results = append(results, result)
	}

	return results
}
