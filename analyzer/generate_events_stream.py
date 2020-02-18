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

EXPERIMENTAL_TIME = 60 * 5 # 10min

def random_dist(call_per_sec):
	# call per minutes
	total = round(call_per_sec * EXPERIMENTAL_TIME)

	ret = []
	for _ in range(total):
		t = math.floor(random.random() * EXPERIMENTAL_TIME * 1000) # msec
		ret.append(t)
	return ret

functions = [
	('W1', WEBSERVER, random_dist(1.1)),
	('W2', WEBSERVER, random_dist(0.8)),
	('W3', WEBSERVER, random_dist(0.9)),
	('W4', WEBSERVER, random_dist(1)),
	('W5', WEBSERVER, random_dist(3.5)),
	('W6', WEBSERVER, random_dist(1.5)),
	('W7', WEBSERVER, random_dist(1)),
	('W8', WEBSERVER, random_dist(0.2)),
	('W9', WEBSERVER, random_dist(0.1)),
    ('D1', DATA_PROCESSING, random_dist(0.033)),
	('D2', DATA_PROCESSING, random_dist(0.05)),
	('D3', DATA_PROCESSING, random_dist(0.033)),
	('D4', DATA_PROCESSING, random_dist(0.066)),
	('D5', DATA_PROCESSING, random_dist(0.066)),
	('T1', THIRD_PARTY, random_dist(0.1)),
	('T2', THIRD_PARTY, random_dist(0.2)),
	('T3', THIRD_PARTY, random_dist(0.3)),
	('I1', INTERNAL_TOOLING, random_dist(0.01)),
	('I2', INTERNAL_TOOLING, random_dist(0.01)),
	('I3', INTERNAL_TOOLING, random_dist(0.01)),
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
