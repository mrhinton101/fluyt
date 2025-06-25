package cmd

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/mrhinton101/fluyt/logger"
	"github.com/openconfig/gnmi/proto/gnmi"
	pb "github.com/openconfig/gnmi/proto/gnmi"
	"github.com/spf13/cobra"
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

var capabilitiesCmd = &cobra.Command{
	Use:   "capabilities",
	Short: "Determine the GNMI capabilities based on connection origin",
	Long: `The gNMI Capabilities RPC is used to discover the capabilities of a gNMI server. 
It allows a client to retrieve information about the gNMI version, supported data models (YANG modules), and supported encodings.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.SLogger(logger.LogEntry{
			Level:     slog.LevelDebug,
			Component: "cli command",
			Action:    "launch capabilities command",
			Msg:       "user selected capabilities",
			Target:    "localhost",
		})

		if addr == "" {
			err := errors.New("missing required flag or env var: addr")
			logger.SLogger(logger.LogEntry{
				Level:     slog.LevelError,
				Err:       err,
				Component: "cli command",
				Action:    "get address var",
				Msg:       "failed to find required flag",
				Target:    "localhost",
			})
			return err
		}

		if username == "" {
			err := errors.New("missing required flag or env var: username")
			logger.SLogger(logger.LogEntry{
				Level:     slog.LevelError,
				Err:       err,
				Component: "authentication",
				Action:    "get username var",
				Msg:       "failed to find required flag",
				Target:    "localhost",
			})
			return err
		}

		if password == "" {
			err := errors.New("missing required flag or env var: password")
			logger.SLogger(logger.LogEntry{
				Level:     slog.LevelError,
				Err:       err,
				Component: "authentication",
				Action:    "get password var",
				Msg:       "failed to find required flag",
				Target:    "localhost",
			})
			return err
		}

		logger.SLogger(logger.LogEntry{
			Level:     slog.LevelDebug,
			Component: "GNMI Client",
			Action:    "launch GNMI Client",
			Msg:       fmt.Sprintf("Successfully authenticated. Launching GNMI client for %s on %s", username, addr),
			Target:    addr,
		})

		client, err := NewGNMIClient(addr, 3*time.Second, username, password)
		if err != nil {
			logger.SLogger(logger.LogEntry{
				Level:     slog.LevelError,
				Err:       err,
				Component: "authentication",
				Action:    "launch GNMI Client",
				Msg:       fmt.Sprintf("Authentication failed for user: %s on device: %s", username, addr),
				Target:    addr,
			})
			return err
		}
		defer client.Conn.Close()

		logger.SLogger(logger.LogEntry{
			Level:     slog.LevelInfo,
			Component: "GNMI Client",
			Action:    "launch GNMI Client",
			Msg:       fmt.Sprintf("successfully launched GNMI Client for user: %s on device: %s", username, addr),
			Target:    addr,
		})

		resp, err := client.Capabilities()
		if err != nil {
			logger.SLogger(logger.LogEntry{
				Level:     slog.LevelError,
				Err:       err,
				Component: "GNMI Client",
				Action:    "get capabilities",
				Msg:       fmt.Sprintf("failed to get capabilities on device: %s", addr),
				Target:    addr,
			})
			return err
		}

		fmt.Println(resp)
		return nil
	},
}

func (c *GNMIClient) Capabilities() (*pb.CapabilityResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	md := metadata.Pairs(
		"username", c.Username,
		"password", c.Password,
	)
	ctx = metadata.NewOutgoingContext(ctx, md)

	req := &gnmi.CapabilityRequest{}
	resp, err := c.GNMI.Capabilities(ctx, req)
	if err == nil {
		slog.Info("capabilities retrieved", "target", c.Target)
	}
	return resp, err
}
