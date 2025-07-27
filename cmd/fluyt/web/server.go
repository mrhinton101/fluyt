package web

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mrhinton101/fluyt/domain/cue"
)

func StartServer(devices *cue.DeviceList) {
	server := Server{Devices: *devices}

	fmt.Println("Registering routes...")
	r := chi.NewRouter()

	// Routes
	r.Get("/", IndexHandler)
	r.Get("/api/{device}/snapshot", server.SnapshotHandler)
	r.Get("/search", server.SearchHandler)
	// r.Get("/api/diff", server.DiffHandler)

	r.Get("/test/{device}", func(w http.ResponseWriter, r *http.Request) {
		target := chi.URLParam(r, "device")
		fmt.Fprintf(w, "Test device: %s\n", target)
	})

	port := 8080
	fmt.Printf("Starting web UI on http://localhost:%d\n", port)
	fmt.Println("Router attached to server correctly")
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), r)
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

// package web

// import (
// 	"fmt"
// 	"log"
// 	"net/http"

// 	"github.com/go-chi/chi/v5"
// 	"github.com/mrhinton101/fluyt/domain/cue"
// )

// func StartServer(devices *cue.DeviceList) {
// 	server := Server{Devices: *devices}
// 	r := chi.NewRouter()

// 	r.Get("/test/{device}", func(w http.ResponseWriter, r *http.Request) {
// 		device := chi.URLParam(r, "device")
// 		fmt.Printf("Extracted device: %q\n", device)
// 		fmt.Fprintf(w, "Test device: %s\n", device)
// 	})
// 	r.Get("/api/{device}/snapshot", server.SnapshotHandler)

// 	port := 8080
// 	fmt.Printf("Starting server on http://localhost:%d\n", port)
// 	err := http.ListenAndServe(fmt.Sprintf(":%d", port), r)
// 	if err != nil {
// 		log.Fatalf("Server failed: %v", err)
// 	}
// }
