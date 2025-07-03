package fluyt

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

