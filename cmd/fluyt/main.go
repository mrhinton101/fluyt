package main

import (
	cmd "github.com/mrhinton101/fluyt/cmd/fluyt/commands"
	"github.com/mrhinton101/fluyt/internal/adapter/cueHandler"
	logger "github.com/mrhinton101/fluyt/internal/adapter/logger"
)

func main() {
	logfile := logger.InitLogger("fluytLogs.json")
	defer logfile.Close()
	schemaDir := "../../schema/"
	invFile := "./inventory.yml"
	DeviceSubsList := cueHandler.LoadDeviceSubsList(schemaDir, invFile)

	cmd.DeviceSubsList(DeviceSubsList)
	cmd.Execute()

}
