package fluyt

import "list"
import "strings"
import "net"



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
    
  _invalid_telemetry: [for x, y in telemetry 
    let tel_name = y
    if !list.Contains(#supported_paths, y) 
    {tel_name} ]

      if len(_invalid_telemetry) > 0 {   
    invalid_telemetry: _invalid_telemetry & error("unsupported telemetry paths for \(name): \(strings.Join(_invalid_telemetry, ", "))")
      
}
    
}
//create inventory and add device name to device definition
#inventory: {
    inventory: [k=string]: #device & { name: k }
}