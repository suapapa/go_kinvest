
#!/usr/bin/env python

import os
import yaml

# open kinvest_prod.yaml
with open(os.path.join(os.path.dirname(__file__), 'kinvest_prod.yaml'), 'r', encoding='utf-8') as f:
    data = yaml.safe_load(f)
    # print (data).path.*.summary
    paths = data['paths']
    for path, methods in paths.items():
        for method, details in methods.items():
            if 'summary' in details:
                summary = details['summary'].replace('J_', '')
                print(f"- [ ] {path} ({method}) : {summary}")
            else:
                print(f"- [ ] {path} ({method})")
