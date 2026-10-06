#!/bin/sh
# Server acceptance for the v0.0 foundation, including later compatible releases.
set -eu
if ! command -v python3 >/dev/null 2>&1; then
  echo '[FAIL] 缺少 python3（仅使用标准库，无需 pip 安装）' >&2
  exit 2
fi
exec python3 - "$@" <<'PY'
import argparse
import datetime
import json
import re
import shutil
import subprocess
import sys
import time
import urllib.error
import urllib.parse
import urllib.request

parser = argparse.ArgumentParser(description='v0.0 服务器验收；默认不停止或重启服务。')
parser.add_argument('--base-url', default='http://127.0.0.1:8080')
parser.add_argument('--expected-version', help='可选：严格核对运行版本，如 0.0.0；默认允许后续兼容版本')
parser.add_argument('--lifecycle', action='store_true', help='允许停止并启动 query-api，短暂中断网关；仅在验收环境启用')
args = parser.parse_args()
url = urllib.parse.urlsplit(args.base_url)
if url.scheme not in ('http', 'https') or not url.netloc or url.query or url.fragment or url.username:
    parser.error('--base-url 必须是无凭证、查询参数和片段的 HTTP(S) 地址')
base = args.base_url.rstrip('/')
opener = urllib.request.build_opener(urllib.request.ProxyHandler({}))
results = {'PASS': 0, 'FAIL': 0, 'SKIP': 0}
traces = []

def record(state, name, detail=''):
    results[state] += 1
    print('[{}] {}{}'.format(state, name, '：' + detail if detail else ''), flush=True)

def check(name, action):
    try:
        detail = action()
        record('PASS', name, detail or '')
    except Exception as exc:
        record('FAIL', name, str(exc))

def require(condition, message):
    if not condition:
        raise AssertionError(message)

def command(argv, timeout=20, input=None):
    return subprocess.run(argv, input=input, text=True, stdout=subprocess.PIPE,
                          stderr=subprocess.PIPE, timeout=timeout)

def docker(argv, **kwargs):
    p = command(['docker'] + argv, **kwargs)
    require(p.returncode == 0, 'Docker 命令失败：' + p.stderr.strip()[:300])
    return p.stdout

def request(path, method='GET', headers=None):
    req = urllib.request.Request(base + path, method=method, headers=headers or {})
    try:
        response = opener.open(req, timeout=10)
    except urllib.error.HTTPError as exc:
        response = exc
    with response:
        return response.status, response.headers, response.read()

def response_json(path, expected=200, method='GET', code=None, headers=None):
    status, h, body = request(path, method, headers)
    require(status == expected, 'HTTP {}，期望 {}'.format(status, expected))
    require(h.get_content_type() == 'application/json', 'Content-Type 不是 application/json')
    data = json.loads(body)
    trace = h.get('X-Trace-Id')
    require(isinstance(trace, str) and re.fullmatch(r'trace_[0-9a-f]{32}', trace), 'trace ID 格式错误')
    require(trace not in traces, '不同请求使用了重复 trace ID')
    traces.append(trace)
    if code:
        require(data.get('code') == code, '错误码不符合约定')
        require(isinstance(data.get('message'), str) and data['message'], '缺少错误说明')
        require(data.get('trace_id') == trace, '错误体与响应头 trace ID 不一致')
    else:
        require(data['meta']['trace_id'] == trace, '响应体与响应头 trace ID 不一致')
        require(data['meta']['schema_version'] == 1, 'schema_version 不为 1')
        require(data['data']['service'] == 'query-api', 'service 不为 query-api')
    return data, h

health = {}
version = {}

def health_check():
    data, _ = response_json('/api/v1/health')
    require(data['data']['status'] == 'ok', '健康状态不是 ok')
    started = data['data']['started_at']
    require(datetime.datetime.fromisoformat(started.replace('Z', '+00:00')).tzinfo is not None,
            'started_at 必须包含时区')
    health['started_at'] = started

