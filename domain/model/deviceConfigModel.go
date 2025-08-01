package model

type DeviceConfigRouting struct {
	Bgp BgpProc
	// OSPF DeviceConfigOSPF
}

type DeviceConfig struct {
	// Physical
	// Switching
	Routing DeviceConfigRouting
	// Overlay
	// Platform
}

type Vlans struct {
}
type BgpNeighbor struct {
	PeerGroup   string `yaml:"peer_group"`
	Description string `yaml:"description"`
	RemoteAs    uint32 `yaml:"remote_as"`
}

type BgpProc struct {
	AS        uint32                 `yaml:"as"`
	RouterID  string                 `yaml:"router_id"`
	Neighbors map[string]BgpNeighbor `yaml:"neighbors"`
}
