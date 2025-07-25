package web

import (
	"fmt"
	"net/http"

	"github.com/mrhinton101/fluyt/domain/gnmi"
	"github.com/mrhinton101/fluyt/internal/adapter/cueHandler"
	"github.com/mrhinton101/fluyt/internal/adapter/gnmiClient"
	"github.com/mrhinton101/fluyt/internal/app/usecase"
)

var (
	schemaDir = "../../schema/"
	invFile   = "./inventory.yml"
)

type SnapshotPageData struct {
	Title   string
	Results gnmi.BgpRibs
}

func SnapshotHandler(w http.ResponseWriter, r *http.Request) {
	cue := cueHandler.NewCueHandler()

	devices, err := cue.LoadDeviceList(schemaDir, invFile)
	if err != nil {
		http.Error(w, "failed to load devices", 500)
		return
	}

	bgpResults := usecase.CollectBgpRib(devices, gnmiClient.NewGNMIClient)

	data := SnapshotPageData{
		Title:   "RIB Snapshot",
		Results: bgpResults,
	}

	renderTemplate(w, "snapshot.html", data)
}

func DiffHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: parse from/to and run diff
	fmt.Fprintln(w, "Diff results here")
}
