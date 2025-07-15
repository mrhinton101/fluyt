package cueHandler

import (
	"errors"
	"fmt"
	"log/slog"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	"cuelang.org/go/encoding/yaml"
	logger "github.com/mrhinton101/fluyt/internal/adapter/logger"
)

type DeviceSubsList struct {
	Devices     []DeviceSubPaths
	dedupTarget map[string]struct{}
}

type DeviceSubPaths struct {
	Name    string
	Address string
	Port    string
	Paths   []string
}

type concreteInv struct {
	inventory cue.Value
}

func initDeviceSubsList() *DeviceSubsList {
	return &DeviceSubsList{
		Devices:     []DeviceSubPaths{},
		dedupTarget: make(map[string]struct{}),
	}
}

func dedupList(inputSlice []string) []string {
	seen := make(map[string]struct{})
	outputSlice := make([]string, 0, len(inputSlice))
	for _, item := range inputSlice {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			outputSlice = append(outputSlice, item)
		} else {
			logger.SLogger(logger.LogEntry{
				Level:     slog.LevelError,
				Err:       errors.New("duplicate entry found"),
				Component: "Cue",
				Action:    "dedup list",
				Msg:       fmt.Sprintf("duplicate entry %s found in list", item),
				Target:    item,
			})
		}
	}
	return outputSlice
}

func (d *DeviceSubsList) Add(sub DeviceSubPaths) {

	if _, dupe := d.dedupTarget[sub.Address]; dupe {
		logger.SLogger(logger.LogEntry{
			Level:     slog.LevelError,
			Err:       errors.New("duplicate target found"),
			Component: "Cue",
			Action:    "add device subscription",
			Msg:       fmt.Sprintf("target %s already exists in the list", sub.Address),
			Target:    sub.Address,
		})
		return
	}
	d.dedupTarget[sub.Address] = struct{}{}
	cleanPathList := dedupList(sub.Paths)
	sub.Paths = cleanPathList
	d.Devices = append(d.Devices, sub)
}

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

func pathLookup(path string, schema cue.Value) (pathResults cue.Value) {
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

func loadYaml(ctx *cue.Context, yamlfile string) (yamlVal cue.Value) {
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

func loadInventory(ctx *cue.Context, schemaVals []cue.Value, invFile string) (concreteInvVal concreteInv) {
	// Load schema directory and all files in the package

	// use the first(and only) schema package. loadSchemaDir already confirms there is a single schema package
	schema := schemaVals[0]

	invSchema := pathLookup("#inventory", schema)

	invVal := loadYaml(ctx, invFile)

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
			Msg:       "error validating unified schema",
			Target:    "localhost",
		})
	}
	concreteInvVal = concreteInv{
		inventory: unifiedVal.LookupPath(cue.ParsePath("inventory"))}
	// Store the unified value in the concreteInv struct

	return concreteInvVal
}

func (inv *concreteInv) loadSubs() (DeviceSubsList *DeviceSubsList) {
	inventory := inv.inventory
	iter_inventory, err := inventory.Fields()
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
	DeviceSubsList = initDeviceSubsList()
	for iter_inventory.Next() {
		deviceName := iter_inventory.Selector()
		deviceVal := iter_inventory.Value()

		ipVal := deviceVal.LookupPath(cue.ParsePath("ip"))
		deviceNameStr := deviceName.String()
		ipStr, err := ipVal.String()
		telemetryPaths := deviceVal.LookupPath(cue.ParsePath("tel_paths"))
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
		iter_telemetry, err := telemetryPaths.List()
		telemPathList := []string{}
		for iter_telemetry.Next() {

			telemetryVal := iter_telemetry.Value()

			iter_telemetry_paths, err := telemetryVal.Fields()
			if err != nil {
				logger.SLogger(logger.LogEntry{
					Level:     slog.LevelError,
					Err:       err,
					Component: "Cue",
					Action:    "get telemetry paths",
					Msg:       fmt.Sprintf("Could not get telemetry paths for device %s: %v", deviceName, err),
					Target:    "localhost",
				})
			}
			for iter_telemetry_paths.Next() {
				telemetryPathVal := iter_telemetry_paths.Value()
				telemetryPathStr, err := telemetryPathVal.String()
				if err != nil {
					logger.SLogger(logger.LogEntry{
						Level:     slog.LevelError,
						Err:       err,
						Component: "Cue",
						Action:    "convert telemetry value to string",
						Msg:       fmt.Sprintf("Could not convert telemetry value to string for device %s: %v", deviceName, err),
						Target:    "localhost",
					})

				}
				telemPathList = append(telemPathList, telemetryPathStr)
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
			}

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
		device := DeviceSubPaths{
			Name:    deviceNameStr,
			Address: ipStr,
			Paths:   telemPathList,
		}
		DeviceSubsList.Add(device)
		// fmt.Printf("Connecting to device %s at IP %s\n", deviceName, ipStr)

	}
	// fmt.Println(DeviceSubsList.Devices)
	return DeviceSubsList
}

func LoadDeviceSubsList(schemaDir, invFile string) (DeviceSubsList *DeviceSubsList) {
	ctx, schemaVals := loadSchemaDir(schemaDir)
	concreteInv := loadInventory(ctx, schemaVals, invFile)
	DeviceSubsList = concreteInv.loadSubs()
	return DeviceSubsList
}
