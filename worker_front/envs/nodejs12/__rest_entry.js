var http = require('http');
var childProcess = require('child_process');

var server = http.createServer();

server.addListener('request', function (request, response) {  
	childProcess.exec('node __entry.js', function (error, stdout, stderr) {
		response.writeHead(200, {'Content-Type' : 'text/plain'});
	    response.write(stdout);
	    response.end();
	});
});

console.log('Server listening at :8080');
server.listen(8080, '0.0.0.0');
