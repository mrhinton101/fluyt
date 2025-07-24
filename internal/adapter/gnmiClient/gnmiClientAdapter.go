package gnmiClient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"

	"github.com/mrhinton101/fluyt/domain/cue"
	"github.com/mrhinton101/fluyt/domain/gnmi"
	"github.com/mrhinton101/fluyt/internal/app/core/logger"
	"github.com/mrhinton101/fluyt/internal/app/ports"
	pb "github.com/openconfig/gnmi/proto/gnmi"
	"github.com/openconfig/gnmic/pkg/api"
	"github.com/openconfig/gnmic/pkg/api/target"
)

type GNMIClientImpl struct {
	Name    string
	Address string
	Port    string
	tg      *target.Target
}

type GNMICapabilityResponse struct {
	Target   string
	Response pb.CapabilityResponse
}

func NewGNMIClient(device cue.Device) ports.GNMIClient {
	return &GNMIClientImpl{
		Name:    device.Name,
		Address: fmt.Sprintf("%s:6030", device.Address),
		Port:    device.Port,
	}
}

func (c *GNMIClientImpl) Init(ctx context.Context) error {
	// Run the actual gNMI Capabilities RPC and get result
	tg, err := api.NewTarget(
		api.Name(c.Name),
		api.Address(c.Address),
		api.Username("admin"),
		api.Password("admin"),
		api.SkipVerify(true),
		api.Insecure(true))

	fmt.Println(err)

	// create a gNMI client
	err = tg.CreateGNMIClient(ctx)
	if err != nil {
		logger.SLogger(logger.LogEntry{
			Level:     slog.LevelError,
			Component: "gnmiClient",
			Msg:       fmt.Sprintf("failed to create gNMI client for target"),
			Err:       err,
			Target:    c.Name,
		})
		return err
	}
	c.tg = tg
	return nil
}

func (c *GNMIClientImpl) Capabilities(ctx context.Context) (map[string]interface{}, error) {
	// Run the actual gNMI Capabilities RPC and get result
	fmt.Println("inside capabilities")
	if c.tg == nil {
		logger.SLogger(logger.LogEntry{
			Level:     slog.LevelError,
			Err:       errors.New("GNMI Target is null, Did you call Init() first?"),
			Component: "gnmiClient",
			Action:    "Confirm Target",
			Msg:       "gNMI client not initialized",
			Target:    c.Name,
		})
	}
	capResp, err := c.tg.Capabilities(ctx)
	fmt.Println("capresp:")
	if err != nil {
		logger.SLogger(logger.LogEntry{
			Level:     slog.LevelError,
			Component: "gnmiClient",
			Msg:       "failed to get capabilities from target",
			Err:       err,
			Target:    c.Name,
		})
		log.Fatal(err)
	}

	resp, err := ValidateCapabilityResponse(c.Address, *capResp)
	if err != nil {
		return nil, err
	}

	// convert to map[string]interface{}
	result, err := UnmarshalCapabilityResponse(resp)
	if err != nil {
		return nil, err
	}

	// for now flatten all capability fields into a map
	return map[string]interface{}{
		"target":    result.Target,
		"encodings": result.Encodings,
		"models":    result.Models,
		"versions":  result.Versions,
	}, nil
}

func (c *GNMIClientImpl) GetBgpRibs(ctx context.Context) (gnmi.BgpRibs, error) {

	// create a GetRequest
	getReq, err := api.NewGetRequest(
		api.Path("network-instances/network-instance/protocols/protocol/bgp/rib/afi-safis/afi-safi/ipv4-unicast/loc-rib/routes"),
		api.Encoding("json_ietf"))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	// fmt.Println(prototext.Format(getReq))

	// send the created gNMI GetRequest to the created target
	getResp, err := c.tg.Get(ctx, getReq)
	if err != nil {
		log.Fatal(err)
	}
	val := getResp.GetNotification()
	BgpRib := ParseBgpRibResp(val)
	if len(BgpRib) == 0 {
		logger.SLogger(logger.LogEntry{
			Level:     slog.LevelError,
			Err:       fmt.Errorf("nil BGP RIB received"),
			Component: "gNMI",
			Action:    "GetBgpRib",
			Msg:       "nil BGP RIB received. Is provider bgp enabled on the device?",
			Target:    c.Address,
		})
	}

	return BgpRib, nil
}

