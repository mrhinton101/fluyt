package web

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/mrhinton101/fluyt/domain/cue"
	"github.com/mrhinton101/fluyt/domain/gnmi"
	"github.com/mrhinton101/fluyt/internal/adapter/gnmiClient"
	"github.com/mrhinton101/fluyt/internal/app/usecase"
)

type Server struct {
	Devices      cue.DeviceList
	TempRibSnaps []TempRibSnap
}

type TempRibSnap struct {
	Devices   cue.DeviceList
	Timestamp string
	ribs      gnmi.BgpRibs
}

type SnapshotPageData struct {
	Title   string
	Results gnmi.BgpRibs
}

// func RibHandler(w http.ResponseWriter, r *http.Request) {
// 	// TODO: parse from/to and run diff
// 	fmt.Fprintln(w, "Diff results here")
// }

// func SnapshotHandler(w http.ResponseWriter, r *http.Request) {

// 	bgpResults := usecase.CollectBgpRib(devices, gnmiClient.NewGNMIClient)

// 	data := SnapshotPageData{
// 		Title:   "RIB Snapshot",
// 		Results: bgpResults,
// 	}

// 	renderTemplate(w, "snapshot.html", data)
// }

func DiffHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: parse from/to and run diff
	fmt.Fprintln(w, "Diff results here")
}

func (s *Server) SnapshotHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Target from chi.URLParam: %q\n", chi.URLParam(r, "device"))
	target := chi.URLParam(r, "device")

	var devices *cue.DeviceList
	fmt.Printf("Target device: %s\n", target)
	switch target {
	case "all":
		devices = s.Devices.All()
	default:
		// Get specific device by name
		var found bool
		devices, found = s.Devices.GetByName(target)

		if !found {
			http.Error(w, "device not found", http.StatusNotFound)
			return
		}

	}
	var snapshotDevices cue.DeviceList
	snapshotDevices = *devices

	results := usecase.CollectBgpRib(devices, gnmiClient.NewGNMIClient)
	if len(results) != 0 {
		s.AddRibDiff(snapshotDevices, results)
	}
	renderTemplate(w, "snapshot.html", SnapshotPageData{
		Title:   fmt.Sprintf("RIB Snapshot for %s", target),
		Results: results,
	})

}

func (s *Server) SearchHandler(w http.ResponseWriter, r *http.Request) {
	device := r.URL.Query().Get("device")
	if device == "" {
		http.Error(w, "missing device name", http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/api/%s/snapshot", device), http.StatusSeeOther)
}

func (s *Server) AddRibDiff(DeviceList cue.DeviceList, RibSnap gnmi.BgpRibs) {
	// Placeholder for adding RIB diff logic
	fmt.Println("Adding RIB diff logic...")
	timestamp := time.Now().Format(time.RFC3339)
	snapshot := TempRibSnap{
		Devices:   DeviceList,
		Timestamp: timestamp,
		ribs:      RibSnap}
	s.TempRibSnaps = append(s.TempRibSnaps, snapshot)
	fmt.Println(s.TempRibSnaps)
	fmt.Println("\n")

	// This method can be expanded to handle RIB differences as needed
}
