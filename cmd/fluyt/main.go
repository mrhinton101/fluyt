package main

import (
	"fmt"
	"log/slog"

	cmd "github.com/mrhinton101/fluyt/cmd/fluyt/commands"
	logger "github.com/mrhinton101/fluyt/internal/app/core"
	cue "github.com/mrhinton101/fluyt/internal/app/core/cueHandler"
)

func main() {
	logfile := logger.InitLogger("fluytLogs.json")
	defer logfile.Close()
	schemaDir := "../../schema/"
	ctx, schemaVals := cue.CueLoadSchemaDir(schemaDir)

	concreteInvVal := cue.CueLoadInventory(ctx, schemaVals, "./inventory.yml")
	// cue.CueLoadTelPaths(ctx, schemaVals)
	CueInputs := cue.CueGrabSubs(concreteInvVal)
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
