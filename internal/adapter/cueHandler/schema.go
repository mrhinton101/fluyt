package cueHandler

import (
	"errors"
	"fmt"
	"log/slog"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	logger "github.com/mrhinton101/fluyt/internal/adapter/logger"
)

func loadSchemaDir(schema_dir string) (ctx *cue.Context, schemaVals []cue.Value) {
	ctx = cuecontext.New()
	instances := load.Instances([]string{schema_dir}, nil)

	if len(instances) > 1 {
		logger.SLogger(logger.LogEntry{
			Level:     slog.LevelError,
			Err:       errors.New("more than 1 schema was loaded"),
			Component: "Cue",
			Action:    "load Cue schema directory",
			Msg: `This function is designed to load a single schema package\n
			you likely loaded a schema with multiple packages.`,
			Target: "localhost",
		})
	}
	if instances[0].Err != nil {
		logger.SLogger(logger.LogEntry{
			Level:     slog.LevelError,
			Err:       instances[0].Err,
			Component: "Cue",
			Action:    "load Cue schema directory",
			Msg:       fmt.Sprintf("error loading schema %s. Check your schema syntax", schema_dir),
			Target:    "localhost",
		})
	}

	schemaVals, err := ctx.BuildInstances(instances)
	if err != nil {
		logger.SLogger(logger.LogEntry{
			Level:     slog.LevelError,
			Err:       err,
			Component: "Cue",
			Action:    "load Cue schema directory",
			Msg:       fmt.Sprintf("error loading schema %s", schema_dir),
			Target:    "localhost",
		})
	}
	logger.SLogger(logger.LogEntry{
		Level:     slog.LevelDebug,
		Err:       nil,
		Component: "Cue",
		Action:    "load Cue schema directory",
		Msg:       fmt.Sprintf("successfully loaded %s", schema_dir),
		Target:    "localhost",
	})
	// Leave schema value as a slice of Values for flexibility and remain idiomatic with Cue.
	// calling functions will still need to gather first entry
	return ctx, schemaVals
}
