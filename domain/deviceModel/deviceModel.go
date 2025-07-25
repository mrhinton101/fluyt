package devicemodel

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
