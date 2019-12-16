const m = require('./index.js');

const startTime = (new Date()).getTime()
const result = m.handle();
const endTime = (new Date()).getTime()

const out = {
	startTime: startTime,
	result: result,
	endTime: endTime,
}

const dummpedData = JSON.stringify(out)
console.log('-----' + dummpedData.length + '-----=' + dummpedData)
