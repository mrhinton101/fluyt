package fluyt

#telemetry_paths: {
    aft : #aft
}

#aft: {
    path: ["network-instances","network-instance","afts","ipv4-unicast","ipv4-entry"]
    description: "AFT IPv4 Table"
    rpc: {
        "get": true
        "set": false
        "subscribe": {
            "supported": true
            "mode": "on_change"
            "interval": null
    }
    }
    provider: "openconfig"
    tags: ["aft"]
}

