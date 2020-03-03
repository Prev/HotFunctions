function sleep(ms){
    return new Promise(resolve=>{
        setTimeout(resolve,ms)
    })
}

exports.handle = async (event, context) => {
	let sum = [0, 0, 0, 0, 0, 0, 0, 0, 0, 0];

	for (let k = 0; k < 50; k++) {
		for (let i = 0; i < 2000 * 500; i++) {
			sum[k%7] += Math.random();
		}
		await sleep(2000);
	}

	return {
		'statusCode': 200,
		'body': JSON.stringify(sum),
	};
};
