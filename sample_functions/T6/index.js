function sleep(ms){
    return new Promise(resolve=>{
        setTimeout(resolve,ms)
    })
}

exports.handle = async (event, context) => {
	await sleep(1000);
	await sleep(800);
	await sleep(500);

	return {
		'statusCode': 200,
		'body': '3 times for sleep'
	};
};
