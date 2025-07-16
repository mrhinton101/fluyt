package gnmi

import (
	"fmt"

	"github.com/openconfig/gnmi/proto/gnmi"
)

type RawCapabilityResult struct {
	Target    string
	Encodings []gnmi.Encoding
	Models    []*gnmi.ModelData
	Versions  string
}

func ValidateCapabilityResponse(target string, capResp gnmi.CapabilityResponse) (*RawCapabilityResult, error) {
	if len(capResp.GNMIVersion) == 0 && len(capResp.SupportedModels) == 0 && len(capResp.SupportedEncodings) == 0 {
		return nil, fmt.Errorf("no capabilities received for target %s", target)
	}
	return &RawCapabilityResult{
		Target:    target,
		Encodings: capResp.SupportedEncodings,
		Models:    capResp.SupportedModels,
		Versions:  capResp.GNMIVersion,
	}, nil
}
