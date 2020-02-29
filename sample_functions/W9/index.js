exports.handle = (event, context) => {
	let sum = 0;
	for (let i = 0; i < 13000000; i++) {
		sum += Math.random();
	}

	return {
		'statusCode': 200,
		'body': 'Hello ' + sum
	};
};
