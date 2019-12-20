import random
import math
import json
import sys

if len(sys.argv) < 2:
	# print('Usage: python generate_event_stream.py <path/to/output.json>')
	print('Usage: python generate_event_stream.py <path/to/output.csv>')
	sys.exit(-1)

WEBSERVER = 'webserver'
DATA_PROCESSING = 'data_processing'
THIRD_PARTY = '3rd-party'
INTERNAL_TOOLING = 'internal_tooling'

def random_dist(call_per_sec):
	# call per minutes
	total = round(call_per_sec * 60)

	ret = []
	for _ in range(total):
		t = math.floor(random.random() * 60 * 1000) # msec
		ret.append(t)
	return ret

functions = [
	('W1', WEBSERVER, random_dist(0.1)),
	('W2', WEBSERVER, random_dist(0.8)),
	('W3', WEBSERVER, random_dist(0.9)),
	('W4', WEBSERVER, random_dist(2)),
	('W5', WEBSERVER, random_dist(3)),
	('W6', WEBSERVER, random_dist(1.5)),
	('W7', WEBSERVER, random_dist(1)),
	('W8', WEBSERVER, random_dist(0.2)),
	('W9', WEBSERVER, random_dist(0.0)),
	('D1', DATA_PROCESSING, [5000, 30000]),
	('D2', DATA_PROCESSING, [6000, 6000, 6000]),
	('D3', DATA_PROCESSING, [3000, 33000]),
	('D4', DATA_PROCESSING, [5000, 20000, 30000, 40000]),
	('D5', DATA_PROCESSING, [6000, 23000, 33000, 43000]),
	('T1', THIRD_PARTY, random_dist(0.1)),
	('T2', THIRD_PARTY, random_dist(0.2)),
	('T3', THIRD_PARTY, random_dist(0.3)),
	('I1', INTERNAL_TOOLING, [0, 10000, 20000, 30000, 40000, 50000]),
	('I2', INTERNAL_TOOLING, [2000, 12000, 22000, 32000, 42000, 52000]),
	('I3', INTERNAL_TOOLING, [4000, 14000, 24000, 34000, 44000, 54000]),
]

event_stream = []
for f in functions:
	name, ftype, scheds = f
	for msec in scheds:
		event_stream.append((name, msec))

event_stream.sort(key=lambda e: e[1]) # sort by start time


with open(sys.argv[1], 'w') as file:
	# file.write(json.dumps(event_stream))
	for event in event_stream:
		file.write("%s,%s\n" % (event[0], event[1]))
