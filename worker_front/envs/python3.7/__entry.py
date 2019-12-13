import os
from main import lambda_handler
import time
import json

if __name__ == '__main__':
	event = os.environ.get('event', {})
	context = os.environ.get('context', {})

	start_time = time.time()
	out = lambda_handler(event, context)
	end_time = time.time()

	out['startTime'] = int(start_time * 1000)
	out['endTime'] = int(end_time * 1000)

	dummped_data = json.dumps(out)
	print('-----%d-----=%s' % (len(dummped_data), dummped_data))
