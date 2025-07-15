package main

import (
	"log/slog"

	cmd "github.com/mrhinton101/fluyt/cmd/fluyt/commands"
	"github.com/mrhinton101/fluyt/internal/adapter/cueHandler"
	logger "github.com/mrhinton101/fluyt/internal/adapter/logger"
)

func main() {
	logger.ProgramLevel.Set(slog.LevelError)
	logfile := logger.InitLogger("fluytLogs.json")
	defer logfile.Close()
	schemaDir := "../../schema/"
	invFile := "./inventory.yml"
	DeviceCapsList := cueHandler.LoadDeviceCapsList(schemaDir, invFile)

	cmd.DeviceCapsList(DeviceCapsList)
	cmd.Execute()

}
