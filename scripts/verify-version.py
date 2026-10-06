#!/usr/bin/env python3
"""Server acceptance shared by verify-v0.1.sh through verify-v0.4.sh.
Uses Python standard library, existing Compose containers and read-only SQL.
"""
import argparse
import datetime as dt
from decimal import Decimal, localcontext
import hashlib
import json
from pathlib import Path
import re
import subprocess
import sys
import tempfile
import time
import urllib.error
import urllib.parse
import urllib.request


class Skip(Exception):
    pass


class Verification:
    def __init__(self, args):
        self.args = args
        self.minor = int(args.version.rsplit('.', 1)[1])
        self.base = args.base_url.rstrip('/')
        self.opener = urllib.request.build_opener(urllib.request.ProxyHandler({}))
        self.results = {'PASS': 0, 'FAIL': 0, 'SKIP': 0}
        self.account = args.account.lower() if args.account else None
        self.network = None
        self.key = None
        self.trace_ids = set()
        self.live_result = None
        self.initial_status = None
        self.final_status = None
        self.container_ids = {}

    def record(self, state, name, detail=''):
        self.results[state] += 1
        print('[{}] {}{}'.format(state, name, '：' + detail if detail else ''), flush=True)

    def check(self, name, action):
        try:
            detail = action()
            self.record('PASS', name, detail or '')
        except Skip as exc:
            self.record('SKIP', name, str(exc))
        except Exception as exc:
            self.record('FAIL', name, str(exc))

    @staticmethod
    def require(condition, message):
        if not condition:
            raise AssertionError(message)

    def command(self, argv, timeout=30):
        p = subprocess.run(argv, stdout=subprocess.PIPE, stderr=subprocess.PIPE,
                           text=True, timeout=timeout)
        self.require(p.returncode == 0, '命令失败（退出码 {}）：{}'.format(
            p.returncode, p.stderr.strip()[:300]))
        return p.stdout

    def compose(self, *argv, timeout=30):
        return self.command(['docker', 'compose', *argv], timeout)

    def sql(self, expression):
        # Only SELECT expressions are supplied by this module. Never interpolate
        # unvalidated identifiers or a caller-supplied SQL expression.
        text = self.compose('exec', '-T', 'postgres', 'psql', '-U', 'robotech_admin',
                            '-d', 'robotech', '-qAt', '-v', 'ON_ERROR_STOP=1', '-c',
                            "BEGIN READ ONLY; SET LOCAL statement_timeout='15s'; "
                            + expression + '; COMMIT;', timeout=25)
        return json.loads(text.strip())

    def http(self, path, params=None, expected=200):
        if params:
            path += '?' + urllib.parse.urlencode(params)
        req = urllib.request.Request(self.base + path)
        try:
            response = self.opener.open(req, timeout=self.args.request_timeout)
        except urllib.error.HTTPError as exc:
            response = exc
        with response:
            status, headers, raw = response.status, response.headers, response.read()
        self.require(status == expected, 'HTTP {}，期望 {}（{}）'.format(status, expected, path.split('?')[0]))
        self.require(headers.get_content_type() == 'application/json', '响应不是 JSON')
        body = json.loads(raw)
        trace = headers.get('X-Trace-Id')
        self.require(isinstance(trace, str) and bool(trace) and trace not in self.trace_ids,
                     '响应 trace 缺失或重复')
        self.trace_ids.add(trace)
        if status == 200:
            self.require(body['meta']['schema_version'] == 1 and body['meta']['trace_id'] == trace,
                         '成功信封或 trace 不一致')
            return body['data']
        self.require(body.get('trace_id') == trace and isinstance(body.get('message'), str),
                     '错误信封或 trace 不一致')
        return body

    @staticmethod
    def timestamp(value):
        parsed = dt.datetime.fromisoformat(value.replace('Z', '+00:00'))
        if parsed.tzinfo is None:
            raise AssertionError('时间缺少时区')
        return parsed

    def deployment(self):
        services = ['query-api', 'trade-log-query']
        if self.minor >= 2:
            services.append('postgres')
        if self.minor >= 3:
            services.append('trade-collector')
        for service in services:
            ids = self.compose('ps', '-q', service).split()
            self.require(len(ids) == 1, service + ' 需要恰好一个运行容器')
            info = json.loads(self.command(['docker', 'inspect', ids[0]]))[0]
            self.require(info['State']['Running'], service + ' 未运行')
            self.container_ids[service] = ids[0]
            if service != 'query-api':
                self.require(not any(info['NetworkSettings'].get('Ports', {}).values()),
                             service + ' 不应发布宿主机端口')
            if service != 'postgres':
                self.require(info['Config']['User'] not in ('', '0', 'root', '0:0'), service + ' 不应以 root 运行')
                self.require(info['HostConfig']['ReadonlyRootfs'], service + ' 根文件系统应只读')
                configs = [m for m in info['Mounts'] if m['Destination'].startswith('/etc/robotech/')]
                self.require(configs and all(not m['RW'] for m in configs), service + ' 配置应只读挂载')
        return '、'.join(services)

    def foundation(self):
        self.require(self.http('/api/v1/health')['status'] == 'ok', '网关健康异常')
        v = self.http('/api/v1/version')['version']
        self.require(re.fullmatch(r'\d+\.\d+\.\d+(?:[-+][\w.-]+)?', v), '版本格式错误')
        if self.args.expected_version:
            self.require(v == self.args.expected_version, '版本不符，实际 ' + v)
        else:
            major, minor = [int(x) for x in v.split('.')[:2]]
            self.require(major > 0 or minor >= self.minor, '运行版本低于被验收版本')
        return '运行版本 ' + v + '；仅验收 ' + self.args.version + ' 能力'

    def choose_account(self):
        if self.minor >= 3:
            items = self.http('/api/v1/watch-accounts')['items']
            if not self.account:
                self.require(len(items) == 1, '不能自动选择账户，请指定 --account')
                self.account = items[0]['account']
            selected = [s for s in items if s['account'] == self.account]
            self.require(len(selected) == 1, '指定账户不在采集配置中')
            self.initial_status = selected[0]
            self.network = selected[0]['network']
        elif not self.account:
            raise AssertionError('v0.1/v0.2 必须通过 --account 指定验收账户')
        self.require(re.fullmatch(r'0x[0-9a-f]{40}', self.account or ''), '账户格式错误')
        if self.network:
            self.require(self.network in ('mainnet', 'testnet'), '网络不支持')
            self.key = 'hyperliquid:{}:hyperliquid:{}'.format(self.network, self.account)
        return self.account

    def validation(self):
        self.require(self.account, '账户选择失败')
        cases = [[], [('account', 'invalid')], [('account', self.account), ('limit', '0')],
                 [('account', self.account), ('limit', '2001')],
                 [('account', self.account), ('account', self.account)],
                 [('account', self.account), ('unknown', '1')]]
        if self.minor >= 2:
            for extra in ({'source': 'bad'}, {'source': 'stored', 'cursor': 'invalid'},
                          {'source': 'stored', 'start_time': '2026-10-06T00:00:00Z',
                           'end_time': '2026-10-06T00:00:00Z'},
                          {'source': 'stored', 'start_time': '2026-10-06T00:00:00'}):
                cases.append([('account', self.account), *extra.items()])
        for params in cases:
            result = self.http('/api/v1/trade-events', params, 400)
            self.require(result['code'] == 'VALIDATION_ERROR', '非法参数错误码不符')
        return '{} 个非法参数请求均返回 JSON 400'.format(len(cases))

    def facts(self, trades):
        identities = set()
        for f in trades:
            self.require(f['fact_id'] not in identities, '响应中出现重复 fact_id')
            identities.add(f['fact_id'])
            self.require(f['account'] == self.account and f['fact_type'] == 'TRADE', '事实账户或类型错误')
            self.require(f['revision'] >= 1 and f['source'] == 'OFFICIAL_API', '事实版本或来源错误')
            self.timestamp(f['occurred_at'])
            p = f['payload']
            self.require(p['instrument_type'] == 'PERPETUAL' and p['side'] in ('BUY', 'SELL'), '市场类型或方向错误')
            with localcontext() as context:
                context.prec = 100
                self.require(Decimal(p['quantity']) > 0 and Decimal(p['price']) * Decimal(p['quantity']) == Decimal(p['notional']),
                             '数量或精确成交额不一致')
            self.require(isinstance(p['copy_eligible'], bool), 'copy_eligible 类型错误')
        return identities

    def live(self):
        self.require(self.account, '账户选择失败')
        if self.minor != 1 and not self.args.live:
            raise Skip('使用 --live 执行官方查询，会保存采集证据及成交')
        data = self.http('/api/v1/trade-events', {'account': self.account, 'limit': '5'})
        self.require(data['account'] == self.account and data['network'] in ('mainnet', 'testnet'), 'live 账户或网络错误')
        if self.network:
            self.require(data['network'] == self.network, 'live 与配置账户网络不一致')
        c = data['counts']
        self.require(all(type(v) is int and v >= 0 for v in c.values()), '分类计数必须为非负整数')
        self.require(c['source_records'] == sum(c[k] for k in ('duplicate_records', 'spot_records', 'unsupported_records',
                                                             'invalid_records', 'perpetual_records')), '来源分类计数不守恒')
        self.require(c['returned_records'] == len(data['trades']) == min(c['perpetual_records'], 5), '展示数量不符')
        if c['source_records'] >= 2000:
            self.require(data['coverage'] == 'LIMITED' and 'SOURCE_RECORD_LIMIT' in data['warnings'], '来源上限未披露')
        self.require(data['display_truncated'] == (c['perpetual_records'] > len(data['trades'])), '展示截断标志不符')
        self.require(data['evidence_ref'] == data['query_id'], '证据引用不一致')
        self.require(re.fullmatch(r'query_[A-Za-z0-9_]+', data['query_id']), 'query_id 不合法')
        self.facts(data['trades'])
        self.live_result = data
        if not data['trades']:
            raise Skip('查询成功但没有合约成交；不能认定真实成交验收通过，请更换 --account')
        if self.minor >= 2:
            p = data['persistence']
            self.require(p['status'] == 'COMMITTED' and p['inserted_records'] + p['existing_records'] == c['perpetual_records'],
                         '数据库持久化回执不符')
        return 'query_id={}；展示 {} 条；完整合约数 {}'.format(data['query_id'], len(data['trades']), c['perpetual_records'])

    def repeated_live(self):
        if not self.live_result or not self.live_result['trades']:
            raise Skip('没有可比较的首次 live 成交')
        first = self.live_result
        try:
            self.live()
            second = self.live_result
        finally:
            self.live_result = first
        original = {f['fact_id']: f for f in first['trades']}
        repeated = {f['fact_id']: f for f in second['trades']}
        overlap = original.keys() & repeated.keys()
        if not overlap:
            raise Skip('来源窗口已变化，两次展示没有共同交易；不能判定稳定身份')
        for identity in overlap:
            a, b = original[identity], repeated[identity]
            self.require(a['occurred_at'] == b['occurred_at'] and a['revision'] == b['revision'], '重复成交时间或 revision 变化')
            self.require(all(a['payload'][k] == b['payload'][k] for k in ('side', 'market')), '重复成交方向或市场变化')
            self.require(all(Decimal(a['payload'][k]) == Decimal(b['payload'][k]) for k in ('price', 'quantity', 'fee')), '重复成交金额变化')
        return '{} 条共同成交的身份、版本及核心字段一致'.format(len(overlap))

    def evidence(self):
        if not self.live_result:
            raise Skip('未执行或未成功执行 live 查询')
        query = self.live_result['query_id']
        config = self.compose('exec', '-T', 'trade-log-query', 'cat', '/etc/robotech/trade-log.toml')
        section = re.search(r'(?ms)^\[evidence\]\s*(.*?)(?=^\[|\Z)', config)
        self.require(section is not None, '无法确定证据目录')
        match = re.search(r'^directory\s*=\s*"([^"]+)"', section[1], re.M)
        self.require(match is not None, '无法读取证据目录')
        directory = match[1].rstrip('/') + '/' + query
        self.require('trade-log-query' in self.container_ids, '部署检查未找到查询容器')
        with tempfile.TemporaryDirectory(prefix='robotech-verify-') as target:
            self.command(['docker', 'cp', self.container_ids['trade-log-query'] + ':' + directory, target], timeout=60)
            root = Path(target) / query
            manifest = json.loads((root / 'manifest.json').read_text())
            result = json.loads((root / 'result.json').read_text())
            self.require(manifest['status'] == 'COMPLETED' and manifest['query_id'] == query, '证据任务未完成')
            self.require(len(result['trades']) == self.live_result['counts']['perpetual_records'], '证据只保存了展示集合')
            files = list((root / 'metadata').glob('*.json'))
            self.require(files, '缺少原始响应元数据')
            for entry in files:
                meta = json.loads(entry.read_text())
                raw = (root / 'responses' / (meta['raw_log_id'] + '.body')).read_bytes()
                self.require(hashlib.sha256(raw).hexdigest() == meta['sha256'], '原始响应摘要不符')
        return '完整事实与 {} 份原始响应摘要一致'.format(len(files))

    def stored(self):
        self.require(self.account, '账户选择失败')
        data = self.http('/api/v1/trade-events', {'account': self.account, 'source': 'stored', 'limit': '5'})
        self.require(data['coverage'] == 'STORED_RECORDS_ONLY' and data['query_scope'] == 'STORED_TIME_RANGE', '库存范围契约错误')
        self.require(data['returned_records'] == len(data['trades']) <= 5 and data['matched_records'] >= len(data['trades']), '库存计数错误')
        self.facts(data['trades'])
        if self.network:
            self.require(self.network == data['network'], '库存与采集状态网络不一致')
        self.network = data['network']
        self.require(self.network in ('mainnet', 'testnet'), '库存网络错误')
        self.key = 'hyperliquid:{}:hyperliquid:{}'.format(self.network, self.account)
        count = self.sql("SELECT count(*) FROM trade_log.account_facts_current WHERE account_key='{}' AND fact_type='TRADE' AND NOT is_retracted AND ingest_seq<={}".format(self.key, int(data['snapshot_seq'])))
        self.require(count == data['matched_records'], '数据库快照计数与接口不一致')
        return '快照库存 {} 条；不是账户完整历史'.format(count)

    def paging(self):
        self.require(self.key, '库存检查未确定账户网络')
        params = {'account': self.account, 'source': 'stored', 'limit': '2'}
        first = self.http('/api/v1/trade-events', params)
        if not first['has_more']:
            raise Skip('库存不足以覆盖跨页验证')
        self.require(first['next_cursor'], 'has_more=true 但缺少游标')
        second = self.http('/api/v1/trade-events', dict(params, cursor=first['next_cursor']))
        self.require(first['snapshot_seq'] == second['snapshot_seq'] and first['matched_records'] == second['matched_records'], '分页快照或总量改变')
        self.require(not self.facts(first['trades']) & self.facts(second['trades']), '跨页重复')
        combined = first['trades'] + second['trades']
        keys = [(self.timestamp(f['occurred_at']).timestamp(), int(f['source_ref']), f['fact_id']) for f in combined]
        self.require(keys == sorted(keys, key=lambda k: (-k[0], -k[1], k[2])), '分页顺序错误')
        if combined:
            start = combined[-1]['occurred_at']
            end = combined[0]['occurred_at']
            if self.timestamp(start) < self.timestamp(end):
                bounded = self.http('/api/v1/trade-events', dict(params, start_time=start, end_time=end))
                self.require(all(self.timestamp(start) <= self.timestamp(f['occurred_at']) < self.timestamp(end) for f in bounded['trades']), '半开时间区间错误')
        return '快照一致、跨页无重复、时间范围及排序符合约定'

    def integrity(self):
        self.require(self.key, '库存检查未确定账户网络')
        duplicate = self.sql("SELECT count(*) FROM (SELECT v.payload->'payload'->>'market',c.source_tid FROM trade_log.account_facts_current c JOIN trade_log.account_fact_versions v ON v.fact_id=c.fact_id AND v.revision=c.current_revision WHERE c.account_key='{}' AND NOT c.is_retracted GROUP BY 1,2 HAVING count(*)>1) q".format(self.key))
        self.require(duplicate == 0, '存在同市场同 tid 的重复事实')
        if self.live_result and self.live_result['trades']:
            ids = ','.join("'{}'".format(f['fact_id']) for f in self.live_result['trades'])
            self.require(all(re.fullmatch(r'hl_fill_v1_[0-9a-f]{64}', f['fact_id']) for f in self.live_result['trades']), 'fact_id 格式错误')
            found = self.sql("SELECT COALESCE(json_agg(v.payload),'[]'::json) FROM trade_log.account_facts_current c JOIN trade_log.account_fact_versions v ON v.fact_id=c.fact_id AND v.revision=c.current_revision WHERE c.account_key='{}' AND c.fact_id IN ({})".format(self.key, ids))
            self.require(self.facts(found) == self.facts(self.live_result['trades']), 'live 返回事实未全部入库')
            saved = {f['fact_id']: f for f in found}
            for f in self.live_result['trades']:
                old = saved[f['fact_id']]
                self.require(old['occurred_at'] == f['occurred_at'] and all(Decimal(old['payload'][k]) == Decimal(f['payload'][k]) for k in ('price', 'quantity', 'notional', 'fee')), '库存与 live 核心事实不一致')
        return '同账户同市场同 tid 无重复；有 live 结果时核对其入库事实'

    def status(self):
        self.require(self.account, '账户选择失败')
        selected = [s for s in self.http('/api/v1/watch-accounts')['items'] if s['account'] == self.account]
        self.require(len(selected) == 1, '状态账户不唯一')
        return selected[0]

    def observe(self):
        initial = self.status()
        old_watermark, old_success = initial['scanned_through'], initial['last_success_at']
        old_pong = initial.get('websocket', {}).get('last_pong_at')
        deadline = time.monotonic() + self.args.wait_seconds
        next_report = 0
        while True:
            current = self.status()
            self.final_status = current
            progressed = (current['scanned_through'] and current['last_success_at'] and
                          (not old_watermark or self.timestamp(current['scanned_through']) > self.timestamp(old_watermark)) and
                          (not old_success or self.timestamp(current['last_success_at']) > self.timestamp(old_success)))
            ws_ready = True
            if self.minor >= 4:
                w, r = current['websocket'], current['recovery']
                self.require(w['enabled'], 'WebSocket 未启用，无法验收 v0.4 双通道')
                self.require(r['status'] != 'BLOCKED', '恢复缺口 BLOCKED：' + json.dumps(r.get('last_error'), ensure_ascii=False))
                ws_ready = (w['status'] == 'LIVE' and w['last_pong_at'] and
                            (not old_pong or self.timestamp(w['last_pong_at']) > self.timestamp(old_pong)) and r['open_gap_count'] == 0)
            if progressed and ws_ready and current['consecutive_failures'] == 0 and current['last_error'] is None:
                return 'HTTP 水位及成功时间推进' + ('；WS pong 更新，未完成缺口为 0' if self.minor >= 4 else '')
            remaining = deadline - time.monotonic()
            self.require(remaining > 0, '等待超时：HTTP 状态={}，水位={}，失败次数={}{}'.format(
                current['status'], current['scanned_through'], current['consecutive_failures'],
                '，WS/恢复=' + json.dumps(current.get('websocket', {}).get('status')) + '/' + json.dumps(current.get('recovery', {}).get('status')) if self.minor >= 4 else ''))
            if time.monotonic() >= next_report:
                print('[WAIT] 等待自动采集推进，剩余约 {} 秒'.format(int(remaining)), flush=True)
                next_report = time.monotonic() + 15
            time.sleep(min(3, remaining))

    def reparse(self, transport='HTTP'):
        if self.minor == 1:
            if not self.live_result:
                raise Skip('没有成功 live 任务可重放')
            query = self.live_result['query_id']
        else:
            self.require(self.key, '库存检查未确定账户网络')
            query = self.sql("SELECT COALESCE(to_json((SELECT query_id FROM trade_log.collection_jobs WHERE account='{}' AND network='{}' AND status='COMPLETED' {} ORDER BY updated_at DESC LIMIT 1)), 'null'::json)".format(
                self.account, self.network, ("AND transport='{}'".format(transport) if self.minor >= 4 else "AND source_id='official_http'") + (" AND job_origin='COLLECTOR'" if self.minor >= 3 else '')))
            if not query:
                raise Skip('没有已完成的 {} 任务'.format(transport))
        self.require(re.fullmatch(r'query_[A-Za-z0-9_]+', query), '重放 query_id 格式错误')
        output = self.compose('exec', '-T', 'trade-log-query', 'trade-log-query', 'reparse', '--query-id', query, timeout=90)
        report = json.loads(output.strip())
        self.require(report['comparison'] == 'SAME', '重放结果不是 SAME')
        return query + ' comparison=SAME'

    def ws_evidence(self):
        self.require(self.key, '库存检查未确定账户网络')
        data = self.sql("SELECT COALESCE(json_agg(q),'[]'::json) FROM (SELECT s.session_id,s.status,s.received_messages,s.committed_messages,j.query_id,j.message_mode FROM trade_log.collection_stream_sessions s LEFT JOIN LATERAL (SELECT query_id,message_mode FROM trade_log.collection_jobs WHERE session_id=s.session_id AND status='COMPLETED' ORDER BY message_sequence DESC LIMIT 1) j ON true WHERE s.checkpoint_key='{}' ORDER BY s.created_at DESC LIMIT 1) q".format(self.key))
        self.require(data and data[0]['status'] == 'LIVE', '当前数据库 WS 会话不是 LIVE')
        self.require(data[0]['query_id'], '当前 WS 会话尚无成功提交消息')
        self.require(data[0]['message_mode'] in ('UNKNOWN', 'SNAPSHOT', 'LIVE_UPDATE'), 'WS 消息模式不合法')
        both = self.sql("SELECT count(*) FROM (SELECT o.fact_id FROM trade_log.fact_observations o JOIN trade_log.raw_logs r ON r.id=o.raw_log_id JOIN trade_log.account_facts_current c ON c.fact_id=o.fact_id WHERE c.account_key='{}' GROUP BY o.fact_id HAVING count(DISTINCT r.transport)=2) q".format(self.key))
        if not both:
            raise Skip('WS 消息已提交，但暂无同时被 HTTP/WS 观察的交易，跨通道去重待验证')
        return '{} 条事实保留 HTTP 与 WS 双方证据；UNKNOWN 不是失败'.format(both)

    def lifecycle(self):
        if not self.args.lifecycle:
            raise Skip('使用 --lifecycle 显式启用服务重启；会短暂中断对应能力')
        if self.results['FAIL']:
            raise Skip('前置检查失败，不执行停止服务测试')
        service = 'trade-collector' if self.minor >= 3 else 'trade-log-query'
        self.require(service in self.container_ids, '部署检查未确认目标容器，不执行停止测试')
        baseline = self.status() if self.minor >= 3 else None
        before = self.sql("SELECT COALESCE(json_agg(fact_id),'[]'::json) FROM (SELECT fact_id FROM trade_log.account_facts_current WHERE account_key='{}' ORDER BY occurred_at DESC LIMIT 5) q".format(self.key)) if self.minor >= 2 and self.key else []
        try:
            self.compose('stop', '-t', '15', service, timeout=30)
            cid = self.container_ids.get(service)
            self.require(cid, '没有已验证的目标容器')
            info = json.loads(self.command(['docker', 'inspect', cid]))[0]
            self.require(not info['State']['Running'] and info['State']['ExitCode'] == 0, '服务未正常退出')
            self.require(self.http('/api/v1/health')['status'] == 'ok', '内部服务停止影响网关存活')
            if self.minor >= 3:
                self.http('/api/v1/watch-accounts', expected=503)
                self.http('/api/v1/trade-events', {'account': self.account, 'source': 'stored', 'limit': '1'})
        finally:
            self.compose('start', service, timeout=30)
        self.wait_restored(baseline, before)
        return service + ' 正常退出并恢复' + ('；抽样旧事实保留' if before else '；无库存样本，旧事实保留未验证') + ('；HTTP 水位未回退' if baseline else '')

    def wait_restored(self, baseline, before):
        deadline = time.monotonic() + self.args.wait_seconds
        error = None
        while time.monotonic() < deadline:
            try:
                self.require(self.http('/api/v1/health')['status'] == 'ok', '网关未恢复')
                if self.minor == 1:
                    self.require(self.live_result, '缺少重启前已完成任务')
                    self.reparse()
                if self.minor >= 2:
                    self.http('/api/v1/trade-events', {'account': self.account, 'source': 'stored', 'limit': '1'})
                    if before:
                        self.require(all(re.fullmatch(r'hl_fill_v1_[0-9a-f]{64}', f) for f in before), '抽样事实 ID 不合法')
                        ids = ','.join("'{}'".format(f) for f in before)
                        n = self.sql("SELECT count(*) FROM trade_log.account_facts_current WHERE account_key='{}' AND fact_id IN ({})".format(self.key, ids))
                        self.require(n == len(before), '重启前抽样事实丢失')
                if baseline:
                    current = self.status()
                    self.require(current['last_success_at'] and current['last_success_at'] != baseline['last_success_at'], '等待新的成功轮次')
                    if baseline['scanned_through']:
                        self.require(self.timestamp(current['scanned_through']) >= self.timestamp(baseline['scanned_through']), '水位回退')
                    self.require(current['consecutive_failures'] == 0 and current['last_error'] is None, '故障未清除')
                    if self.minor >= 4:
                        self.require(current['websocket']['status'] == 'LIVE' and current['websocket']['last_pong_at'] and current['recovery']['open_gap_count'] == 0,
                                     '等待 WS 心跳与缺口恢复')
                return
            except Exception as exc:
                error = exc
                print('[WAIT] 等待服务恢复：' + str(exc), flush=True)
                time.sleep(3)
        raise AssertionError('恢复等待超时：' + str(error))

    def database_fault(self):
        if not self.args.database_fault:
            raise Skip('使用 --database-fault 显式启用停库测试；会影响所有使用该数据库的服务')
        if self.results['FAIL']:
            raise Skip('前置检查失败，不执行停库测试')
        self.require(self.key, '库存检查未确定账户网络')
        self.require('postgres' in self.container_ids, '部署检查未确认数据库容器，不执行停库测试')
        baseline = self.status() if self.minor >= 3 else None
        before = self.sql("SELECT COALESCE(json_agg(fact_id),'[]'::json) FROM (SELECT fact_id FROM trade_log.account_facts_current WHERE account_key='{}' ORDER BY occurred_at DESC LIMIT 5) q".format(self.key))
        try:
            self.compose('stop', '-t', '15', 'postgres', timeout=30)
            self.http('/api/v1/trade-events', {'account': self.account, 'source': 'stored', 'limit': '1'}, 503)
            self.require(self.http('/api/v1/health')['status'] == 'ok', '数据库停止影响网关存活')
            if self.minor >= 3:
                self.http('/api/v1/watch-accounts', expected=503)
            print('[WAIT] 数据库保持停止 {} 秒'.format(self.args.fault_seconds), flush=True)
            time.sleep(self.args.fault_seconds)
        finally:
            self.compose('start', 'postgres', timeout=30)
        self.wait_restored(baseline, before)
        return '停库期间业务 503、网关存活' + ('；抽样旧事实保留' if before else '；无库存样本，旧事实保留未验证') + ('，采集水位推进及错误清除' if baseline else '')

    def run(self):
        print(self.args.version + ' 服务器验收；PASS/FAIL/SKIP，未覆盖项目不算通过', flush=True)
        self.check('基础健康、版本与信封', self.foundation)
        self.check('Compose 运行、内部端口及挂载', self.deployment)
        self.check('验收账户', self.choose_account)
        self.check('参数错误及 trace 契约', self.validation)
        self.check('官方 live 查询、计数及精确金额', self.live)
        if self.minor == 1 or self.args.live:
            self.check('重复 live 查询的事实身份稳定', self.repeated_live)
            self.check('完整文件证据及原始摘要', self.evidence)
        if self.minor >= 2:
            self.check('stored 契约及数据库快照计数', self.stored)
            self.check('库存分页及时间区间', self.paging)
            self.check('事实去重与 live 入库核对', self.integrity)
        if self.minor >= 3:
            self.check('观察 HTTP 推进' + ('、WS 心跳及缺口恢复' if self.minor >= 4 else ''), self.observe)
        self.check('HTTP 任务重新解析', self.reparse)
        if self.minor >= 4:
            self.check('WS 持久化及跨通道证据', self.ws_evidence)
            self.check('WS 任务重新解析', lambda: self.reparse('WEBSOCKET'))
        self.check('应用停止及恢复', self.lifecycle)
        if self.minor >= 2:
            self.check('数据库停止及恢复', self.database_fault)
        self.record('SKIP', '来源限流、业务冲突、饱和分页及真实历史完整性', '受控开发测试或独立交易明细对账；本脚本不伪造来源故障')
        if self.minor >= 4:
            self.record('SKIP', '订阅失败、pong 超时、队列溢出与关闭 WS', '固定开发测试或文档手动步骤覆盖；本脚本不修改部署配置')
        print('\n结果：通过 {PASS} 项，失败 {FAIL} 项，跳过 {SKIP} 项'.format(**self.results))
        print('结论：' + ('不通过' if self.results['FAIL'] else '已执行项目通过；存在跳过项，非完整验收'))
        return 1 if self.results['FAIL'] else 0


