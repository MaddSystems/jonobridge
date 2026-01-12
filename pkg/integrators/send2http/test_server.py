import http.server
import json

class PayloadHandler(http.server.BaseHTTPRequestHandler):
    def do_POST(self):
        if self.path == '/receivers/json':
            content_length = int(self.headers['Content-Length'])
            post_data = self.rfile.read(content_length)
            try:
                data = json.loads(post_data.decode('utf-8'))
                print("Received payload:")
                print(json.dumps(data, indent=4))
            except json.JSONDecodeError as e:
                print("Error decoding JSON:", e)
                self.send_response(400)
                self.end_headers()
                return
            self.send_response(200)
            self.end_headers()
            self.wfile.write(b'OK')
        else:
            self.send_response(404)
            self.end_headers()

if __name__ == '__main__':
    server = http.server.HTTPServer(('localhost', 8095), PayloadHandler)
    print("Test server running on http://localhost:8095")
    print("Listening for POST requests to /receivers/json")
    try:
        server.serve_forever()
    except KeyboardInterrupt:
        print("\nServer stopped.")