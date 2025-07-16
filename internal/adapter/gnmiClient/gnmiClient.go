package gnmiClient

import (
	"context"
	"fmt"
	"log"
	"log/slog"

	"github.com/mrhinton101/fluyt/domain/cue"
	"github.com/mrhinton101/fluyt/domain/gnmi"
	"github.com/mrhinton101/fluyt/internal/app/core/logger"
	"github.com/mrhinton101/fluyt/internal/app/ports"
	"github.com/openconfig/gnmic/pkg/api"
)

type GNMIClientImpl struct {
	Name    string
	Address string
	Port    string
}

func NewGNMIClient(device cue.Device) ports.GNMIClient {
	return &GNMIClientImpl{
		Name:    device.Name,
		Address: fmt.Sprintf("%s:6030", device.Address),
		Port:    device.Port,
	}
}

func (c *GNMIClientImpl) Capabilities() (map[string]interface{}, error) {
	// Run the actual gNMI Capabilities RPC and get result
	tg, err := api.NewTarget(
		api.Name(c.Name),
		api.Address(c.Address),
		api.Username("admin"),
		api.Password("admin"),
		api.SkipVerify(true),
		api.Insecure(true),
	)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create a gNMI client
	err = tg.CreateGNMIClient(ctx)
	if err != nil {
		logger.SLogger(logger.LogEntry{
			Level:     slog.LevelError,
			Component: "gnmiClient",
			Msg:       fmt.Sprintf("failed to create gNMI client for target"),
			Err:       err,
			Target:    c.Name,
		})
		log.Fatal(err)
	}
	defer tg.Close()

	// send a gNMI capabilities request to the created target
	capResp, err := tg.Capabilities(ctx)
	fmt.Println("capresp:")
	if err != nil {
		logger.SLogger(logger.LogEntry{
			Level:     slog.LevelError,
			Component: "gnmiClient",
			Msg:       "failed to get capabilities from target",
			Err:       err,
			Target:    c.Name,
		})
		log.Fatal(err)
	}

	resp, err := gnmi.ValidateCapabilityResponse(c.Address, *capResp)
	if err != nil {
		return nil, err
	}

	// convert to map[string]interface{}
	result, err := gnmi.UnmarshalCapabilityResponse(resp)
	if err != nil {
		return nil, err
	}

	// for now flatten all capability fields into a map
	return map[string]interface{}{
		"target":    result.Target,
		"encodings": result.Encodings,
		"models":    result.Models,
		"versions":  result.Versions,
	}, nil
}

func (c *GNMIClientImpl) Close() error {
	return nil
}
