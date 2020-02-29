import json
import time
import random

def lambda_handler(event, context):
	ret = [0] * 10
	for i in range(10):
		for _ in range(800000):
			ret[i] += random.random()

		time.sleep(3)

	return {
		'statusCode': 200,
		'body': json.dumps({'data': ret})
    }
