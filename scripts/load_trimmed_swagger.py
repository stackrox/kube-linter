import gzip
import json
import sys
import urllib.request

def remove_description_keys(obj):
    if type(obj) is not dict:
        return obj
    new_obj = {}
    for k, v in obj.items():
        if k == "description":
            new_obj[k] = ""
            continue
        new_obj[k] = remove_description_keys(v)
    return new_obj

if len(sys.argv) < 3:
    print(f"Usage: {sys.argv[0]} <k8s_tag> <out_file>")
    sys.exit(1)

k8s_tag, out_file = sys.argv[1], sys.argv[2]

swagger_json_url = f"https://raw.githubusercontent.com/kubernetes/kubernetes/{k8s_tag}/api/openapi-spec/swagger.json"

with urllib.request.urlopen(swagger_json_url) as response:
    json_contents = response.read()

swagger_loaded = json.loads(json_contents)
swagger_trimmed = remove_description_keys(swagger_loaded)

with gzip.open(out_file, "w") as f:
    f.write(json.dumps(swagger_trimmed).encode())
