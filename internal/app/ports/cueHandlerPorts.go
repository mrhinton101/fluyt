package ports

import "github.com/mrhinton101/fluyt/domain/cue"

type CueHandler interface {
	LoadDeviceSubs(schemaDir, invFile string) (*cue.DeviceSubsList, error)
	LoadDeviceList(schemaDir, invFile string) (*cue.DeviceList, error)
}
