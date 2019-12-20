import math
import json
import sys

if len(sys.argv) < 2:
	print('Usage: python analyzer.py <path/to/log>')
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

with open(sys.argv[1], 'r') as file:
	log = file.readlines()

NUM_NODES = 8

ret = [{} for _ in range(0, NUM_NODES)]
std_time = 999999999999

num_warm_starts = 0
num_cold_starts = 0

durations = []
latencies = []
durations_per_functions = {}
latencies_per_functions = {}
warm_per_functions = {}

for row in log:
	if len(row) <= 1:
		continue
	date, time, node_index, function_name, start_type, start_time, duration, latency = row.split(" ")
	
	start_time = int(start_time)
	duration = int(duration)
	latency = int(latency)
	node_index = int(node_index)
	
	s = math.floor(start_time / 1000)
	e = math.ceil((start_time + duration) / 1000)

	node_ret = ret[node_index]
	for i in range(s, e):
		if i not in node_ret:
			node_ret[i] = {}
		
		node_ret[i][function_name] = node_ret[i].get(function_name, 0) + 1

	std_time = min(std_time, s)

	if function_name not in durations_per_functions:
		durations_per_functions[function_name] = []
		latencies_per_functions[function_name] = []
		warm_per_functions[function_name] = 0

	if start_type == 'warm':
		num_warm_starts += 1
		warm_per_functions[function_name] += 1
	else:
		num_cold_starts += 1

	durations_per_functions[function_name].append(duration)
	latencies_per_functions[function_name].append(latency)
	durations.append(duration)
	latencies.append(latency)

# print('--------------- # of distinct functions ---------------')
avg_num_df = []
for node_index, data in enumerate(ret):
	arr = [len(d.keys()) for d in data.values()]
	# print('[Node %d] %s (avg. %.1f)' % (
	# 	node_index,
	# 	join2(arr),
	# 	avg(arr)
	# ))
	avg_num_df.append(avg(arr))

print('# of distinct functions for each node: %.1f (stddev.: %.1f)' % (
	avg(avg_num_df),
	stdd(avg_num_df),
))


# print('-------------------- loads --------------------')
avg_executions = []
avg_capacities = []
for node_index, data in enumerate(ret):
	arr = [sum(d.values()) for d in data.values()]
	# print('[Node %d] %s (avg. %.1f)' % (
	# 	node_index,
	# 	join2(arr),
	# 	avg(arr)
	# ))
	avg_executions.append(avg(arr))
	avg_capacities.append(8-avg(arr))

print('executions per node/sec: %.1f (stddev.: %.1f)' % (
	avg(avg_executions),
	stdd(avg_executions),
))
print('capacities per node/sec: %.1f (stddev.: %.1f)' % (
	avg(avg_capacities),
	stdd(avg_capacities),
))

# print('-------------------- warm/cold --------------------')
print('# of warm starts:', num_warm_starts)
print('# of cold starts:', num_cold_starts)
print('warm (%%): %d%%' % (num_warm_starts / (num_warm_starts + num_cold_starts) * 100))


# print('-------------------- exec time / latency --------------------')
# for key, arr in sorted(durations_per_functions.items(), key=lambda e: e[0]):
# 	print('[%s]: avg exec time: %dms, avg latency: %dms, warm: %d/%d' % (
# 		key,
# 		avg(arr),
# 		avg(latencies_per_functions[key]),
# 		warm_per_functions[key],
# 		len(latencies_per_functions[key]),
# 	))

print('avg exec time: %dms\navg latency: %dms' % (
	avg(durations),
	avg(latencies),
))