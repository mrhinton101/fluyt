package cmd

import (
	"log/slog"
	"os"

<<<<<<< HEAD:cmd/fluyt/commands/root.go
	"github.com/mrhinton101/fluyt/internal/adapter/cueHandler"
	"github.com/mrhinton101/fluyt/internal/adapter/logger"
=======
	"github.com/mrhinton101/fluyt/cue"
	"github.com/mrhinton101/fluyt/logger"
>>>>>>> origin/main:cmd/root.go
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

<<<<<<< HEAD:cmd/fluyt/commands/root.go
// CueInputs is a global variable that holds the device subscription list
var CueSubsInputs *cueHandler.DeviceSubsList

// CueInputs is a global variable that holds the device subscription list
var CueCapsInputs *cueHandler.DeviceList

// DeviceSubsList is a function that sets the CueInputs variable
func DeviceSubsList(DeviceSubsList *cueHandler.DeviceSubsList) { CueSubsInputs = DeviceSubsList }

// DeviceCapsList is a function that sets the CueInputs variable
func DeviceCapsList(DeviceCapsList *cueHandler.DeviceList) { CueCapsInputs = DeviceCapsList }

=======
var CueInputs *cue.DeviceSubsList

func SetCueInputs(l *cue.DeviceSubsList) { CueInputs = l }
>>>>>>> origin/main:cmd/root.go
func init() {
	rootCmd.PersistentFlags().StringVar(&addr, "addr", os.Getenv("GNMI_ADDR"), "Target device address (or GNMI_ADDR)")
	rootCmd.PersistentFlags().StringVar(&username, "username", os.Getenv("GNMI_USERNAME"), "Username (or GNMI_USERNAME)")
	rootCmd.PersistentFlags().StringVar(&password, "password", os.Getenv("GNMI_PASSWORD"), "Password (or GNMI_PASSWORD)")
	rootCmd.AddCommand(capabilitiesCmd)
}
