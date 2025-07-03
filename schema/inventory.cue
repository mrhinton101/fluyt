package fluyt
#ipcidr: =~"^((25[0-5]|2[0-4][0-9]|1?[0-9]{1,2})\\.){3}(25[0-5]|2[0-4][0-9]|1?[0-9]{1,2})/(3[0-2]|[12]?[0-9])$"

//device entry schema
#device: {
    name: string
    ip:   #ipcidr
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
#inventory: {
    inventory: [k=string]: #device & { name: k }
}