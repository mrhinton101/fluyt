package cueHandler

import (
	"errors"
	"fmt"
	"log/slog"

	"cuelang.org/go/cue"
	logger "github.com/mrhinton101/fluyt/internal/adapter/logger"
)

func (inv *concreteInv) loadSubs() (DeviceSubsList *DeviceSubsList) {
	inventory := inv.inventory
	iter_inventory, err := inventory.Fields()
	if err != nil {
		logger.SLogger(logger.LogEntry{
			Level:     slog.LevelError,
			Err:       err,
			Component: "Cue",
			Action:    "load devices fields",
			Msg:       "failed to iterate over devices fields. are they defined?",
			Target:    "localhost",
		})
	}
	DeviceSubsList = initDeviceSubsList()
	for iter_inventory.Next() {
		deviceName := iter_inventory.Selector()
		deviceVal := iter_inventory.Value()

		ipVal := deviceVal.LookupPath(cue.ParsePath("ip"))
		deviceNameStr := deviceName.String()
		ipStr, err := ipVal.String()
		telemetryPaths := deviceVal.LookupPath(cue.ParsePath("tel_paths"))
		if !telemetryPaths.Exists() {
			logger.SLogger(logger.LogEntry{
				Level:     slog.LevelError,
				Err:       errors.New("no telemetry paths defined"),
				Component: "Cue",
				Action:    "get telemetry paths",
				Msg:       fmt.Sprintf("Could not get telemetry paths for device %s: %v", deviceName, err),
				Target:    "localhost",
			})
		}
		iter_telemetry, err := telemetryPaths.List()
		telemPathList := []string{}
		for iter_telemetry.Next() {

			telemetryVal := iter_telemetry.Value()

			iter_telemetry_paths, err := telemetryVal.Fields()
			if err != nil {
				logger.SLogger(logger.LogEntry{
					Level:     slog.LevelError,
					Err:       err,
					Component: "Cue",
					Action:    "get telemetry paths",
					Msg:       fmt.Sprintf("Could not get telemetry paths for device %s: %v", deviceName, err),
					Target:    "localhost",
				})
			}
			for iter_telemetry_paths.Next() {
				telemetryPathVal := iter_telemetry_paths.Value()
				telemetryPathStr, err := telemetryPathVal.String()
				if err != nil {
					logger.SLogger(logger.LogEntry{
						Level:     slog.LevelError,
						Err:       err,
						Component: "Cue",
						Action:    "convert telemetry value to string",
						Msg:       fmt.Sprintf("Could not convert telemetry value to string for device %s: %v", deviceName, err),
						Target:    "localhost",
					})

				}
				telemPathList = append(telemPathList, telemetryPathStr)
				if err != nil {
					logger.SLogger(logger.LogEntry{
						Level:     slog.LevelError,
						Err:       err,
						Component: "Cue",
						Action:    "get telemetry value",
						Msg:       fmt.Sprintf("Could not get telemetry value for device %s: %v", deviceName, err),
						Target:    "localhost",
					})
				}
			}

		}

		if err != nil {
			logger.SLogger(logger.LogEntry{
				Level:     slog.LevelError,
				Err:       err,
				Component: "Cue",
				Action:    "find ip",
				Msg:       fmt.Sprintf("Could not get IP for device %s: %v", deviceName, err),
				Target:    "localhost",
			})
		}
		device := DeviceSubPaths{
			DeviceInfo: DeviceInfo{
				Name:    deviceNameStr,
				Address: ipStr},
			Paths: telemPathList,
		}
		DeviceSubsList.Add(device)

	}
	return DeviceSubsList
}
func (inv *concreteInv) loadCaps() (DeviceList *DeviceList) {
	inventory := inv.inventory
	iter_inventory, err := inventory.Fields()
	if err != nil {
		logger.SLogger(logger.LogEntry{
			Level:     slog.LevelError,
			Err:       err,
			Component: "Cue",
			Action:    "load devices fields",
			Msg:       "failed to iterate over devices fields. are they defined?",
			Target:    "localhost",
		})
	}
	for iter_inventory.Next() {
		deviceName := iter_inventory.Selector()
		deviceVal := iter_inventory.Value()

		ipVal := deviceVal.LookupPath(cue.ParsePath("ip"))
		deviceNameStr := deviceName.String()
		ipStr, err := ipVal.String()
		telemetryPaths := deviceVal.LookupPath(cue.ParsePath("tel_paths"))
		if !telemetryPaths.Exists() {
			logger.SLogger(logger.LogEntry{
				Level:     slog.LevelError,
				Err:       errors.New("no telemetry paths defined"),
				Component: "Cue",
				Action:    "get telemetry paths",
				Msg:       fmt.Sprintf("Could not get telemetry paths for device %s: %v", deviceName, err),
				Target:    "localhost",
			})
		}

		if err != nil {
			logger.SLogger(logger.LogEntry{
				Level:     slog.LevelError,
				Err:       err,
				Component: "Cue",
				Action:    "find ip",
				Msg:       fmt.Sprintf("Could not get IP for device %s: %v", deviceName, err),
				Target:    "localhost",
			})
		}
		device := DeviceInfo{
			Name:    deviceNameStr,
			Address: ipStr}

		DeviceList.Add(device)

	}
	return DeviceList
}

func LoadDeviceSubsList(schemaDir, invFile string) (DeviceSubsList *DeviceSubsList) {
	ctx, schemaVals := loadSchemaDir(schemaDir)
	concreteInv := loadInventory(ctx, schemaVals, invFile)
	DeviceSubsList = concreteInv.loadSubs()
	return DeviceSubsList
}
