import random
import math
import json
import sys

if len(sys.argv) < 2:
	print('Usage: python generate_event_stream.py <path/to/output.csv>')
	sys.exit(-1)

EXPERIMENTAL_TIME = 60 * 5 # 5min

def random_dist(call_per_sec):
	# call per minutes
	total = round(call_per_sec * EXPERIMENTAL_TIME)

	ret = []
	for _ in range(total):
		t = math.floor(random.random() * EXPERIMENTAL_TIME * 1000) # msec
		ret.append(t)
	return ret

functions = [
	('W1', random_dist(1.6355)),
	('W2', random_dist(1.2737)),
	('W3', random_dist(0.9920)),
	('W4', random_dist(0.7725)),
	('W5', random_dist(0.6017)),
	('W6', random_dist(0.4686)),
	('W7', random_dist(0.3649)),
	('W8', random_dist(0.2842)),
	('W9', random_dist(0.2213)),
	('W10', random_dist(0.1724)),
	('W11', random_dist(0.1343)),
	('T1', random_dist(0.1046)),
	('T2', random_dist(0.0814)),
	('T3', random_dist(0.0634)),
	('T4', random_dist(0.0494)),
	('T5', random_dist(0.0384)),
	('T6', random_dist(0.0300)),
	('I1', random_dist(0.0233)),
	('I2', random_dist(0.0182)),
	('I3', random_dist(0.0141)),
	('I4', random_dist(0.0110)),
	('I5', random_dist(0.0086)),
	('D1', random_dist(0.0066)),
	('D2', random_dist(0.0052)),
	('D3', random_dist(0.0041)),
	('D4', random_dist(0.0031)),
	('D5', random_dist(0.0024)),
	('D6', random_dist(0.0019)),
	('D7', random_dist(0.0015)),
	('D8', random_dist(0.0012)),
]

event_stream = []
for f in functions:
	name, scheds = f
	for msec in scheds:
		event_stream.append((name, msec))

event_stream.sort(key=lambda e: e[1]) # sort by start time

summary = {}

with open(sys.argv[1], 'w') as file:
	for name, start_time in event_stream:
		file.write("%s,%s\n" % (name, start_time))

		summary[name] = summary.get(name, 0) + 1

for name, _ in functions:
	cnt = summary.get(name, 0)
	print('%s: %d' % (name, cnt))
