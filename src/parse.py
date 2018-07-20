#!/usr/bin/env python
# coding=utf-8

import json

{
	"Alloc": 328696,
	"TotalAlloc": 10804936,
	"Sys": 7903480,
	"Lookups": 3421,
	"Mallocs": 17749,
	"Frees": 16919,
	"HeapAlloc": 328696,
	"HeapSys": 5079040,
	"HeapIdle": 4374528,
	"HeapInuse": 704512,
	"HeapReleased": 0,
	"HeapObjects": 830,
	"StackInuse": 163840,
	"StackSys": 163840,
	"MSpanInuse": 12312,
	"MSpanSys": 81920,
	"MCacheInuse": 6816,
	"MCacheSys": 16384,
	"BuckHashSys": 1442755,
	"GCSys": 372736,
	"OtherSys": 746805,
	"NextGC": 4194304,
	"LastGC": 1532050229662446300,
	"PauseTotalNs": 502200,
	"NumGC": 3,
	"EnableGC": false,
	"DebugGC": false,
	"InUseBytes": 0,
	"InUseObjects": 0,
	"NumGoroutine": 0
}

f = open('./profile.dump')
content = f.read()
js = json.loads(content)
print "Alloc:" + str(js['Alloc'])

f.close()
