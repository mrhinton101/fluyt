package cue

import (
	"fmt"
	"log"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/encoding/yaml"
)

const cueSource = `
//ip address definition
#IPCidr: =~"^((25[0-5]|2[0-4][0-9]|1?[0-9]{1,2})\\.){3}(25[0-5]|2[0-4][0-9]|1?[0-9]{1,2})/(3[0-2]|[12]?[0-9])$"

//device entry schema
#Device: {
    name: string
    ip:   #IPCidr
    telemetry: [...string]
    tags: {
    region : string
    env: string
    }
    description?: string
    config?: [...string]
    pushmode?: "GNMI" | "Terraform" | "Pulumi"
}
//create inventory and add device name to device definition
#Inventory: {
    inventory: [k=string]: #Device & { name: k }
}
`

func CueLoad() {
	ctx := cuecontext.New()
	schema := ctx.CompileString(cueSource)
	inventorySchema := schema.LookupPath(cue.ParsePath("#Inventory"))

	yamlFile, err := yaml.Extract("inventory.yml", nil)
	if err != nil {
		log.Fatal(err)
	}
	yamlVal := ctx.BuildFile(yamlFile)

	unified := inventorySchema.Unify(yamlVal)
	// Validate unified value
	if err := unified.Validate(cue.Concrete(true)); err != nil {
		fmt.Println("YAML: NOT ok")
		fmt.Println(unified)
		log.Fatal(err)
	}

	fmt.Println("YAML: ok")

	devices := unified.LookupPath(cue.ParsePath("inventory"))
	// ip := unified.LookupPath(cue.ParsePath("inventory.device1.ip"))

	// fmt.Println(ip)
	iter, err := devices.Fields()
	if err != nil {
		log.Fatal(err)
	}

	for iter.Next() {
		deviceName := iter.Selector()
		deviceVal := iter.Value()

		ipVal := deviceVal.LookupPath(cue.ParsePath("ip"))
		ipStr, err := ipVal.String()
		if err != nil {
			log.Fatalf("Could not get IP for device %s: %v", deviceName, err)
		}

		fmt.Printf("Connecting to device %s at IP %s\n", deviceName, ipStr)

	}
}
