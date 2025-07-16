package main

import (
	"log/slog"

	"github.com/mrhinton101/fluyt/cmd/fluyt/commands"
	"github.com/mrhinton101/fluyt/internal/app/core/logger"
)

func main() {
	logger.ProgramLevel.Set(slog.LevelError)
	logfile := logger.InitLogger("fluytLogs.json")
	defer logfile.Close()

	commands.Execute()
}
