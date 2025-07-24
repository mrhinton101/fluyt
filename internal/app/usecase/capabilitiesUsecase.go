package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/mrhinton101/fluyt/domain/cue"
	"github.com/mrhinton101/fluyt/domain/gnmi"
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
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			initErr := client.Init(ctx)
			if initErr != nil {
				fmt.Println(initErr)
			}
			caps, err := client.Capabilities(ctx)
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
			client.Close()

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

func CollectBgpRib(devices *cue.DeviceList, clientFactory func(cue.Device) ports.GNMIClient) gnmi.BgpRibs {
	results := gnmi.BgpRibs{}
	resultChan := make(chan gnmi.BgpRibs)
	var wg sync.WaitGroup
	for _, device := range devices.Devices {
		wg.Add(1)
		client := clientFactory(device)
		go func() {
			defer wg.Done()
			fmt.Printf("starting bgp goroutine: %v", client.GetAddress())
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			initErr := client.Init(ctx)
			if initErr != nil {
				fmt.Println(initErr)
			}
			bgpRib, err := client.GetBgpRibs(ctx)
			if err != nil {
				logger.SLogger(logger.LogEntry{
					Level:     slog.LevelError,
					Err:       err,
					Component: "gNMI",
					Action:    "GetBgpRib",
					Msg:       "failed to fetch BGP RIB",
					Target:    client.GetAddress(),
				})
			}
			client.Close()

			resultChan <- bgpRib
			fmt.Println("goroutine written to resultchan")
		}()
	}
	go func() {
		wg.Wait()
		close(resultChan)
	}()
	for result := range resultChan {
		for vrfName, rib := range result {
			if _, exists := results[vrfName]; !exists {
				results[vrfName] = make(map[gnmi.BgpRibKey]gnmi.BgpRibRoute)
			}
			for key, route := range rib {
				results[vrfName][key] = route
			}
		}
	}
	return results

}
