package gnmi

import (
	"fmt"

	"github.com/openconfig/gnmi/proto/gnmi"
)

type GNMICapabilityResponse struct {
	Target   string
	Response gnmi.CapabilityResponse
}

type CleanCapabilityResponse struct {
	Target    string
	Encodings []string
	Models    []string
	Versions  string
}

func ValidateCapabilityResponse(target string, capResp gnmi.CapabilityResponse) (*GNMICapabilityResponse, error) {
	if len(capResp.GNMIVersion) == 0 && len(capResp.SupportedModels) == 0 && len(capResp.SupportedEncodings) == 0 {
		return nil, fmt.Errorf("no capabilities received for target %s", target)
	}
	return &GNMICapabilityResponse{
		Target:   target,
		Response: capResp,
	}, nil
}

func UnmarshalCapabilityResponse(capResp *GNMICapabilityResponse) (CleanCapabilityResponse, error) {
	if capResp == nil {
		fmt.Errorf("capability response is nil")
	}

	result := capResp.Response
	models := make([]string, len(result.SupportedModels))
	for i, m := range result.SupportedModels {
		models[i] = m.Name
	}

	encodings := make([]string, len(result.SupportedEncodings))
	for i, e := range result.SupportedEncodings {
		encodings[i] = e.String()
	}
	// for now flatten all capability fields into a map
	return CleanCapabilityResponse{
		Target:    capResp.Target,
		Encodings: encodings,
		Models:    models,
		Versions:  result.GNMIVersion,
	}, nil
}
