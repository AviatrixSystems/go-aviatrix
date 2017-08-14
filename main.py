import requests
import sys
import os
import time
import re
from requests.packages.urllib3.exceptions import InsecureRequestWarning
requests.packages.urllib3.disable_warnings(InsecureRequestWarning)

response = requests.get("https://13.126.166.7/v1/api?action=login&username=admin&password=Aviatrix123#", verify=False)
response_json = response.json()
print response_json