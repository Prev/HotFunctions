from http.server import BaseHTTPRequestHandler, HTTPServer
import subprocess

class ProxyHTTPServer_RequestHandler(BaseHTTPRequestHandler):
	def do_GET(self):
		self.send_response(200)

		self.send_header('Content-type','text/plain')
		self.end_headers()
		
		out = subprocess.check_output(["python", "__entry.py"])
		self.wfile.write(out)
		return

if __name__ == '__main__':
	server_address = ('0.0.0.0', 8080)
	httpd = HTTPServer(server_address, ProxyHTTPServer_RequestHandler)

	print('running server at %s:%d' % server_address)
	httpd.serve_forever()
