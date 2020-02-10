import os
from main import lambda_handler
import time
import json

if __name__ == '__main__':
	event = os.environ.get('event', {})
	context = os.environ.get('context', {})
	
	start_time = time.time()
	result = lambda_handler(event, context)
	end_time = time.time()

	out = {
		'startTime': int(start_time * 1000),
		'endTime': int(end_time * 1000),
		'result': result,
	}

	dump_data = json.dumps(out)
	print('-=-=-=-=-=%d-=-=-=-=-=>%s==--==--==--==--==' % (len(dump_data), dump_data))
