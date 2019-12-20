import json
import time

def lambda_handler(event, context):
	cnt = 3
	for _ in range(cnt):
		time.sleep(0.8)

	return {
		'statusCode': 200,
		'body': json.dumps({'sleep': '0.8 * 3'})
    }
