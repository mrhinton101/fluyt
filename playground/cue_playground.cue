import ("list"
"strings"
"net"
)

#telemetry_paths: {
    aft : #aft
}

#aft: {
    path: ["network-instances","network-instance","afts","ipv4-unicast","ipv4-entry"]
    description: "AFT IPv4 Table"
    rpc: {
        "get": true
        "set": false
        "subscribe": true
    }
    provider: "openconfig"
    tags: ["aft"]
}



#telemetry_paths: {
  [name=string]: #gnmi_path_meta
}

#gnmi_path_meta: {
    path: [...string]                 // gNMI path segments
    description: string              // what the path captures
    rpc: {  "get" : bool
            "set": bool
            "subscribe": bool}        // whether it's Get-only or supports Set
    provider: string                 // "openconfig", "vendor-x", etc.
    tags?: [...string]               // optional tags like "interface", "bgp", etc.
    type?: "string" | "int" | "bool" | "float" | "bytes"
    arista_eos?: string // eos version path confirmed
    cisco_nxos?: string // nxos version path confirmed
    junos?: string // junos version path confirmed
}

//ip address definition

//device entry schema
#device: {
    name: string
    ip:   net.IPCIDR
    telemetry: [...string]
    tags: {
    "region" : string
    "env": string 
    }
    description?: string
    config?: [...string]
    pushmode?: "GNMI" | "Terraform" | "Pulumi" 
    tel_paths: [for x, y in #telemetry_paths 
    let tel_name = x
    
    if list.Contains(telemetry, x) 
    let tel_path = strings.Join(y.path, "/")
    { "\(tel_name)" : tel_path} ]
}


#inventory: {
    inventory: [k=string]: #device & { name: k }
}

#inventory & {
inventory:{
  device1:{
    ip: "10.1.1.1/22"
    tags: {
      "region": "amer"
      "env": "lab"}
    telemetry:[ "aft", "fu"]
      }
  device2: {
    ip: "10.1.1.2/22"
    tags:{
      "region": "amer"
      "env": "lab"}
    telemetry:["fo"]
      }
}
}