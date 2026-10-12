#!/usr/bin/env python3
"""Read deployment state; mutations are confined to optional isolated fixtures."""
import argparse
import importlib.util
from pathlib import Path
import time

spec = importlib.util.spec_from_file_location('verify_common', Path(__file__).with_name('verify-version.py'))
common = importlib.util.module_from_spec(spec)
spec.loader.exec_module(common)

class Multi(common.Verification):
    def accounts(self):
        items, cursor = [], None
        while True:
            page = self.http('/api/v1/monitoring/accounts/list', {'limit': 100, **({'cursor': cursor} if cursor else {})})
            self.require(isinstance(page.get('items'), list), '列表结构错误')
            self.require(0 <= page['active_count'] <= page['max_active_accounts'] <= 10, '活动容量不合法')
            items.extend(page['items'])
            cursor = page.get('next_cursor')
            if not cursor:
                break
        self.require(len({r['id'] for r in items}) == len(items), '分页有重复账户')
        self.rows = items
        for row in items:
            self.require(row['desired_state'] in ('ACTIVE', 'PAUSED') and row['generation'] >= 1, '账户状态或代次错误')
            detail = self.http('/api/v1/monitoring/accounts/detail', {'id': row['id']})
            self.require(detail['account_key'] == row['account_key'], '详情账户串账')
        return '{} 个登记账户，{} 个活动账户'.format(len(items), sum(r['desired_state']=='ACTIVE' for r in items))

    def progression(self):
        active = [r for r in self.rows if r['desired_state'] == 'ACTIVE']
        if len(active) < 2:
            raise common.Skip('至少需要两个真实活动账户；脚本不会自动添加部署账户')
        baseline = {r['id']: (r.get('collection') or {}).get('scanned_through') for r in active}
        pongs = {r['id']: ((r.get('collection') or {}).get('websocket') or {}).get('last_pong_at') for r in active}
        end = time.monotonic() + self.args.wait_seconds
        progressed = set()
        while time.monotonic() < end:
            for row in active:
                current = self.http('/api/v1/monitoring/accounts/detail', {'id': row['id']})
                collection = current.get('collection') or {}
                websocket = collection.get('websocket') or {}
                ws_ready = not websocket.get('enabled') or (websocket.get('status') == 'LIVE' and websocket.get('last_pong_at') and websocket['last_pong_at'] != pongs[row['id']])
                if ws_ready and collection.get('scanned_through') and collection['scanned_through'] != baseline[row['id']]:
                    self.require(collection.get('consecutive_failures') == 0, '水位推进后仍有采集失败')
                    progressed.add(row['id'])
            if len(progressed) == len(active):
                return '{} 个账户 HTTP 水位独立推进；不代表来源历史完整'.format(len(active))
            time.sleep(2)
        raise AssertionError('等待超时；未推进账户数 {}'.format(len(active)-len(progressed)))

    def publications(self):
        summary = self.http('/api/v1/publishing/status')
        self.require(isinstance(summary.get('items'), list), '发布汇总格式错误')
        keys = {row['account_key'] for row in self.rows}
        for item in summary['items']:
            self.require(item['account_key'] in keys, '发布账户未登记')
            registration = next(row for row in self.rows if row['account_key'] == item['account_key'])
            detail = self.http('/api/v1/monitoring/accounts/publishing-status', {'id': registration['id']})
            self.require(detail['account_key'] == item['account_key'], '发布详情串账')
        if not any(item['enabled'] for item in summary['items']):
            raise common.Skip('发布未启用；汇总与分账户格式已检查，真实投递未验证')
        return '{} 个账户发布状态可查询；真实接收证据仍用 v0.5 验收'.format(len(summary['items']))

    def fixtures(self):
        if not self.args.fixture_tests:
            raise common.Skip('使用 --fixture-tests 启用隔离双账户故障与暂停恢复测试')
        spec = importlib.util.spec_from_file_location('v06_fixtures', Path(__file__).with_name('verify-v0.6-fixtures.py'))
        module = importlib.util.module_from_spec(spec)
        spec.loader.exec_module(module)
        return module.run()

    def lifecycle(self):
        if not self.args.lifecycle:
            raise common.Skip('使用 --lifecycle 启用采集和发布进程重启')
        if self.results['FAIL']:
            raise common.Skip('前置检查失败，不重启服务')
        before = {r['id']: (r['desired_state'], r['generation']) for r in self.rows}
        publishing = {r['account_key']: (r['activation_epoch'], r['activated_at']) for r in self.http('/api/v1/publishing/status')['items']}
        queue = self.sql("SELECT COALESCE(json_agg(q),'[]'::json) FROM (SELECT event_id,body_sha256,payload->>'expires_at' expires_at FROM trade_log.outbox_events ORDER BY created_at DESC LIMIT 100) q")
        self.compose('restart', 'trade-collector', 'trade-parser-publisher', timeout=90)
        end = time.monotonic()+self.args.wait_seconds
        while time.monotonic()<end:
            try:
                self.accounts()
                self.require(before == {r['id']: (r['desired_state'], r['generation']) for r in self.rows}, '重启改变账户状态或代次')
                after = {r['account_key']: (r['activation_epoch'], r['activated_at']) for r in self.http('/api/v1/publishing/status')['items']}
                self.require(publishing == after, '重启改变发布激活边界')
                self.publishing_queue_preserved(queue)
                return '账户清单、暂停状态、代次和发布激活边界保留'
            except Exception:
                time.sleep(2)
        raise AssertionError('进程重启恢复超时或状态不一致')

    def database(self):
        if not self.args.database_fault:
            raise common.Skip('使用 --database-fault 启用停库测试')
        if self.results['FAIL']:
            raise common.Skip('前置检查失败，不停止数据库')
        self.accounts()
        before = {r['id']: (r['desired_state'], r['generation'], (r.get('collection') or {}).get('scanned_through')) for r in self.rows}
        queue = self.sql("SELECT COALESCE(json_agg(q),'[]'::json) FROM (SELECT event_id,body_sha256,payload->>'expires_at' expires_at FROM trade_log.outbox_events ORDER BY created_at DESC LIMIT 100) q")
        try:
            self.compose('stop', 'postgres', timeout=90)
            self.http('/api/v1/health')
            self.http('/api/v1/monitoring/accounts/list', expected=503)
            self.http('/api/v1/publishing/status', expected=503)
            print('[WAIT] 数据库停止 {} 秒'.format(self.args.fault_seconds), flush=True)
            time.sleep(self.args.fault_seconds)
        finally:
            self.compose('start', 'postgres', timeout=90)
        end = time.monotonic()+self.args.wait_seconds
        while time.monotonic()<end:
            try:
                self.accounts()
                self.require(set(before) == {r['id'] for r in self.rows}, '账户清单丢失')
                for row in self.rows:
                    state, generation, watermark = before[row['id']]
                    self.require((state, generation) == (row['desired_state'], row['generation']), '管理状态丢失')
                    current = (row.get('collection') or {}).get('scanned_through')
                    if watermark:
                        self.require(current and self.timestamp(current)>=self.timestamp(watermark), '水位回退')
                self.publishing_queue_preserved(queue)
                return '共享依赖故障明确返回 503；恢复后清单、状态、水位及抽样 outbox 保留'
            except Exception:
                time.sleep(2)
        raise AssertionError('数据库恢复超时')

    def run(self):
        print('v0.6 多账户验收；部署账户仅查看，固定故障使用独立测试库', flush=True)
        for name, action in [('基础健康和版本', self.foundation), ('账户列表、详情和容量', self.accounts),
                             ('多账户持续更新', self.progression), ('分账户发布状态', self.publications),
                             ('隔离双账户固定样本', self.fixtures), ('进程重启', self.lifecycle), ('数据库故障恢复', self.database)]:
            self.check(name, action)
        print('结果：通过 {PASS} 项，失败 {FAIL} 项，跳过 {SKIP} 项'.format(**self.results))
        print('结论：'+('不通过' if self.results['FAIL'] else ('通过' if not self.results['SKIP'] else '已执行项目通过；存在未验证项')))
        return int(bool(self.results['FAIL']))

if __name__ == '__main__':
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--base-url', default='http://127.0.0.1:8080')
    parser.add_argument('--fixture-tests', action='store_true')
    parser.add_argument('--lifecycle', action='store_true')
    parser.add_argument('--database-fault', action='store_true')
    parser.add_argument('--wait-seconds', type=int, default=180)
    parser.add_argument('--fault-seconds', type=int, default=15)
    parser.add_argument('--request-timeout', type=int, default=30)
    args = parser.parse_args()
    if not 1<=args.wait_seconds<=3600 or not 1<=args.fault_seconds<=300 or not 1<=args.request_timeout<=120:
        parser.error('等待时间或超时参数超出范围')
    url=common.urllib.parse.urlsplit(args.base_url)
    if url.scheme not in ('http','https') or not url.netloc or url.username or url.query or url.fragment:
        parser.error('base-url 必须是无凭证和查询参数的 HTTP(S) 地址')
    args.version='v0.6'; args.account=None; args.expected_version=None
    raise SystemExit(Multi(args).run())
