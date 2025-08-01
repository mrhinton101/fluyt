package gnmi

type CleanCapabilityResponse struct {
	Target    string
	Encodings []string
	Models    []string
	Versions  string
}

type BgpRibRoute struct {
	Origin string `json:"origin"`
	PathID int    `json:"path-id"`
	Prefix string `json:"prefix"`
	State  struct {
		AttrIndex    string `json:"attr-index"`
		LastModified string `json:"last-modified"`
		Origin       string `json:"origin"`
		PathID       int    `json:"path-id"`
		Prefix       string `json:"prefix"`
		ValidRoute   bool   `json:"valid-route"`
	} `json:"state"`
}
type BgpRibKey struct {
	Prefix string
}

type BgpVrfName struct {
	Name string
}

type Device string

type BgpRib map[BgpVrfName]map[BgpRibKey]BgpRibRoute

type BgpRibs map[Device]BgpRib

type GnmiBgpRibRoutes struct {
	Routes []BgpRibRoute `json:"openconfig-network-instance:route"`
}
