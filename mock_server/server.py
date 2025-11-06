#!/usr/bin/env python3
import http.server
import json
import os
import urllib.parse
import uuid
from http import HTTPStatus

CONFIG_PATH = os.path.join(os.path.dirname(os.path.dirname(__file__)), 'config', 'settings.json')

class Handler(http.server.BaseHTTPRequestHandler):
    def _set_cors(self):
        self.send_header('Access-Control-Allow-Origin', '*')
        self.send_header('Access-Control-Allow-Methods', 'GET, POST, PUT, DELETE, OPTIONS')
        self.send_header('Access-Control-Allow-Headers', 'Content-Type')

    def do_OPTIONS(self):
        self.send_response(HTTPStatus.OK)
        self._set_cors()
        self.end_headers()

    def _read_config(self):
        try:
            with open(CONFIG_PATH, 'r', encoding='utf-8') as f:
                return json.load(f)
        except Exception:
            return {"servers": [], "log_level": "info", "auto_connect": False}

    # Simple runtime state to simulate a connected server and stats
    _connected_server_id = None
    _data_sent = 0
    _data_received = 0

    def _write_config(self, cfg):
        try:
            os.makedirs(os.path.dirname(CONFIG_PATH), exist_ok=True)
            with open(CONFIG_PATH, 'w', encoding='utf-8') as f:
                json.dump(cfg, f, ensure_ascii=False, indent=2)
            return True
        except Exception as e:
            print('Error writing config:', e)
            return False

    def _send_json(self, obj, status=200):
        data = json.dumps(obj, ensure_ascii=False).encode('utf-8')
        self.send_response(status)
        self._set_cors()
        self.send_header('Content-Type', 'application/json; charset=utf-8')
        self.send_header('Content-Length', str(len(data)))
        self.end_headers()
        self.wfile.write(data)

    def do_GET(self):
        parsed = urllib.parse.urlparse(self.path)
        # status and stats endpoints
        if parsed.path == '/api/status':
            if Handler._connected_server_id:
                self._send_json({'status': 'connected', 'server_id': Handler._connected_server_id})
            else:
                self._send_json({'status': 'disconnected'})
            return
        if parsed.path == '/api/stats':
            # simulate some traffic
            if Handler._connected_server_id:
                Handler._data_sent += 1024
                Handler._data_received += 2048
            self._send_json({'data_sent': Handler._data_sent, 'data_received': Handler._data_received})
            return
        if parsed.path == '/api/config':
            cfg = self._read_config()
            self._send_json(cfg)
            return
        if parsed.path == '/api/servers':
            cfg = self._read_config()
            self._send_json(cfg.get('servers', []))
            return
        # /api/servers/{id}
        if parsed.path.startswith('/api/servers/'):
            id = parsed.path[len('/api/servers/'):] or ''
            cfg = self._read_config()
            for s in cfg.get('servers', []):
                if s.get('id') == id or s.get('uuid') == id:
                    self._send_json(s)
                    return
            self._send_json({'error':'not found'}, status=404)
            return
        self.send_response(404)
        self.end_headers()

    def do_POST(self):
        parsed = urllib.parse.urlparse(self.path)
        # connect/disconnect endpoints for simulating client actions
        if parsed.path == '/api/connect':
            length = int(self.headers.get('content-length', 0))
            body = self.rfile.read(length) if length else b''
            try:
                data = json.loads(body.decode('utf-8')) if body else {}
            except Exception:
                self._send_json({'error': 'invalid json'}, status=400)
                return
            sid = data.get('server_id') or data.get('id') or data.get('uuid')
            if not sid:
                self._send_json({'error': 'no server_id provided'}, status=400)
                return
            # mark connected
            Handler._connected_server_id = sid
            Handler._data_sent = 0
            Handler._data_received = 0
            self._send_json({'status': 'connected', 'server_id': sid})
            return
        if parsed.path == '/api/disconnect':
            Handler._connected_server_id = None
            self._send_json({'status': 'disconnected'})
            return

        # allow existing POST handlers (servers) by delegating to original do_POST
        # fallthrough to existing server add handled earlier
        if parsed.path == '/api/servers':
            length = int(self.headers.get('content-length', 0))
            body = self.rfile.read(length) if length else b''
            try:
                s = json.loads(body.decode('utf-8')) if body else {}
            except Exception:
                self._send_json({'error':'invalid json'}, status=400)
                return
            # assign id
            s.setdefault('id', str(uuid.uuid4()))
            cfg = self._read_config()
            cfg.setdefault('servers', []).append(s)
            self._write_config(cfg)
            self._send_json(s, status=201)
            return

        self.send_response(404)
        self.end_headers()

    

    def do_PUT(self):
        parsed = urllib.parse.urlparse(self.path)
        # status/stats endpoints for queries
        if parsed.path == '/api/status':
            if Handler._connected_server_id:
                self._send_json({'status': 'connected', 'server_id': Handler._connected_server_id})
            else:
                self._send_json({'status': 'disconnected'})
            return
        if parsed.path == '/api/stats':
            # return simulated stats
            self._send_json({'data_sent': Handler._data_sent, 'data_received': Handler._data_received})
            return
        parsed = urllib.parse.urlparse(self.path)
        if parsed.path == '/api/config':
            length = int(self.headers.get('content-length', 0))
            body = self.rfile.read(length) if length else b''
            try:
                newcfg = json.loads(body.decode('utf-8'))
            except Exception:
                self._send_json({'error':'invalid json'}, status=400)
                return
            ok = self._write_config(newcfg)
            if not ok:
                self._send_json({'error':'could not write config'}, status=500)
                return
            self._send_json(newcfg)
            return
        if parsed.path.startswith('/api/servers/'):
            id = parsed.path[len('/api/servers/'):] or ''
            length = int(self.headers.get('content-length', 0))
            body = self.rfile.read(length) if length else b''
            try:
                s = json.loads(body.decode('utf-8')) if body else {}
            except Exception:
                self._send_json({'error':'invalid json'}, status=400)
                return
            cfg = self._read_config()
            servers = cfg.setdefault('servers', [])
            for i, existing in enumerate(servers):
                if existing.get('id') == id or existing.get('uuid') == id:
                    s.setdefault('id', existing.get('id'))
                    servers[i] = s
                    self._write_config(cfg)
                    self._send_json(s)
                    return
            self._send_json({'error':'not found'}, status=404)
            return
        self.send_response(404)
        self.end_headers()

    def do_DELETE(self):
        parsed = urllib.parse.urlparse(self.path)
        if parsed.path.startswith('/api/servers/'):
            id = parsed.path[len('/api/servers/'):] or ''
            cfg = self._read_config()
            servers = cfg.setdefault('servers', [])
            for i, existing in enumerate(servers):
                if existing.get('id') == id or existing.get('uuid') == id:
                    servers.pop(i)
                    self._write_config(cfg)
                    self.send_response(204)
                    self._set_cors()
                    self.end_headers()
                    return
            self._send_json({'error':'not found'}, status=404)
            return
        self.send_response(404)
        self.end_headers()

if __name__ == '__main__':
    port = 8080
    server = http.server.ThreadingHTTPServer(('0.0.0.0', port), Handler)
    print(f"mock_server (python): listening on http://0.0.0.0:{port}")
    try:
        server.serve_forever()
    except KeyboardInterrupt:
        server.server_close()
        print('server stopped')
