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

func CueLoadSchemaDir(schema_dir string) (ctx *cue.Context, schemaVals []cue.Value) {
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
			Msg:       fmt.Sprintf("error loading schema %s", schema_dir),
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

func cuePathLookup(path string, schema cue.Value) (pathResults cue.Value) {
	// Extract the `#Inventory` definition from schema
	pathResults = schema.LookupPath(cue.ParsePath(path))
	if !pathResults.Exists() {
		logger.SLogger(logger.LogEntry{
			Level:     slog.LevelError,
			Err:       errors.New("failed to load top level schema key"),
			Component: "Cue",
			Action:    "load Cue top level key",
			Msg:       fmt.Sprintf("schema missing %s definition", path),
			Target:    "localhost",
		})
	}
	return pathResults

}

func cueLoadYaml(ctx *cue.Context, yamlfile string) (yamlVal cue.Value) {
	// Load YAML inventory file
	yamlFile, err := yaml.Extract(yamlfile, nil)
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
	yamlVal = ctx.BuildFile(yamlFile)
	return yamlVal

}

func CueLoadTelPaths(ctx *cue.Context, schemaVals []cue.Value) {
	schema := schemaVals[0]

	telSchemaVal := cuePathLookup("#telemetry_paths", schema)

	fmt.Println(telSchemaVal)
	// iter, err := telSchemaVal.Fields()

	// for iter.Next() {
	// 	telemetryVal := iter.Value()
	// 	telemetryStr, err := telemetryVal.String()
	// 	if err != nil {
	// 		logger.SLogger(logger.LogEntry{
	// 			Level:     slog.LevelError,
	// 			Err:       err,
	// 			Component: "Cue",
	// 			Action:    "get telemetry value",
	// 			Msg:       fmt.Sprintf("Could not get telemetry value for device %v", err),
	// 			Target:    "localhost",
	// 		})
	// 	}
	// 	fmt.Printf("Device %s has telemetry %s\n",  telemetryStr)
	// }

}

func CueLoadInventory(ctx *cue.Context, schemaVals []cue.Value, invFile string) (concreteInvVal cue.Value) {
	// Load schema directory and all files in the package

	// use the first(and only) schema package. CueLoadSchemaDir already confirms there is a single schema package
	schema := schemaVals[0]

	invSchema := cuePathLookup("#inventory", schema)

	invVal := cueLoadYaml(ctx, invFile)

	unifiedVal := invSchema.Unify(invVal)
	// fmt.Println("Unified CUE value:")
	// fmt.Println(unifiedVal)
	// Validate unified value. all values must be defined
	if err := unifiedVal.Validate(cue.Concrete(true)); err != nil {
		logger.SLogger(logger.LogEntry{
			Level:     slog.LevelError,
			Err:       err,
			Component: "Cue",
			Action:    "unify schema and inventory",
			Msg:       "failed to validate unified schema",
			Target:    "localhost",
		})
	}
	concreteInvVal = unifiedVal.LookupPath(cue.ParsePath("inventory"))
	return concreteInvVal
}

func CueMatchTags(concreteInvVal cue.Value) {
	iter, err := concreteInvVal.Fields()
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
