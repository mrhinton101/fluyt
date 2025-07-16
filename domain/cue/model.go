package cue

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/mrhinton101/fluyt/internal/app/core/logger"
)

type DeviceSubPaths struct {
	Name    string
	Address string
	Port    string // optional
	Paths   []string
}

type DeviceCapsPaths struct {
	Name    string
	Address string
	Port    string // optional
	Paths   []string
}

type DeviceSubsList struct {
	Devices     []DeviceSubPaths
	dedupTarget map[string]struct{}
}

func NewDeviceSubsList() *DeviceSubsList {
	return &DeviceSubsList{
		Devices:     []DeviceSubPaths{},
		dedupTarget: make(map[string]struct{}),
	}
}

func (d *DeviceSubsList) Add(sub DeviceSubPaths) {
	if _, exists := d.dedupTarget[sub.Address]; exists {
		logger.SLogger(logger.LogEntry{
			Level:     slog.LevelError,
			Err:       errors.New("duplicate target found"),
			Component: "Cue",
			Action:    "add device subscription",
			Msg:       fmt.Sprintf("target %s already exists", sub.Address),
			Target:    sub.Address,
		})
		return
	}
	d.dedupTarget[sub.Address] = struct{}{}
	sub.Paths = dedupList(sub.Paths)
	d.Devices = append(d.Devices, sub)
}

func dedupList(paths []string) []string {
	seen := make(map[string]struct{})
	unique := make([]string, 0, len(paths))
	for _, path := range paths {
		if _, ok := seen[path]; !ok {
			seen[path] = struct{}{}
			unique = append(unique, path)
		}
	}
	return unique
}
