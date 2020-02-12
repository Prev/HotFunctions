const m = require('./index.js');

const startTime = (new Date()).getTime();
let result = m.handle();

function printResult(result) {
	const endTime = (new Date()).getTime();
	const dumpedData = JSON.stringify({
		startTime: startTime,
		result: result,
		endTime: endTime,
	});
	console.log('-=-=-=-=-=' + dumpedData.length + '-=-=-=-=-=>' + dumpedData + '==--==--==--==--==');
}

if (result instanceof Promise) {
	result.then(realResult => {
		printResult(realResult);
	});

} else {
	printResult(result);
}
