package main

import (
	"context"
	"fmt"
	"log"

	"github.com/openconfig/gnmic/pkg/api"
)

// type BgpRibRoute struct {
// 	Origin string `json:"origin"`
// 	PathID int    `json:"path-id"`
// 	Prefix string `json:"prefix"`
// 	State  struct {
// 		AttrIndex    string `json:"attr-index"`
// 		LastModified string `json:"last-modified"`
// 		Origin       string `json:"origin"`
// 		PathID       int    `json:"path-id"`
// 		Prefix       string `json:"prefix"`
// 		ValidRoute   bool   `json:"valid-route"`
// 	} `json:"state"`
// }

// type BgpRibRoutes struct {
// 	Routes []BgpRibRoute `json:"openconfig-network-instance:route"`
// }

func main() {
	// create a target
	tg, err := api.NewTarget(
		api.Name("srl1"),
		api.Address("192.168.121.101:6030"),
		api.Username("admin"),
		api.Password("admin"),
		api.SkipVerify(true),
		api.Insecure(true),
	)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create a gNMI client
	err = tg.CreateGNMIClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer tg.Close()

	// create a GetRequest
	getReq, err := api.NewGetRequest(
		api.Path("network-instances/network-instance/protocols/protocol/bgp/neighbors/neighbor"),
		api.Encoding("json_ietf"))
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(prototext.Format(getReq))

	// send the created gNMI GetRequest to the created target
	getResp, err := tg.Get(ctx, getReq)
	if err != nil {
		log.Fatal(err)
	}
	val := getResp.GetNotification()
	for _, n := range val {
		for _, v := range n.Update {
			fmt.Println()
			// fmt.Println(v)
			test1 := v.Val
			// fmt.Println(test1)
			test2 := test1.GetJsonIetfVal()
			// err := json.Unmarshal(test2, &bgpRib)
			// if err != nil {
			// 	log.Fatalf("failed to unmarshal jsonIetfVal: %v", err)
			// }

			fmt.Printf("%s\n", string(test2))
			// fmt.Println(prototext.Format(getResp))
		}
	}
	// var bgpRib BgpRibRoutes
	// val := getResp.GetNotification()
	// for _, n := range val {
	// 	for _, v := range n.Update {
	// 		fmt.Println()
	// 		fmt.Println(v)
	// 		test1 := v.Val
	// 		// fmt.Println(test1)
	// 		test2 := test1.GetJsonIetfVal()
	// 		err := json.Unmarshal(test2, &bgpRib)
	// 		if err != nil {
	// 			log.Fatalf("failed to unmarshal jsonIetfVal: %v", err)
	// 		}

	// 	}
	// }

	// fmt.Println(bgpRib.Routes[0].Prefix)
	// for _, route := range bgpRib.Routes {
	// 	fmt.Println("Prefix:", route.Prefix)
	// 	fmt.Println("Valid:", route.State.ValidRoute)
	// }
	// fmt.Println(prototext.Format(getResp.Data.Value))
}
