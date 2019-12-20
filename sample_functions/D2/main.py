import json
import time
import random

def lambda_handler(event, context):
	ret = [0] * 6
	for i in range(25):
		for _ in range(1000 * 10000):
			ret[i % 6] += random.random()

		time.sleep(2)

	return {
		'statusCode': 200,
		'body': json.dumps({'data': ret})
    }
