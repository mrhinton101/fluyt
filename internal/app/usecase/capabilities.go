package usecase

import (
	cmd "github.com/mrhinton101/fluyt/cmd/fluyt/commands"
	"github.com/mrhinton101/fluyt/internal/adapter/cueHandler"
)

func RunCapabilities(schemaDir, invFile string) (err error) {

	// Load the device capabilities list
	DeviceCapsList := cueHandler.LoadDeviceCapsList(schemaDir, invFile)

	// Execute the command to display the device capabilities
	cmd.DeviceCapsList(DeviceCapsList)
	cmd.Execute()
}
