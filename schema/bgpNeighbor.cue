package fluyt

import "net"

#bgp_valid_asn: <=4294967295

#neighbor: {
  peer_group?: "LEAF" | "SPINE" 
  description?: string
  route_map_in?: string
  remote_as: int & #bgp_valid_asn
}

router_bgp: {
  neighbors: {
    [net.IP]: #neighbor
  }
}