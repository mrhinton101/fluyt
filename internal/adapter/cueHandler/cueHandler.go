package cueHandler

import (
	"errors"
	"fmt"
	"log/slog"

	"cuelang.org/go/cue"
	"cuelang.org/go/encoding/yaml"
	logger "github.com/mrhinton101/fluyt/internal/adapter/logger"
)

type DeviceSubsList struct {
	Devices     []DeviceSubPaths
	dedupTarget map[string]struct{}
}

type DeviceList struct {
	Devices []DeviceInfo
}

type DeviceInfo struct {
	Name    string
	Address string
	Port    string
}

type DeviceSubPaths struct {
	DeviceInfo
	Paths []string
}

type concreteInv struct {
	inventory cue.Value
}

func initDeviceSubsList() *DeviceSubsList {
	return &DeviceSubsList{
		Devices:     []DeviceSubPaths{},
		dedupTarget: make(map[string]struct{}),
	}
}

func dedupList(inputSlice []string) []string {
	seen := make(map[string]struct{})
	outputSlice := make([]string, 0, len(inputSlice))
	for _, item := range inputSlice {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			outputSlice = append(outputSlice, item)
		} else {
			logger.SLogger(logger.LogEntry{
				Level:     slog.LevelError,
				Err:       errors.New("duplicate entry found"),
				Component: "Cue",
				Action:    "dedup list",
				Msg:       fmt.Sprintf("duplicate entry %s found in list", item),
				Target:    item,
			})
		}
	}
	return outputSlice
}

func (d *DeviceList) Add(device DeviceInfo) {
	d.Devices = append(d.Devices, device)
}

func (d *DeviceSubsList) Add(sub DeviceSubPaths) {

	if _, dupe := d.dedupTarget[sub.Address]; dupe {
		logger.SLogger(logger.LogEntry{
			Level:     slog.LevelError,
			Err:       errors.New("duplicate target found"),
			Component: "Cue",
			Action:    "add device subscription",
			Msg:       fmt.Sprintf("target %s already exists in the list", sub.Address),
			Target:    sub.Address,
		})
		return
	}
	d.dedupTarget[sub.Address] = struct{}{}
	cleanPathList := dedupList(sub.Paths)
	sub.Paths = cleanPathList
	d.Devices = append(d.Devices, sub)
}

func pathLookup(path string, schema cue.Value) (pathResults cue.Value) {
	// Extract the `#Inventory` definition from schema
	pathResults = schema.LookupPath(cue.ParsePath(path))
	if !pathResults.Exists() {
		logger.SLogger(logger.LogEntry{
			Level:     slog.LevelError,
			Err:       errors.New("failed to load top level schema key"),
			Component: "Cue",
			Action:    "load Cue top level key",
			Msg:       fmt.Sprintf("schema missing %s definition", path),
			Target:    "localhost",
		})
	}
	return pathResults

}

func loadYaml(ctx *cue.Context, yamlfile string) (yamlVal cue.Value) {
	// Load YAML inventory file
	yamlFile, err := yaml.Extract(yamlfile, nil)
	if err != nil {
		logger.SLogger(logger.LogEntry{
			Level:     slog.LevelError,
			Err:       err,
			Component: "Cue",
			Action:    "load yaml file",
			Msg:       fmt.Sprintf("failed to load %v. Is the location correct", yamlFile),
			Target:    "localhost",
		})
	}
	yamlVal = ctx.BuildFile(yamlFile)
	return yamlVal

}
