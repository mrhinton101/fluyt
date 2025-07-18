package cueHandler

import (
	"errors"
	"fmt"
	"log/slog"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	"cuelang.org/go/encoding/yaml"
	cueInternal "github.com/mrhinton101/fluyt/domain/cue"
	"github.com/mrhinton101/fluyt/internal/app/core/logger"
	"github.com/mrhinton101/fluyt/internal/app/ports"
)

type CueHandler struct{}

func NewCueHandler() ports.CueHandler {
	return &CueHandler{}
}

func (h *CueHandler) LoadDeviceSubs(schemaDir, invFile string) (*cueInternal.DeviceSubsList, error) {
	ctx, schemaVals := h.loadSchemaDir(schemaDir)
	concrete := h.loadInventory(ctx, schemaVals, invFile)
	return concrete.LoadSubs()
}

func (h *CueHandler) LoadDeviceList(schemaDir, invFile string) (*cueInternal.DeviceList, error) {
	ctx, schemaVals := h.loadSchemaDir(schemaDir)
	concrete := h.loadInventory(ctx, schemaVals, invFile)
	return concrete.LoadDevices()
}

// ------------------- Internals -------------------

func (h *CueHandler) loadSchemaDir(schemaDir string) (*cue.Context, []cue.Value) {
	ctx := cuecontext.New()
	instances := load.Instances([]string{schemaDir}, nil)

	if len(instances) > 1 {
		logger.SLogger(logger.LogEntry{
			Level:     slog.LevelError,
			Err:       errors.New("more than 1 schema was loaded"),
			Component: "Cue",
			Action:    "load Cue schema directory",
			Msg:       "Expected a single CUE package in schema dir",
			Target:    "localhost",
		})
	}

	if instances[0].Err != nil {
		logger.SLogger(logger.LogEntry{
			Level:     slog.LevelError,
			Err:       instances[0].Err,
			Component: "Cue",
			Action:    "load Cue schema directory",
			Msg:       fmt.Sprintf("error loading schema %s. Check your schema syntax", schemaDir),
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
			Msg:       fmt.Sprintf("error loading schema %s", schemaDir),
			Target:    "localhost",
		})
	}

	return ctx, schemaVals
}

func (h *CueHandler) loadInventory(ctx *cue.Context, schemaVals []cue.Value, invFile string) cueInternal.ConcreteInv {
	schema := schemaVals[0]
	invSchema := schema.LookupPath(cue.ParsePath("#inventory"))
	yamlVal := h.loadYaml(ctx, invFile)

	unifiedVal := invSchema.Unify(yamlVal)

	if err := unifiedVal.Validate(cue.Concrete(true)); err != nil {
		logger.SLogger(logger.LogEntry{
			Level:     slog.LevelError,
			Err:       err,
			Component: "Cue",
			Action:    "validate schema + inventory",
			Msg:       "validation failed",
			Target:    "localhost",
		})
	}

	return cueInternal.ConcreteInv{
		Inventory: unifiedVal.LookupPath(cue.ParsePath("inventory")),
	}
}

func (h *CueHandler) loadYaml(ctx *cue.Context, invFile string) cue.Value {
	yamlFile, err := yaml.Extract(invFile, nil)
	if err != nil {
		logger.SLogger(logger.LogEntry{
			Level:     slog.LevelError,
			Err:       err,
			Component: "Cue",
			Action:    "load yaml file",
			Msg:       fmt.Sprintf("could not load %s", invFile),
			Target:    "localhost",
		})
	}
	return ctx.BuildFile(yamlFile)
}
