package webui

import (
	"fmt"
	"log"
	"net/http"
)

func StartServer() {
	// Register routes
	fmt.Println("Registering routes...")
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/api/snapshot", SnapshotHandler)
	http.HandleFunc("/api/diff", DiffHandler)

	// Serve
	port := 8080
	log.Printf("Starting web UI on http://localhost:%d\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
