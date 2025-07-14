package main

import (
	"fmt"
	"log/slog"

	cmd "github.com/mrhinton101/fluyt/cmd/fluyt/commands"
	cue "github.com/mrhinton101/fluyt/internal/adapter/cueHandler"
	logger "github.com/mrhinton101/fluyt/internal/adapter/logger"
)

func main() {
	logfile := logger.InitLogger("fluytLogs.json")
	defer logfile.Close()
	schemaDir := "../../schema/"
	invFile := "./inventory.yml"
	CueInputs := cue.LoadCueInputs(schemaDir, invFile)
	defer func() {
		if r := recover(); r != nil {

			logger.SLogger(logger.LogEntry{
				Level:     slog.LevelError,
				Err:       fmt.Errorf("%v", r),
				Component: "runtime",
				Action:    "panic recovery",
				Msg:       "unexpected panic has occured",
				Target:    "localhost",
			})
		}
	}()
	cmd.SetCueInputs(CueInputs)
	cmd.Execute()

}
