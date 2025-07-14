package gnmiClient

import (
	"context"
	"fmt"
	"log/slog"

	cue "github.com/mrhinton101/fluyt/internal/adapter/cueHandler"
	logger "github.com/mrhinton101/fluyt/internal/adapter/logger"
	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/openconfig/gnmic/pkg/api"
)

type Credentials struct {
	Username string
	Password string
}

func Capabilities(gnmiTarget cue.DeviceSubPaths, creds Credentials) (capResp *gnmi.CapabilityResponse) {
	// create a target
	tg, err := api.NewTarget(
		api.Name(gnmiTarget.Name),
		api.Address(fmt.Sprintf("%s:6030", gnmiTarget.Address)),
		api.Username(creds.Username),
		api.Password(creds.Password),
		api.Insecure(true),
		api.SkipVerify(true),
	)
	if err != nil {
		logger.SLogger(logger.LogEntry{
			Level:     slog.LevelError,
			Err:       err,
			Component: "gnmiClient",
			Action:    "create target",
			Msg:       "failed to create gNMI target",
			Target:    gnmiTarget.Name,
		})
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create a gNMI client
	err = tg.CreateGNMIClient(ctx)
	if err != nil {
		logger.SLogger(logger.LogEntry{
			Level:     slog.LevelError,
			Err:       err,
			Component: "gnmiClient",
			Action:    "create gNMI client",
			Msg:       "failed to create gNMI client",
			Target:    gnmiTarget.Name,
		})
	}
	defer tg.Close()

	// send a gNMI capabilities request to the created target
	capResp, err = tg.Capabilities(ctx)
	if err != nil {
		logger.SLogger(logger.LogEntry{
			Level:     slog.LevelError,
			Err:       err,
			Component: "gnmiClient",
			Action:    "get capabilities",
			Msg:       "failed to get gNMI capabilities",
			Target:    gnmiTarget.Name,
		})
	}
	// response, err = prototext.Marshal(capResp)
	if err != nil {
		logger.SLogger(logger.LogEntry{
			Level:     slog.LevelError,
			Err:       err,
			Component: "gnmiClient",
			Action:    "marshal capabilities response",
			Msg:       "failed to marshal gNMI capabilities response",
			Target:    gnmiTarget.Name,
		})
	}
	return capResp
}
