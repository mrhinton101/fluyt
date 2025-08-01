package avd

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/mrhinton101/fluyt/domain/model"

	"github.com/mrhinton101/fluyt/internal/app/core/utils"
	"gopkg.in/yaml.v3"
)

func LoadAvdDeviceConfig(reader utils.Reader, path string) (*AvdDeviceConfig, error) {
	data, err := reader.Read(path)
	if err != nil {
		return nil, err
	}
	var cfg AvdDeviceConfig
	switch filepath.Ext(path) {
	case ".json":
		err = json.Unmarshal(data, &cfg)
	case ".yaml", ".yml":
		err = yaml.Unmarshal(data, &cfg)
	default:
		return nil, fmt.Errorf("unsupported file format: %s", path)
	}

	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func AvdToModel(path string) (model.DeviceConfig, error) {
	fmt.Printf("avd path: %s\n", path)
	reader := utils.LocalReader{}
	avdConfigs, err := LoadAvdDeviceConfig(reader, path)
	fmt.Println("avd Configs:\n", avdConfigs)
	if err != nil {
		return model.DeviceConfig{}, err
	}
	DeviceConfig := model.DeviceConfig{
		Routing: model.DeviceConfigRouting{Bgp: avdConfigs.Bgp},
	}

	return DeviceConfig, nil
}
