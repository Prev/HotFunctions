function sleep(ms){
    return new Promise(resolve=>{
        setTimeout(resolve,ms)
    })
}

exports.handle = async (event, context) => {
	let sum = [0, 0, 0, 0, 0, 0, 0, 0, 0, 0];

	for (let k = 0; k < 60; k++) {
		for (let i = 0; i < 2000 * 1000; i++) {
			sum[k%8] += Math.random();
		}
		await sleep(1500);
	}

	return {
		'statusCode': 200,
		'body': JSON.stringify(sum),
	};
};
