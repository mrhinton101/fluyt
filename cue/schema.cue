package schema

#IPCidr: =~"^((25[0-5]|2[0-4][0-9]|1?[0-9]{1,2})\\.){3}(25[0-5]|2[0-4][0-9]|1?[0-9]{1,2})/(3[0-2]|[12]?[0-9])$"

#Device: {
name: string
ip: #IPCidr
description? : string
telemetry: [... string]
config?: [... string]
pushmode?: "GNMI" | "Terraform"
}

