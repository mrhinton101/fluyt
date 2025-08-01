package ports

import "github.com/mrhinton101/fluyt/domain/model"

type CueHandler interface {
	LoadDeviceSubs(schemaDir, invFile string) (*model.DeviceSubsList, error)
	LoadDeviceList(schemaDir, invFile string) (*model.DeviceList, error)
}
