package ports

import "github.com/mrhinton101/fluyt/domain/cue"

type CueAdapter interface {
	LoadDeviceSubs(schemaDir, invFile string) (*cue.DeviceSubsList, error)
}
