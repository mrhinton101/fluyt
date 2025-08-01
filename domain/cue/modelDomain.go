package cue

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/mrhinton101/fluyt/domain/model"
	"github.com/mrhinton101/fluyt/internal/app/core/logger"
)

type DeviceSubPaths struct {
	Device
	Paths []string
}
type Device struct {
	Name    string
	Address string
	Port    string // optional
	Config  model.DeviceConfig
}

type DeviceList struct {
	Devices []Device
}

func (d DeviceList) GetByName(name string) (*DeviceList, bool) {
	for _, device := range d.Devices {
		if device.Name == name {
			return &DeviceList{Devices: []Device{device}}, true
		}
	}
	return nil, false
}

// func (d DeviceList) GetByAddress(address string) *Device {
// 	for _, device := range d.Devices {
// 		if device.Address == address {
// 			return &device
// 		}
// 	}
// 	return nil
// }

func (d DeviceList) All() *DeviceList {
	return &DeviceList{Devices: d.Devices}
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

func NewDeviceList() *DeviceList {
	return &DeviceList{
		Devices: []Device{},
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

func (d *DeviceList) Add(sub Device) {
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
