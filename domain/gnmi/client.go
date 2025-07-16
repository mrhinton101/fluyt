package gnmi

import (
	"fmt"
)

type CapabilityResult struct {
	Target    string
	Encodings []string
	Models    []string
	Versions  []string
}

func ValidateCapabilityResponse(target string, encodings, models, versions []string) (*CapabilityResult, error) {
	if len(encodings) == 0 && len(models) == 0 && len(versions) == 0 {
		return nil, fmt.Errorf("no capabilities received for target %s", target)
	}
	return &CapabilityResult{
		Target:    target,
		Encodings: encodings,
		Models:    models,
		Versions:  versions,
	}, nil
}
