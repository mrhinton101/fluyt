package cmd

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/mrhinton101/fluyt/internal/adapter/gnmiClient"
	"github.com/mrhinton101/fluyt/internal/adapter/logger"
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

		fmt.Println(CueInputs)
		creds := gnmiClient.Credentials{
			Username: username,
			Password: password}

		for _, target := range CueInputs.Devices {
			fmt.Println("capabilities")
			capResp := gnmiClient.Capabilities(target.DeviceInfo, creds)
			// fmt.Println(resp)
			for _, m := range capResp.SupportedModels {
				fmt.Printf("model: %s %s (rev %s)\n", m.Name, m.Organization, m.Version)
			}

			for _, enc := range capResp.SupportedEncodings {
				fmt.Println("encoding:", enc.String()) // enc is an enum
			}

			fmt.Println("gNMI versions:", capResp.GNMIVersion)
		}
		return nil
	},
}
