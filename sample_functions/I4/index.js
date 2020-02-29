function sleep(ms){
    return new Promise(resolve=>{
        setTimeout(resolve,ms)
    })
}

exports.handle = async (event, context) => {
	let sum = [0, 0, 0, 0, 0, 0, 0, 0, 0, 0];

	for (let k = 0; k < 10; k++) {
		for (let i = 0; i < 200000; i++) {
			sum[0] += Math.random();
		}
		await sleep(3000);
	}

	return {
		'statusCode': 200,
		'body': JSON.stringify(sum),
	};
};
