import json
import time

def lambda_handler(event, context):
	cnt = 4
	for _ in range(4):
		time.sleep(0.5)

	return {
		'statusCode': 200,
		'body': json.dumps({'sleep': '0.5 * 4'})
    }
