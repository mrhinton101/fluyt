package cmd

import (
	"log/slog"
	"os"

	"github.com/mrhinton101/fluyt/cue"
	"github.com/mrhinton101/fluyt/logger"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "fluyt",
	Short: "Tool to interact with grpc for network devices",
	Long:  `This tool will be used to interact via RPCs over GNMI. RPCs supported will be Get, Subscribe and Set`,
}

var (
	username string
	addr     string
	password string
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logger.SLogger(logger.LogEntry{
			Level:     slog.LevelError,
			Err:       err,
			Component: "cli",
			Action:    "execute",
			Msg:       "fatal error in CLI",
			Target:    "self",
		})
	}
}

var CueInputs *cue.DeviceSubsList

func SetCueInputs(l *cue.DeviceSubsList) { CueInputs = l }
func init() {
	rootCmd.PersistentFlags().StringVar(&addr, "addr", os.Getenv("GNMI_ADDR"), "Target device address (or GNMI_ADDR)")
	rootCmd.PersistentFlags().StringVar(&username, "username", os.Getenv("GNMI_USERNAME"), "Username (or GNMI_USERNAME)")
	rootCmd.PersistentFlags().StringVar(&password, "password", os.Getenv("GNMI_PASSWORD"), "Password (or GNMI_PASSWORD)")
	rootCmd.AddCommand(capabilitiesCmd)
}
