package model

import "github.com/mrhinton101/fluyt/domain/gnmi"

type Device struct {
	Name  string
	IP    string
	Make  string
	Model string
	// Switching
	Routing Routing
	// Interfaces map[string]Interface

}

type Routing struct {
	Protocols SupportedProtocols
}

type SupportedProtocols struct {
	BGP BgpProcess
	// OSPF OSPF
	// EIGRP EIGRP
	// Static Static
}

type BgpProcess struct {
	Runtime BgpRuntime
}

type BgpRuntime struct {
	Rib gnmi.BgpRibs
}

// type Interface struct {
// 	InterfaceName string
// 	IPAddresses []string
// 	MacAddress string
// 	Speed int
// 	AdminStatus string
// 	OperStatus string
// }

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
