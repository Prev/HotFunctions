import random
import math
import json
import sys

if len(sys.argv) < 2:
	# print('Usage: python generate_event_stream.py <path/to/output.json>')
	print('Usage: python generate_event_stream.py <path/to/output.csv>')
	sys.exit(-1)

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
	('W1', random_dist(2.3364)),
	('W2', random_dist(1.8196)),
	('W3', random_dist(1.4171)),
	('W4', random_dist(1.1036)),
	('W5', random_dist(0.8595)),
	('W6', random_dist(0.6694)),
	('W7', random_dist(0.5213)),
	('W8', random_dist(0.4060)),
	('W9', random_dist(0.3162)),
	('W10', random_dist(0.2463)),
	('W11', random_dist(0.1918)),
	('T1', random_dist(0.1494)),
	('T2', random_dist(0.1163)),
	('T3', random_dist(0.0906)),
	('T4', random_dist(0.0706)),
	('T5', random_dist(0.0549)),
	('T6', random_dist(0.0428)),
	('I1', random_dist(0.0333)),
	('I2', random_dist(0.0260)),
	('I3', random_dist(0.0202)),
	('I4', random_dist(0.0157)),
	('I5', random_dist(0.0123)),
	('D1', random_dist(0.0095)),
	('D2', random_dist(0.0074)),
	('D3', random_dist(0.0058)),
	('D4', random_dist(0.0045)),
	('D5', random_dist(0.0035)),
	('D6', random_dist(0.0027)),
	('D7', random_dist(0.0021)),
	('D8', random_dist(0.0017)),
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
