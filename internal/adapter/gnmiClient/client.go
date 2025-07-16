package gnmiClient

import (
	"github.com/mrhinton101/fluyt/domain/gnmi"
	"github.com/mrhinton101/fluyt/internal/app/ports"
)

type GNMIClientImpl struct {
	target string
	// whatever underlying gnmi-specific fields
}

func NewGNMIClient(target string) ports.GNMIClient {
	return &GNMIClientImpl{target: target}
}

func (c *GNMIClientImpl) Connect() error {
	// actual gNMI client dial logic here
	return nil
}

func (c *GNMIClientImpl) Capabilities() (map[string][]string, error) {
	// Run the actual gNMI Capabilities RPC and get result
	rawModels := []string{"openconfig-interfaces", "openconfig-bgp"}
	rawEncodings := []string{"JSON_IETF", "PROTO"}
	rawVersions := []string{"0.7.0"}

	result, err := gnmi.ValidateCapabilityResponse(c.target, rawEncodings, rawModels, rawVersions)
	if err != nil {
		return nil, err
	}

	// for now flatten all capability fields into a map
	return map[string][]string{
		"encodings": result.Encodings,
		"models":    result.Models,
		"versions":  result.Versions,
	}, nil
}

func (c *GNMIClientImpl) Close() error {
	return nil
}

func (c *GNMIClientImpl) Target() string {
	return c.target
}
