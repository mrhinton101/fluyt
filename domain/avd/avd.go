package avd

import "github.com/mrhinton101/fluyt/domain/model"

type AvdDeviceConfig struct {
	// Physical
	// Switching
	Bgp model.BgpProc `yaml:"router_bgp"`
	// Vlans map[string]model.Vlans
	// Overlay
	// Platform
}


func 