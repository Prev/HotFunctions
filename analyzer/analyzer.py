import math
import json
import sys
import re
from pprint import pprint

if len(sys.argv) < 2:
	print('Usage: python analyzer.py <path/to/log> [<detail_mode=0>]')
	sys.exit(-1)

def avg(datalist):
	if len(datalist) == 0:
		return 0
	return sum(datalist) / len(datalist)

def join2(datalist, joinstr=' '):
	return joinstr.join([str(d) for d in datalist])

def stdd(datalist):
	avg_n = avg(datalist)
	v = avg([(d - avg_n) ** 2 for d in datalist])
	return math.sqrt(v)

def percentile(num, total):
	return '%d%%' % (num / total * 100)

with open(sys.argv[1], 'r') as file:
	log = file.readlines()

DETAIL_MODE = len(sys.argv) == 3 and sys.argv[2] != '0'

NUM_NODES = 10
LOG_FORMAT = re.compile('(\S+)\s(\S+)\s(\d+)\s(\d+)\s(\S+)\s(.+)')
std_time = math.inf
last_time = -1

records = []

# Parse log
for row in log[20:-20]:
	if len(row) <= 1:
		continue

	matches = LOG_FORMAT.search(row)
	start_time = int(matches.group(3))
	end_time = int(matches.group(4))
	function_name = matches.group(5)
	data = json.loads(matches.group(6))

	if 'Error' in data:
		continue

	std_time = min(std_time, start_time)
	last_time = max(last_time, end_time)

	records.append((
		start_time,
		end_time,
		function_name,
		data,
	))


# ret = [{} for _ in range(0, NUM_NODES)]
max_timeslot = int((last_time - std_time) / 1000) + 1
ret = []
for _ in range(0, NUM_NODES):
	arr = []
	for _ in range(0, max_timeslot):
		arr.append({})
	ret.append(arr)


# Init data structures
durations = []
latencies = []
durations_per_functions = {}
latencies_per_functions = {}

num_total = 0

num_image_hit = 0
opt1_execution_times = ([], [])
opt1_latencies = ([], [])

num_using_pooled_container = 0
opt2_execution_times = ([], [])
opt2_latencies = ([], [])

num_using_existing_rest_container = 0
opt3_execution_times = ([], [])
opt3_latencies = ([], [])

executions_per_node = [{} for _ in range(0, NUM_NODES)]

for start_time, end_time, function_name, data in records:
	s = math.floor((start_time - std_time) / 1000)
	e = math.ceil((end_time - std_time) / 1000)

	node_id = data['LoadBalancingInfo']['WorkerNodeId']
	
	for i in range(s, e):
		ret[node_id][i][function_name] = ret[node_id][i].get(function_name, 0) + 1

	duration = end_time - start_time
	latency = duration - data['InternalExecutionTime']

	if function_name not in durations_per_functions:
		durations_per_functions[function_name] = []
		latencies_per_functions[function_name] = []

	durations_per_functions[function_name].append(duration)
	latencies_per_functions[function_name].append(latency)
	durations.append(duration)
	latencies.append(latency)

	num_total += 1

	if data['Meta']['ImageBuilt']:
		opt1_execution_times[1].append(duration)
		opt1_latencies[1].append(latency)
	else:
		num_image_hit += 1
		opt1_execution_times[0].append(duration)
		opt1_latencies[0].append(latency)

	if data['Meta']['UsingPooledContainer']:
		num_using_pooled_container += 1
		opt2_execution_times[0].append(duration)
		opt2_latencies[0].append(latency)
	else:
		opt2_execution_times[1].append(duration)
		opt2_latencies[1].append(latency)

	if data['Meta']['UsingExistingRestContainer']:
		num_using_existing_rest_container += 1
		opt3_execution_times[0].append(duration)
		opt3_latencies[0].append(latency)
	else:
		opt3_execution_times[1].append(duration)
		opt3_latencies[1].append(latency)


	executions_per_node[node_id][function_name] = executions_per_node[node_id].get(function_name, 0) + 1

