# Modules
from pygnmi.client import gNMIclient, telemetryParser
import json

# Variables
host = ('192.168.121.101', '6030')

path = "network-instances/network-instance/afts/ipv4-unicast/ipv4-entry" #routing table
# path = "/network-instances/network-instance/name" # vrf list / read/write
# path = "/network-instances/network-instance/afts/state-synced" # vrf list
# path = "network-instances" # vrf list / read/write



# subscribe = {
#     "subscription": [
#         {
#             "path": "/network-instances/network-instance/afts/ipv4-unicast/ipv4-entry",
#             "mode": "Update",
#           #   "sample_interval": 10000000000,
#         },
#     ],
#     "mode": "stream",
#     "encoding": "json",
# }
# Body
if __name__ == '__main__':
    with gNMIclient(target=host, username='admin', password='admin', insecure=True) as gc:
        result = gc.get(path=[path])
            # telemetry_stream = gc.subscribe(subscribe=subscribe)
            # for telemetry_entry in telemetry_stream:
            #    print(telemetryParser(telemetry_entry))
    clean_path = path.replace("/", "_").replace("-", "_").replace(" ", "_")
    print(json.dumps(result, indent=4))
    with open("result.json", "w") as f:
        # Write the result to a file
        json.dump(result, open(f'{clean_path}.json', "w"), indent=4)
