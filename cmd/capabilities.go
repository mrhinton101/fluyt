package cmd

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/mrhinton101/fluyt/gnmiClient"
	"github.com/mrhinton101/fluyt/logger"
	"github.com/spf13/cobra"
)

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
			Component: "gnmi client",
			Action:    "launch GNMI Client",
			Msg:       fmt.Sprintf("Successfully authenticated. Launching GNMI client for %s on %s", username, addr),
			Target:    addr,
		})

		client, err := gnmiClient.NewGNMIClient(addr, 3*time.Second, username, password)
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
			Component: "gnmi client",
			Action:    "launch GNMI Client",
			Msg:       fmt.Sprintf("successfully launched GNMI Client for user: %s on device: %s", username, addr),
			Target:    addr,
		})

		resp, err := client.Capabilities()
		if err != nil {
			logger.SLogger(logger.LogEntry{
				Level:     slog.LevelError,
				Err:       err,
				Component: "gnmi client",
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