if DETAIL_MODE:
	print('Total:', num_total)
	for node_id, d in enumerate(executions_per_node):
		print('Node %d | total: %d | avg: %.1f/s | %s' % (
			node_id,
			sum(d.values()),
			sum(d.values()) / max_timeslot,
			d
		))

print('------------------------Locality------------------------')
avg_num_df = []
per_time = [{} for _ in range(0, max_timeslot)]

for node_id, data in enumerate(ret):
	# Sample format of `data`:
	#   [{'W7': 1, 'T2': 1, 'W4': 2, 'W5': 2, 'D3': 1}, ...]

	arr = [len(d.keys()) for d in data]
	# `arr` means number of distinct functions of the node where index is timeslot
	# Sample format of `arr`:
	#	[1, 2, 1, 2, 2, 3, 1, ...]

	avg_num_df.append(avg(arr))

	for timeslot, dic in enumerate(data):
		per_time[timeslot][node_id] = len(dic.keys())

print('# of distinct functions for each node: %.1f (stddev.: %.1f)' % (
	avg(avg_num_df),
	stdd(avg_num_df),
))

if DETAIL_MODE:
	for timeslot, row in enumerate(per_time):
		if timeslot > 300: continue
		print('(%d, %.1f)' % (timeslot, avg(row.values())), end=' ')

	print('\n')

print('\n------------------------Imbalance------------------------')
avg_executions = []
per_time = [{} for _ in range(0, max_timeslot)]

for node_id, data in enumerate(ret):
	arr = [sum(d.values()) for d in data]
	avg_executions.append(avg(arr))

	for timeslot, dic in enumerate(data):
		per_time[timeslot][node_id] = sum(dic.values())

print('CV: %.2f (sttdev.: %.2f, avg: %.2f)' % (
	stdd(avg_executions) / avg(avg_executions),
	stdd(avg_executions),
	avg(avg_executions),
))

if DETAIL_MODE:
	# print('avg load:')
	# for timeslot, row in enumerate(per_time):
	# 	if timeslot > 300 or timeslot < 10: continue
	# 	v = row.values()
	# 	print('(%d, %.2f)' % (timeslot,  avg(v)), end=' ')
	# print('')

	print('CV:')
	for timeslot, row in enumerate(per_time):
		if timeslot > 300 or timeslot < 10: continue
		v = row.values()
		cv = stdd(v) / avg(v) if avg(v) != 0 else 0

		print('(%d, %.2f)' % (timeslot, cv), end=' ')
	print('')

print('\n-------------------- Cache hits --------------------')
print('                           | Hit | Miss | %  | H.E.Time | M.E.Time')
print('ImageReuse                 | %d | %d | %s | %dms | %dms |' % (
	num_image_hit,
	num_total - num_image_hit,
	percentile(num_image_hit, num_total),
	avg(opt1_execution_times[0]),
	avg(opt1_execution_times[1]),
))
print('UsingPooledContainer       | %d | %d | %s | %dms | %dms |' % (
	num_using_pooled_container,
	num_total - num_using_pooled_container,
	percentile(num_using_pooled_container, num_total),
	avg(opt2_execution_times[0]),
	avg(opt2_execution_times[1]),
))
print('UsingExistingRestContainer | %d | %d | %s | %dms | %dms |' % (
	num_using_existing_rest_container,
	num_total - num_using_existing_rest_container,
	percentile(num_using_existing_rest_container, num_total),
	avg(opt3_execution_times[0]),
	avg(opt3_execution_times[1]),
))
print('')


print('-------------------- Exec time / latency --------------------')

print('avg exec time: %dms\navg latency: %dms' % (
	avg(durations),
	avg(latencies),
))
print('')


print('Total executions:', num_total)
functions = sorted(durations_per_functions.keys())
for fname in functions:
	cnt = len(durations_per_functions[fname])
	dur = avg(durations_per_functions[fname])
	latency = avg(latencies_per_functions[fname])

	print('%s: %dms (n=%d)(internal: %dms, latency: %dms)' % (
		fname,
		dur,
		cnt,
		dur - latency,
		latency,
	))