def stable_start():
    data, _ = response_json('/api/v1/health', headers={'X-Trace-Id': 'client_supplied_trace'})
    require(data['data']['started_at'] == health.get('started_at'), '运行期间 started_at 发生变化')
    require(data['meta']['trace_id'] != 'client_supplied_trace', '直接采用了客户端 trace ID')

def version_check():
    data, _ = response_json('/api/v1/version')
    v = data['data']['version']
    require(isinstance(v, str) and re.fullmatch(r'\d+\.\d+\.\d+(?:[-+][\w.-]+)?', v), '版本格式错误')
    if args.expected_version:
        require(v == args.expected_version, '实际版本 {}，期望 {}'.format(v, args.expected_version))
    version['value'] = v
    return '运行版本 ' + v

def method_check(path):
    _, h = response_json(path, 405, 'POST', 'METHOD_NOT_ALLOWED')
    require({'GET', 'HEAD'} <= set(h.get('Allow', '').replace(' ', '').split(',')), 'Allow 缺少 GET 或 HEAD')

def head_check(path):
    status, h, body = request(path, 'HEAD')
    require(status == 200 and body == b'', 'HEAD 状态不为 200 或存在响应体')
    require(h.get_content_type() == 'application/json', 'HEAD Content-Type 错误')
    trace = h.get('X-Trace-Id', '')
    require(re.fullmatch(r'trace_[0-9a-f]{32}', trace) and trace not in traces, 'HEAD trace ID 错误或重复')
    traces.append(trace)

print('v0.0 基础能力验收（不代表后续版本业务验收或历史完整性）', flush=True)
check('健康接口及 JSON 信封', health_check)
check('启动时间固定、客户端 trace 不覆盖服务端 trace', stable_start)
check('版本接口及 JSON 信封', version_check)
check('未知路由返回 JSON 404', lambda: response_json('/api/v1/__verify_v00_missing__', 404, code='RESOURCE_NOT_FOUND') and None)
for path in ('/api/v1/health', '/api/v1/version'):
    check('POST ' + path + ' 返回 405 及 Allow', lambda p=path: method_check(p))
    check('HEAD ' + path + ' 无响应体', lambda p=path: head_check(p))

container = {}
def deployment():
    require(shutil.which('docker'), '缺少 Docker 命令')
    ids = docker(['compose', 'ps', '-q', 'query-api']).split()
    require(len(ids) == 1, '需要恰好一个运行中的 Compose query-api 容器；请从项目目录执行')
    info = json.loads(docker(['inspect', ids[0]]))[0]
    require(info['State']['Running'], 'query-api 未运行')
    mounts = [m for m in info['Mounts'] if m['Destination'] == '/etc/robotech/query-api.toml']
    require(len(mounts) == 1 and not mounts[0]['RW'], '配置文件未按约定只读挂载')
    require(info['HostConfig']['ReadonlyRootfs'], '容器根文件系统不是只读')
    container.update(id=ids[0], image=info['Image'])
    return '容器运行，配置只读挂载，根文件系统只读'
check('Docker 部署与配置挂载', deployment)

