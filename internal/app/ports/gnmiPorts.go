package ports

import (
	"context"

	"github.com/mrhinton101/fluyt/domain/gnmi"
)

type GNMIClient interface {
	Capabilities(context.Context) (map[string]interface{}, error)
	Init(context.Context) error
	GetAddress() string
	Close()
	GetBgpRibs(context.Context) (gnmi.BgpRibs, error)
}
