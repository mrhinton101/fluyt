package cue

import (
	"errors"
	"fmt"
	"log/slog"

	"cuelang.org/go/cue"

	"github.com/mrhinton101/fluyt/internal/app/core/logger"
)

type ConcreteInv struct {
	Inventory cue.Value
}

func (inv *ConcreteInv) LoadDevices() (*DeviceList, error) {
	list := NewDeviceList()

	iter, err := inv.Inventory.Fields()
	if err != nil {
		return nil, fmt.Errorf("cannot iterate inventory: %w", err)
	}

	for iter.Next() {
		deviceName := iter.Selector().String()
		deviceVal := iter.Value()

		ipVal := deviceVal.LookupPath(cue.ParsePath("ip"))
		ipStr, err := ipVal.String()
		if err != nil {
			logError("ip", deviceName, err)
			continue
		}
		portVal := deviceVal.LookupPath(cue.ParsePath("port"))
		portStr, _ := portVal.String() // optional

		list.Add(Device{
			Name:    deviceName,
			Address: ipStr,
			Port:    portStr,
		})
	}

	return list, nil
}

func (inv *ConcreteInv) LoadSubs() (*DeviceSubsList, error) {
	list := NewDeviceSubsList()

	iter, err := inv.Inventory.Fields()
	if err != nil {
		return nil, fmt.Errorf("cannot iterate inventory: %w", err)
	}

	for iter.Next() {
		deviceName := iter.Selector().String()
		deviceVal := iter.Value()

		ipVal := deviceVal.LookupPath(cue.ParsePath("ip"))
		ipStr, err := ipVal.String()
		if err != nil {
			logError("ip", deviceName, err)
			continue
		}
		// Create device entry
		device := Device{
			Name:    deviceName,
			Address: ipStr,
		}

		telPaths := deviceVal.LookupPath(cue.ParsePath("tel_paths"))
		if !telPaths.Exists() {
			logError("tel_paths", deviceName, errors.New("missing tel_paths"))
			continue
		}

		paths := extractTelemetryPaths(telPaths, deviceName)
		list.Add(DeviceSubPaths{
			Device: device,
			Paths:  paths,
		})
	}

	return list, nil
}

func extractTelemetryPaths(telPaths cue.Value, device string) []string {
	var result []string
	iter, _ := telPaths.List()

	for iter.Next() {
		entry := iter.Value()
		subIter, err := entry.Fields()
		if err != nil {
			logError("fields", device, err)
			continue
		}

		for subIter.Next() {
			valStr, err := subIter.Value().String()
			if err == nil {
				result = append(result, valStr)
			}
		}
	}
	return result
}

func logError(field, device string, err error) {
	logger.SLogger(logger.LogEntry{
		Level:     slog.LevelError,
		Err:       err,
		Component: "Cue",
		Action:    "extract inventory",
		Msg:       fmt.Sprintf("Device %s - %s error: %v", device, field, err),
		Target:    device,
	})
}