if container:
    def cli_version():
        actual = docker(['exec', container['id'], 'query-api', '--version']).strip()
        require(actual == 'query-api ' + version.get('value', ''), 'CLI 版本与 HTTP 版本不一致')
    check('CLI 与 HTTP 版本一致', cli_version)

    def missing_config():
        p = command(['docker', 'run', '--rm', '--pull=never', '--network=none', container['image'],
                     '--config', '/__verify_v00_missing__.toml'])
        require(p.returncode == 2, '缺失配置退出码为 {}，期望 2'.format(p.returncode))
    check('隔离容器：缺失配置拒绝启动', missing_config)

    basic_config = ('config_version = 1\n[server]\nhost = "0.0.0.0"\nport = 8080\n'
                    'shutdown_timeout_seconds = 10\n[logging]\nlevel = "info"\nformat = "json"\n')
    invalid_configs = {
        '不支持的配置版本': basic_config.replace('config_version = 1', 'config_version = 99'),
        '缺少必填端口': basic_config.replace('port = 8080\n', ''),
        '非法端口': basic_config.replace('port = 8080', 'port = 0'),
        '非法日志级别': basic_config.replace('level = "info"', 'level = "invalid"'),
        '未知配置字段': basic_config.replace('[server]', 'unknown_field = true\n[server]'),
        'TOML 格式错误': '[server',
    }
    def invalid_config(body):
        p = command(['docker', 'run', '--rm', '--pull=never', '--network=none', '-i',
                     container['image'], '--config', '/dev/stdin'], input=body)
        require(p.returncode == 2, '非法配置退出码为 {}，期望 2'.format(p.returncode))
    for name, body in invalid_configs.items():
        check('隔离容器：' + name + '拒绝启动', lambda b=body: invalid_config(b))

    def request_logs():
        # Application logs are written to stderr; docker logs retains that stream.
        p = command(['docker', 'logs', '--since=5m', container['id']])
        require(p.returncode == 0, '读取日志失败')
        text = p.stdout + p.stderr
        events = []
        for line in text.splitlines():
            try:
                events.append(json.loads(line))
            except ValueError:
                pass
        matching = [e for e in events if e.get('fields', {}).get('trace_id') in traces]
        require(matching, '未找到本次请求的 JSON 日志；核对 URL 对应容器、日志级别及 format=json')
        for event in matching:
            f = event['fields']
            require(f.get('message') == 'request_completed' and f.get('service') == 'query-api', '请求日志结构错误')
            require(all(k in f for k in ('method', 'route', 'status', 'duration_ms')), '请求日志缺少字段')
        return '匹配 {} 条请求日志'.format(len(matching))
    check('响应 trace 与请求日志关联', request_logs)
else:
    for name in ('CLI 与 HTTP 版本一致', '隔离容器：缺失配置拒绝启动', '响应 trace 与请求日志关联'):
        record('SKIP', name, 'Docker 部署检查未通过')

if args.lifecycle and container:
    def lifecycle():
        stopped = False
        try:
            stopped = True
            docker(['compose', 'stop', '-t', '15', 'query-api'], timeout=30)
            info = json.loads(docker(['inspect', container['id']]))[0]
            require(not info['State']['Running'] and info['State']['ExitCode'] == 0, 'SIGTERM 未正常退出（退出码非 0）')
            p = command(['docker', 'logs', '--tail=50', container['id']])
            require('shutdown_completed' in p.stdout + p.stderr, '未找到正常退出日志')
        finally:
            if stopped:
                docker(['compose', 'start', 'query-api'], timeout=30)
        deadline = time.monotonic() + 30
        while time.monotonic() < deadline:
            try:
                data, _ = response_json('/api/v1/health')
                require(data['data']['status'] == 'ok', '重启后健康状态错误')
                require(data['data']['started_at'] != health.get('started_at'), '重启后启动时间未更新')
                return 'SIGTERM 退出码 0，重新启动后健康检查通过'
            except Exception:
                time.sleep(1)
        raise AssertionError('启动后 30 秒内未恢复健康接口')
    check('停止及重新启动 query-api（会中断网关）', lifecycle)
else:
    record('SKIP', 'SIGTERM 停止及重新启动', '使用 --lifecycle 显式启用；会短暂中断网关')
record('SKIP', 'SIGINT、在途请求 draining 与退出超时', '需受控开发测试，本脚本不宣称已覆盖')
record('SKIP', '源码构建及环境覆盖优先级', '本脚本验收部署结果，这些项目由开发测试覆盖')
print('\n结果：通过 {PASS} 项，失败 {FAIL} 项，跳过 {SKIP} 项'.format(**results))
print('结论：' + ('不通过' if results['FAIL'] else '已执行项目全部通过；存在未执行项目，非完整验收'))
sys.exit(1 if results['FAIL'] else 0)
PY
