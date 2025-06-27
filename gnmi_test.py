# Modules
from pygnmi.client import gNMIclient, telemetryParser
import json

# Variables
host = ('192.168.121.101', '6030')

path = "/network-instances/network-instance/afts/ipv4-unicast/ipv4-entry"

subscribe = {
    "subscription": [
        {
            "path": "/network-instances/network-instance/afts/ipv4-unicast/ipv4-entry",
            "mode": "Update",
          #   "sample_interval": 10000000000,
        },
    ],
    "mode": "stream",
    "encoding": "json",
}
# Body
if __name__ == '__main__':
    with gNMIclient(target=host, username='admin', password='admin', insecure=True) as gc:
     #     result = gc.get(path=[path])
          telemetry_stream = gc.subscribe(subscribe=subscribe)
          for telemetry_entry in telemetry_stream:
               print(telemetryParser(telemetry_entry))



    print(json.dumps(result, indent=4))