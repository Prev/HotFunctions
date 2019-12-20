const m = require('./index.js');

const startTime = (new Date()).getTime();
let result = m.handle();

if (result instanceof Promise) {
	result.then(realResult => {
		const endTime = (new Date()).getTime();
		const dummpedData = JSON.stringify({
			startTime: startTime,
			result: realResult,
			endTime: endTime,
		})
		console.log('-----' + dummpedData.length + '-----=' + dummpedData);
	});

} else {
	const endTime = (new Date()).getTime();
	const dummpedData = JSON.stringify({
		startTime: startTime,
		result: result,
		endTime: endTime,
	})
	console.log('-----' + dummpedData.length + '-----=' + dummpedData);
}