def main():
    parser = argparse.ArgumentParser(description='从项目根目录执行；使用现有容器，不构建或部署。')
    parser.add_argument('--version', required=True, choices=['v0.1', 'v0.2', 'v0.3', 'v0.4'])
    parser.add_argument('--base-url', default='http://127.0.0.1:8080')
    parser.add_argument('--account', help='验收账户；v0.1/v0.2 必填，v0.3/v0.4 默认读取配置账户')
    parser.add_argument('--expected-version', help='可选精确核对，如 0.4.0；默认允许后续兼容版本')
    parser.add_argument('--live', action='store_true', help='主动查询官方来源并保存结果；v0.1 默认执行')
    parser.add_argument('--lifecycle', action='store_true', help='停止并启动内部查询服务或采集服务，会中断对应能力')
    parser.add_argument('--database-fault', action='store_true', help='验收环境显式停库并恢复，影响全部数据库调用')
    parser.add_argument('--fault-seconds', type=int, default=15)
    parser.add_argument('--wait-seconds', type=int, default=90, help='每段采集或恢复观察上限，默认 90 秒')
    parser.add_argument('--request-timeout', type=int, default=60)
    args = parser.parse_args()
    url = urllib.parse.urlsplit(args.base_url)
    if url.scheme not in ('http', 'https') or not url.netloc or url.username or url.query or url.fragment:
        parser.error('--base-url 必须是无凭证和查询参数的 HTTP(S) 地址')
    if args.account and not re.fullmatch(r'0x[0-9a-fA-F]{40}', args.account):
        parser.error('--account 地址格式错误')
    if args.version in ('v0.1', 'v0.2') and not args.account:
        parser.error('v0.1/v0.2 请通过 --account 指定实际账户')
    if args.version == 'v0.1' and args.database_fault:
        parser.error('v0.1 不包含数据库故障测试')
    if not 1 <= args.fault_seconds <= 60 or not 5 <= args.wait_seconds <= 3600 or not 1 <= args.request_timeout <= 120:
        parser.error('fault-seconds 范围 1–60，wait-seconds 范围 5–3600，request-timeout 范围 1–120')
    try:
        return Verification(args).run()
    except KeyboardInterrupt:
        print('\n[FAIL] 验收被中断；如启用了停止测试，请确认对应服务已恢复', file=sys.stderr)
        return 130


if __name__ == '__main__':
    sys.exit(main())
