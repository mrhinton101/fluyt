package cue

import (
	"errors"
	"fmt"
	"log/slog"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	"cuelang.org/go/encoding/yaml"
	"github.com/mrhinton101/fluyt/logger"
)

func CueLoadSchemaDir(schema_dir string) (schema_value []cue.Value) {
	ctx := cuecontext.New()
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
			Msg:       fmt.Sprintf("error loading schema %s", schema_dir),
			Target:    "localhost",
		})
	}

	schema_value, err := ctx.BuildInstances(instances)
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
	return schema_value

}

func CueLoadInventory(schemaDir string, inv_file string) {
	ctx := cuecontext.New()
	// schema := ctx.CompileString(cueSource)
	// topLevelSchema := schema.LookupPath(cue.ParsePath("#Inventory"))

	// Load schema directory and all files in the package
	schemaVals := CueLoadSchemaDir(schemaDir)
	// use the first(and only) schema package. CueLoadSchemaDir already confirms there is a single schema package
	schema := schemaVals[0]

	// Extract the `#Inventory` definition from schema
	topLevelKey := "#Inventory"
	topLevelSchema := schema.LookupPath(cue.ParsePath(topLevelKey))
	if !topLevelSchema.Exists() {
		logger.SLogger(logger.LogEntry{
			Level:     slog.LevelError,
			Err:       errors.New("failed to load top level schema key"),
			Component: "Cue",
			Action:    "load Cue top level key",
			Msg:       fmt.Sprintf("schema missing %s definition", topLevelKey),
			Target:    "localhost",
		})
	}

	// Load YAML inventory file
	yamlFile, err := yaml.Extract(inv_file, nil)
	if err != nil {
		logger.SLogger(logger.LogEntry{
			Level:     slog.LevelError,
			Err:       err,
			Component: "Cue",
			Action:    "load yaml file",
			Msg:       fmt.Sprintf("failed to load %v. Is the location correct", yamlFile),
			Target:    "localhost",
		})
	}
	yamlVal := ctx.BuildFile(yamlFile)

	unified := topLevelSchema.Unify(yamlVal)
	// Validate unified value
	if err := unified.Validate(cue.Concrete(true)); err != nil {
		logger.SLogger(logger.LogEntry{
			Level:     slog.LevelError,
			Err:       err,
			Component: "Cue",
			Action:    "unify schema and inventory",
			Msg:       "failed to validate unified schema",
			Target:    "localhost",
		})
	}

	devices := unified.LookupPath(cue.ParsePath("inventory"))

	iter, err := devices.Fields()
	if err != nil {
		logger.SLogger(logger.LogEntry{
			Level:     slog.LevelError,
			Err:       err,
			Component: "Cue",
			Action:    "load devices fields",
			Msg:       "failed to iterate over devices fields. are they defined?",
			Target:    "localhost",
		})
	}

	for iter.Next() {
		deviceName := iter.Selector()
		deviceVal := iter.Value()

		ipVal := deviceVal.LookupPath(cue.ParsePath("ip"))
		ipStr, err := ipVal.String()
		telemetryPaths := deviceVal.LookupPath(cue.ParsePath("telemetry"))
		if !telemetryPaths.Exists() {
			logger.SLogger(logger.LogEntry{
				Level:     slog.LevelError,
				Err:       errors.New("no telemetry paths defined"),
				Component: "Cue",
				Action:    "get telemetry paths",
				Msg:       fmt.Sprintf("Could not get telemetry paths for device %s: %v", deviceName, err),
				Target:    "localhost",
			})
		}
		iter, err := telemetryPaths.List()
		for iter.Next() {

			telemetryVal := iter.Value()
			telemetryStr, err := telemetryVal.String()
			if err != nil {
				logger.SLogger(logger.LogEntry{
					Level:     slog.LevelError,
					Err:       err,
					Component: "Cue",
					Action:    "get telemetry value",
					Msg:       fmt.Sprintf("Could not get telemetry value for device %s: %v", deviceName, err),
					Target:    "localhost",
				})
			}
			fmt.Printf("Device %s has telemetry %s\n", deviceName, telemetryStr)
		}

		if err != nil {
			logger.SLogger(logger.LogEntry{
				Level:     slog.LevelError,
				Err:       err,
				Component: "Cue",
				Action:    "find ip",
				Msg:       fmt.Sprintf("Could not get IP for device %s: %v", deviceName, err),
				Target:    "localhost",
			})
		}

		fmt.Printf("Connecting to device %s at IP %s\n", deviceName, ipStr)

	}
}
