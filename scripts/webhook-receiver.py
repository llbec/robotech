#!/usr/bin/env python3
"""Isolated acceptance receiver. Control endpoints must stay on a test network."""
import argparse
import hashlib
import hmac
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer
import json
from pathlib import Path
import socket
import sqlite3
import threading
import time

class Receiver:
    def __init__(self, database, token):
        self.database = database
        self.token = token
        self.lock = threading.Lock()
        self.control = {'response_status': 200, 'drop_response_after_store_once': False, 'retry_after': None}
        with self.connect() as db:
            db.executescript('''CREATE TABLE IF NOT EXISTS receipts(event_id TEXT PRIMARY KEY, body BLOB NOT NULL, sha256 TEXT NOT NULL, received_at REAL NOT NULL);
            CREATE TABLE IF NOT EXISTS attempts(id INTEGER PRIMARY KEY, event_id TEXT, attempt TEXT, received_at REAL, status INTEGER, duplicate INTEGER, sha256 TEXT);
            ''')
    def connect(self):
        return sqlite3.connect(self.database, timeout=10)
    def receive(self, event, body, attempt):
        with self.lock, self.connect() as db:
            config = dict(self.control)
            status = config['response_status']
            digest = hashlib.sha256(body).hexdigest()
            previous = db.execute('SELECT sha256 FROM receipts WHERE event_id=?', (event,)).fetchone()
            duplicate = previous is not None
            if previous and previous[0] != digest:
                status = 409
            if 200 <= status < 300:
                db.execute('INSERT OR IGNORE INTO receipts VALUES(?,?,?,?)', (event, body, digest, time.time()))
            db.execute('INSERT INTO attempts(event_id,attempt,received_at,status,duplicate,sha256) VALUES(?,?,?,?,?,?)', (event, attempt, time.time(), status, duplicate, digest))
            drop = config['drop_response_after_store_once'] and not duplicate and 200 <= status < 300
            if drop:
                self.control['drop_response_after_store_once'] = False
            return status, drop, config['retry_after']
    def receipts(self):
        with self.connect() as db:
            records = [{'event_id': r[0], 'body': json.loads(r[1]), 'sha256': r[2], 'received_at': r[3]} for r in db.execute('SELECT * FROM receipts ORDER BY received_at')]
            attempts = [dict(zip(('id','event_id','attempt','received_at','status','duplicate','sha256'), r)) for r in db.execute('SELECT * FROM attempts ORDER BY id')]
        return {'receipts': records, 'attempts': attempts}

def handler(receiver):
    class Handler(BaseHTTPRequestHandler):
        def log_message(self, *_):
            pass
        def reply(self, status, data, retry_after=None):
            body = json.dumps(data, separators=(',', ':')).encode()
            self.send_response(status)
            self.send_header('Content-Type', 'application/json')
            self.send_header('Content-Length', str(len(body)))
            if retry_after is not None:
                self.send_header('Retry-After', str(retry_after))
            self.end_headers()
            self.wfile.write(body)
        def do_GET(self):
            if self.path == '/health':
                self.reply(200, {'status': 'ok'})
            elif self.path == '/receipts':
                self.reply(200, receiver.receipts())
            else:
                self.reply(404, {'error': 'not_found'})
        def do_POST(self):
            try:
                length = int(self.headers.get('Content-Length', '0'))
                if not 0 < length <= 1048576:
                    raise ValueError()
                body = self.rfile.read(length)
                data = json.loads(body)
                if self.path == '/control':
                    if not isinstance(data, dict) or set(data) - set(receiver.control):
                        raise ValueError()
                    if 'response_status' in data and (type(data['response_status']) is not int or not 200 <= data['response_status'] <= 599):
                        raise ValueError()
                    if 'drop_response_after_store_once' in data and type(data['drop_response_after_store_once']) is not bool:
                        raise ValueError()
                    if 'retry_after' in data and data['retry_after'] is not None and (not isinstance(data['retry_after'], (str, int)) or '\n' in str(data['retry_after']) or '\r' in str(data['retry_after'])):
                        raise ValueError()
                    with receiver.lock:
                        receiver.control.update(data)
                    return self.reply(200, {'status': 'ok'})
                if self.path != '/events':
                    return self.reply(404, {'error': 'not_found'})
                if not hmac.compare_digest(self.headers.get('Authorization', ''), 'Bearer ' + receiver.token):
                    return self.reply(401, {'error': 'authentication_required'})
                event = self.headers.get('X-Robotech-Event-Id')
                if not event or not isinstance(data, dict) or data.get('event_id') != event:
                    raise ValueError()
                status, drop, after = receiver.receive(event, body, self.headers.get('X-Robotech-Delivery-Attempt'))
                if drop:
                    self.close_connection = True
                    self.connection.shutdown(socket.SHUT_RDWR)
                    self.connection.close()
                else:
                    self.reply(status, {'event_id': event}, after)
            except (ValueError, json.JSONDecodeError):
                self.reply(400, {'error': 'invalid_request'})
    return Handler

def main():
    parser = argparse.ArgumentParser()
    parser.add_argument('--host', default='0.0.0.0')
    parser.add_argument('--port', type=int, default=8080)
    parser.add_argument('--database', required=True)
    parser.add_argument('--credential-file', required=True)
    args = parser.parse_args()
    token = Path(args.credential_file).read_text().strip()
    if len(token) < 32:
        parser.error('credential requires at least 32 characters')
    ThreadingHTTPServer((args.host, args.port), handler(Receiver(args.database, token))).serve_forever()
if __name__ == '__main__':
    main()
