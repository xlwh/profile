#!/usr/bin/env python
# coding=utf-8

import json


f = open('./profile.dump')
content = f.read()
js = json.loads(content)
for k in js.keys():
    print k + ":" + str(js[k])
f.close()