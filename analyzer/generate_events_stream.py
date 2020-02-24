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
	('W1', random_dist(3.500)),
	('W2', random_dist(2.372)),
	('W3', random_dist(1.607)),
	('W4', random_dist(1.089)),
	('W5', random_dist(0.738)),
	('W6', random_dist(0.500)),
	('W7', random_dist(0.339)),
	('W8', random_dist(0.230)),
	('W9', random_dist(0.155)),
	('T1', random_dist(0.105)),
	('T2', random_dist(0.071)),
	('T3', random_dist(0.048)),
	('I1', random_dist(0.0329)),
	('I2', random_dist(0.0223)),
	('I3', random_dist(0.0151)),
    ('D1', random_dist(0.0102)),
	('D2', random_dist(0.0069)),
	('D3', random_dist(0.0047)),
	('D4', random_dist(0.0031)),
	('D5', random_dist(0.0021)),
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