func (c *GNMIClientImpl) GetAddress() string {
	return c.Address
}

func (c *GNMIClientImpl) Close() {
	c.tg.Close()
}

func ValidateCapabilityResponse(target string, capResp pb.CapabilityResponse) (*GNMICapabilityResponse, error) {
	if len(capResp.GNMIVersion) == 0 && len(capResp.SupportedModels) == 0 && len(capResp.SupportedEncodings) == 0 {
		return nil, fmt.Errorf("no capabilities received for target %s", target)
	}
	return &GNMICapabilityResponse{
		Target:   target,
		Response: capResp,
	}, nil
}

func UnmarshalCapabilityResponse(capResp *GNMICapabilityResponse) (gnmi.CleanCapabilityResponse, error) {
	if capResp == nil {
		fmt.Errorf("capability response is nil")
	}

	result := capResp.Response
	models := make([]string, len(result.SupportedModels))
	for i, m := range result.SupportedModels {
		models[i] = m.Name
	}

	encodings := make([]string, len(result.SupportedEncodings))
	for i, e := range result.SupportedEncodings {
		encodings[i] = e.String()
	}
	// for now flatten all capability fields into a map
	return gnmi.CleanCapabilityResponse{
		Target:    capResp.Target,
		Encodings: encodings,
		Models:    models,
		Versions:  result.GNMIVersion,
	}, nil
}

// func bgpRibRespMapper(gnmiRoutesResp []gnmi.BgpRibRoute) (ribRoutes gnmi.BgpRibs) {
// 	for _, route := range gnmiRoutesResp {
// 		prefix := route.Prefix
// 		ribRoutes[prefix] = route
// 	}
// 	return ribRoutes
// }

func extractVRF(path *pb.Path) gnmi.BgpVrfName {
	for _, elem := range path.Elem {
		if elem.Name == "network-instance" {
			return gnmi.BgpVrfName{Name: elem.Key["name"]}
		}
	}

	logger.SLogger(logger.LogEntry{
		Level:     slog.LevelError,
		Err:       errors.New("failed to extract VRF from path"),
		Component: "gnmiClient",
		Action:    "extractVRF",
		Msg:       "VRF name not found in path",
		Target:    path.Target,
	})
	return gnmi.BgpVrfName{Name: ""}
}

func ParseBgpRibResp(notifs []*pb.Notification) gnmi.BgpRibs {
	ribs := make(map[gnmi.BgpVrfName]map[gnmi.BgpRibKey]gnmi.BgpRibRoute)

	for _, n := range notifs {
		for _, u := range n.Update {
			vrf := extractVRF(u.Path)

			var gnmiRib gnmi.GnmiBgpRibRoutes
			err := json.Unmarshal(u.Val.GetJsonIetfVal(), &gnmiRib)
			if err != nil {
				logger.SLogger(logger.LogEntry{
					Level:     slog.LevelError,
					Err:       err,
					Component: "gnmiClient",
					Action:    "ParseBgpRibResp",
					Msg:       "failed to unmarshal BGP RIB routes",
					Target:    vrf.Name})
				continue
			}
			// fmt.Println("routes:", gnmiRib)
			if ribs[vrf] == nil {
				ribs[vrf] = make(map[gnmi.BgpRibKey]gnmi.BgpRibRoute)
			}

			for _, r := range gnmiRib.Routes {
				fmt.Println("route:", r)
				key := gnmi.BgpRibKey{Prefix: r.Prefix}
				ribs[vrf][key] = r
			}
		}
	}

	return ribs
}
