package cue

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"cuelang.org/go/cue"
	"github.com/mrhinton101/fluyt/domain/avd"
	"github.com/mrhinton101/fluyt/domain/model"
	"github.com/mrhinton101/fluyt/internal/app/core/logger"
)

type ConcreteInv struct {
	Inventory cue.Value
}

func (inv *ConcreteInv) LoadDevices() (*model.DeviceList, error) {
	list := model.NewDeviceList()
	var deviceConfig model.DeviceConfig

	iter, err := inv.Inventory.Fields()
	if err != nil {
		logger.SLogger(logger.LogEntry{
			Level:     slog.LevelError,
			Err:       err,
			Component: "Cue",
			Action:    "load inventory",
			Msg:       fmt.Sprintln("failed to iterate over inventory"),
			Target:    "localhost",
		})
		return nil, err
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

		configPathVal := deviceVal.LookupPath(cue.ParsePath("config_file"))
		configPathStr, _ := configPathVal.String() // optional

		configFmtVal := deviceVal.LookupPath(cue.ParsePath("config_format"))
		configFmtStr, _ := configFmtVal.String() // optional
		// fmt.Println(configFmtVal.String())

		lowerCfgFmt := strings.ToLower(configFmtStr)

		switch lowerCfgFmt {
		case "avd":
			deviceConfig, err = avd.AvdToModel(configPathStr)
			if err != nil {
				logger.SLogger(logger.LogEntry{
					Level:     slog.LevelError,
					Err:       err,
					Component: "Cue",
					Action:    "Read local File",
					Msg:       fmt.Sprintf("error reading from path %s", configPathStr),
					Target:    configPathStr,
				})
				return nil, err
			}
		default:
			err := fmt.Errorf("unsupported config file format: %s", configFmtStr)
			logger.SLogger(logger.LogEntry{
				Level:     slog.LevelError,
				Err:       err,
				Component: "Cue",
				Action:    "determine config file format",
				Msg:       fmt.Sprintf("currently only 'AVD' is supported config format. You provided  %s", configFmtStr),
				Target:    configPathStr,
			})
			return nil, err
		}

		// fmt.Println(deviceConfig)

		list.Add(model.Device{
			Name:    deviceName,
			Address: ipStr,
			Port:    portStr,
			Config:  deviceConfig,
		})
	}

	return list, nil
}

func (inv *ConcreteInv) LoadSubs() (*model.DeviceSubsList, error) {
	list := model.NewDeviceSubsList()

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
		device := model.Device{
			Name:    deviceName,
			Address: ipStr,
		}

		telPaths := deviceVal.LookupPath(cue.ParsePath("tel_paths"))
		if !telPaths.Exists() {
			logError("tel_paths", deviceName, errors.New("missing tel_paths"))
			continue
		}

		paths := extractTelemetryPaths(telPaths, deviceName)
		list.Add(model.DeviceSubPaths{
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
