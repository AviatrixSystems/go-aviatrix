import requests
import sys
import os
import time
import re
import json
from requests.packages.urllib3.exceptions import InsecureRequestWarning
requests.packages.urllib3.disable_warnings(InsecureRequestWarning)

# qresponse = requests.get("https://13.126.166.7/v1/api?action=login&username=admin&password=Aviatrix123#", verify=False)
#response_json = response.json()
mydata={"CID":"59951c322abf9", "action":"connect_container","cloud_type":1,"account_name":"devops","gw_name":"avtxgw2","vpc_id":"vpc-2ee4a147","vpc_reg":"ap-south-1","vpc_size":"t2.micro","vpc_net":"10.2.0.0/24"}
print type(json.dumps(mydata))
response = requests.post("https://13.126.166.7/v1/api", mydata, verify=False)
response_json = response.json()
print response_json

