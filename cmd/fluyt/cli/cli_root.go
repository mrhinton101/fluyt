package cli

import (
	"fmt"
	"log"

	"github.com/mrhinton101/fluyt/internal/adapter/cueHandler"
	"github.com/mrhinton101/fluyt/internal/adapter/gnmiClient"
	"github.com/mrhinton101/fluyt/internal/app/usecase"
)

var (
	schemaDir = "../../schema/"
	invFile   = "./inventory.yml"
)

func Execute() {
	// Instantiate the CueHandler
	cue := cueHandler.NewCueHandler()

	// Call the method on the instance
	devices, err := cue.LoadDeviceList(schemaDir, invFile)
	if err != nil {
		log.Fatalf("failed to load devices: %v", err)
	}

	// Use the devices in the usecase layer
	results := usecase.CollectCapabilities(devices, gnmiClient.NewGNMIClient)
	fmt.Println("usecase goroutine finished")
	// Print the results
	for _, r := range results {
		fmt.Printf("Target: %s\n", r.Target)
		fmt.Printf("Encodings:\n")
		for _, encoding := range r.Encodings {
			fmt.Printf("	- %v\n", encoding)
		}
		fmt.Printf("Models:\n")
		for _, model := range r.Models {
			fmt.Printf("	- %v\n", model)
		}

		fmt.Printf("Versions: %v\n\n", r.Versions)
	}
	bgpResults := usecase.CollectBgpRib(devices, gnmiClient.NewGNMIClient)
	fmt.Println("RIB goroutine finished")

	for _, bgpRib := range bgpResults {
		for _, route := range bgpRib {
			fmt.Printf("Prefix: %s, \n", route.Prefix)
		}
	}
}
