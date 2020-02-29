import json
import random

def lambda_handler(event, context):
	N = 8000000
	ret = 0
	for _ in range(N):
		ret += random.random()

	return {
		'statusCode': 200,
		'body': json.dumps({'ret': ret})
    }
