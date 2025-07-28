package web

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/go-cmp/cmp"
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
	Ribs      gnmi.BgpRibs
}

type SnapshotPageData struct {
	Title   string
	Results gnmi.BgpRibs
}

type DiffPageData struct {
	Title   string
	Results []TempRibSnap
}

type DiffResultsPageData struct {
	Title   string
	Results []DiffLine
}

type DiffLine struct {
	Line  string
	Class string // e.g. "add", "remove", or "neutral"
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

func (s *Server) DiffResultsHandler(w http.ResponseWriter, r *http.Request) {

	penulti := s.TempRibSnaps[len(s.TempRibSnaps)-2]

	ulti := s.TempRibSnaps[len(s.TempRibSnaps)-1]
	ignoredSuffixes := map[string]bool{
		"Timestamp":          true,
		"State.LastModified": true,
	}

	diff := cmp.Diff(penulti, ulti,
		cmp.FilterPath(func(p cmp.Path) bool {
			for suffix := range ignoredSuffixes {
				if strings.HasSuffix(p.String(), suffix) {
					return true
				}
			}
			return false
		}, cmp.Ignore()),
	)
	// fmt.Println(diff)

	if diff == "" {
		http.Error(w, "No differences found between snapshots", http.StatusNotFound)
		return
	}
	var diffLines []DiffLine
	for _, line := range strings.Split(diff, "\n") {
		fmt.Println("Processing line:", line)
		switch {
		case strings.Contains(line, "192.168.121."):
			diffLines = append(diffLines, DiffLine{Line: line, Class: "neutral"})
		case strings.HasPrefix(line, "+"):
			diffLines = append(diffLines, DiffLine{Line: line, Class: "add"})
		case strings.HasPrefix(line, "-"):
			diffLines = append(diffLines, DiffLine{Line: line, Class: "remove"})
		default:
			continue
		}
	}

	results := diffLines
	renderTemplate(w, "diffresults.html", DiffResultsPageData{
		Title:   fmt.Sprint("RIB Diff results for"),
		Results: results,
	})

}

func (s *Server) DiffHandler(w http.ResponseWriter, r *http.Request) {

	results := s.TempRibSnaps
	renderTemplate(w, "diff.html", DiffPageData{
		Title:   fmt.Sprint("RIB Differ for"),
		Results: results,
	})

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
		Ribs:      RibSnap}
	s.TempRibSnaps = append(s.TempRibSnaps, snapshot)
	fmt.Println(s.TempRibSnaps)
	fmt.Println("\n")

	// This method can be expanded to handle RIB differences as needed
}
