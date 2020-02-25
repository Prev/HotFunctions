import json
import time
import random

def lambda_handler(event, context):
	ret = [0] * 5
	for i in range(30):
		for _ in range(1000 * 5000):
			ret[i % 5] += random.random()

		time.sleep(2)

	return {
		'statusCode': 200,
		'body': json.dumps({'data': ret})
    }
