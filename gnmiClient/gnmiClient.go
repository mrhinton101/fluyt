package gnmiClient

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/mrhinton101/fluyt/logger"
	pb "github.com/openconfig/gnmi/proto/gnmi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type GNMIClient struct {
	Conn     *grpc.ClientConn
	GNMI     pb.GNMIClient
	Target   string
	Timeout  time.Duration
	Username string
	Password string
}

func NewGNMIClient(target string, timeout time.Duration, username string, password string) (*GNMIClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, target, grpc.WithInsecure())
	if err != nil {
		// errLogger(err, "create-grpc-conn", "create grpc DialContext", "failed to create grpc connection", "local")
		os.Exit(1)
	}
	slog.Info("successfully created grpc connection")

	return &GNMIClient{
		Conn:     conn,
		GNMI:     pb.NewGNMIClient(conn),
		Target:   target,
		Timeout:  timeout,
		Username: username,
		Password: password,
	}, nil
}

func (c *GNMIClient) Capabilities() (*pb.CapabilityResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	md := metadata.Pairs(
		"username", c.Username,
		"password", c.Password,
	)
	ctx = metadata.NewOutgoingContext(ctx, md)

	req := &pb.CapabilityRequest{}
	resp, err := c.GNMI.Capabilities(ctx, req)
	if err == nil {
		slog.Info("capabilities retrieved", "target", c.Target)
	}
	return resp, err
}

func (c *GNMIClient) Subscribe(query string, timeoutWithUnit string) error {
	timeout, err := time.ParseDuration(timeoutWithUnit)
	if err != nil {
		logger.SLogger(logger.LogEntry{
			Level:     slog.LevelError,
			Component: "gnmi client",
			Action:    "parse subscribe duration",
			Msg:       fmt.Sprintf("unable to parse timeout duration: %s", timeoutWithUnit),
			Err:       err,
			Target:    "localhost",
		})
	}
	if timeout <= 0 {
		logger.SLogger(logger.LogEntry{
			Level:     slog.LevelError,
			Component: "gnmi client",
			Action:    "validate subscribe duration",
			Msg:       fmt.Sprintf("invalid timeout duration: %s", timeoutWithUnit),
			Target:    "localhost",
		})
		return fmt.Errorf("invalid timeout duration: %s", timeoutWithUnit)
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	md := metadata.Pairs(
		"username", c.Username,
		"password", c.Password,
	)
	ctx = metadata.NewOutgoingContext(ctx, md)

	req := &pb.CapabilityRequest{}
	resp, err := c.GNMI.Capabilities(ctx, req)
	if err == nil {
		slog.Info("capabilities retrieved", "target", c.Target)
	}
	return resp, err
}
