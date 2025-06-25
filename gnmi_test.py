# Modules
from pygnmi.client import gNMIclient
import json

# Variables
host = ('192.168.121.101', '6030')

path = "/network-instances/network-instance/afts/ipv4-unicast/ipv4-entry"

# Body
if __name__ == '__main__':
    with gNMIclient(target=host, username='admin', password='admin', insecure=True) as gc:
         result = gc.get(path=[path])
    print(json.dumps(result, indent=4))