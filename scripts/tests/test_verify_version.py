"""Contract and failure-path tests; no Docker or external account required."""
import argparse
import contextlib
import copy
import importlib.util
import io
import json
from pathlib import Path
import unittest
from unittest.mock import patch

spec = importlib.util.spec_from_file_location('verify_version', Path(__file__).parents[1] / 'verify-version.py')
module = importlib.util.module_from_spec(spec)
spec.loader.exec_module(module)
ROOT = Path(__file__).resolve().parents[2]
FIXTURE = json.loads((ROOT / 'tests/fixtures/v0.2/import-evidence/query_fixture/result.json').read_text())


def verification(version='v0.4', **kwargs):
    values = dict(version=version, account=FIXTURE['account'], base_url='http://127.0.0.1:8080',
                  expected_version=None, live=False, lifecycle=False, database_fault=False,
                  fault_seconds=1, wait_seconds=5, request_timeout=5)
    values.update(kwargs)
    v = module.Verification(argparse.Namespace(**values))
    v.network = 'mainnet'
    v.key = 'hyperliquid:mainnet:hyperliquid:' + v.account
    return v


class AcceptanceTests(unittest.TestCase):
    def test_real_fixture_and_corrupt_amount(self):
        v = verification()
        self.assertEqual(len(v.facts(FIXTURE['trades'])), 1)
        bad = copy.deepcopy(FIXTURE['trades'])
        bad[0]['payload']['notional'] = '0'
        with self.assertRaisesRegex(AssertionError, '成交额'):
            v.facts(bad)

    def test_live_preserves_counts_and_empty_is_skip(self):
        v = verification('v0.1')
        v.http = lambda *args, **kwargs: copy.deepcopy(FIXTURE)
        self.assertIn('query_fixture', v.live())
        broken = copy.deepcopy(FIXTURE)
        broken['counts']['source_records'] += 1
        v.http = lambda *args, **kwargs: broken
        with self.assertRaisesRegex(AssertionError, '不守恒'):
            v.live()
        empty = copy.deepcopy(FIXTURE)
        empty['trades'] = []
        empty['counts'] = {k: 0 for k in empty['counts']}
        v.http = lambda *args, **kwargs: empty
        with self.assertRaises(module.Skip):
            v.live()

    def test_stored_database_count_mismatch_is_failure(self):
        v = verification('v0.2')
        data = dict(coverage='STORED_RECORDS_ONLY', query_scope='STORED_TIME_RANGE',
                    returned_records=1, matched_records=1, trades=FIXTURE['trades'],
                    network='mainnet', snapshot_seq='1')
        v.http = lambda *args, **kwargs: data
        v.sql = lambda statement: 2
        with self.assertRaisesRegex(AssertionError, '快照计数'):
            v.stored()

    def test_duplicate_page_and_changed_snapshot_are_failures(self):
        v = verification('v0.2')
        first = dict(has_more=True, next_cursor='cursor', snapshot_seq='1', matched_records=5,
                     trades=FIXTURE['trades'])
        v.http = lambda *args, **kwargs: first
        with self.assertRaisesRegex(AssertionError, '跨页重复'):
            v.paging()
        second = dict(first, snapshot_seq='2', trades=[])
        replies = iter([first, second])
        v.http = lambda *args, **kwargs: next(replies)
        with self.assertRaisesRegex(AssertionError, '快照'):
            v.paging()

    def test_quiet_ws_requires_new_pong_and_http_progress(self):
        v = verification()
        initial = dict(scanned_through='2026-10-06T00:00:00Z', last_success_at='2026-10-06T00:00:01Z',
                       consecutive_failures=0, last_error=None, status='WAITING',
                       websocket=dict(enabled=True, status='LIVE', last_pong_at='2026-10-06T00:00:00Z'),
                       recovery=dict(status='HTTP_SCANNED', open_gap_count=0))
        final = copy.deepcopy(initial)
        final.update(scanned_through='2026-10-06T00:00:30Z', last_success_at='2026-10-06T00:00:31Z')
        final['websocket']['last_pong_at'] = '2026-10-06T00:00:20Z'
        replies = iter([initial, initial, final])
        v.status = lambda: next(replies)
        with patch.object(module.time, 'sleep'), contextlib.redirect_stdout(io.StringIO()):
            self.assertIn('pong 更新', v.observe())

    def test_blocked_gap_and_disabled_ws_are_not_passes(self):
        for enabled, recovery, expected in [(True, 'BLOCKED', 'BLOCKED'), (False, 'DISABLED', '未启用')]:
            v = verification()
            v.status = lambda: dict(scanned_through=None, last_success_at=None,
                                    websocket=dict(enabled=enabled, last_pong_at=None),
                                    recovery=dict(status=recovery, last_error={'code': 'VERSION_CONFLICT'}))
            with self.assertRaisesRegex(AssertionError, expected):
                v.observe()

    def test_fault_failure_still_starts_database(self):
        v = verification('v0.2', database_fault=True)
        commands = []
        v.container_ids['postgres'] = 'test-postgres'
        v.sql = lambda statement: []
        v.compose = lambda *args, **kwargs: commands.append(args)
        def failure(*args, **kwargs):
            raise AssertionError('未返回预期 503')
        v.http = failure
        with self.assertRaisesRegex(AssertionError, '503'):
            v.database_fault()
        self.assertEqual(commands[-1], ('start', 'postgres'))

    def test_faults_are_opt_in(self):
        v = verification()
        v.compose = lambda *args, **kwargs: self.fail('默认不应停止服务')
        with self.assertRaises(module.Skip):
            v.lifecycle()
        with self.assertRaises(module.Skip):
            v.database_fault()

    def test_summary_fails_on_real_error_and_keeps_skip_distinct(self):
        v = verification('v0.1')
        for name in ('foundation', 'deployment', 'choose_account', 'validation', 'live', 'repeated_live', 'evidence', 'reparse'):
            setattr(v, name, lambda: None)
        v.validation = lambda: (_ for _ in ()).throw(AssertionError('验证失败'))
        output = io.StringIO()
        with contextlib.redirect_stdout(output):
            self.assertEqual(v.run(), 1)
        self.assertIn('[FAIL]', output.getvalue())
        self.assertIn('[SKIP]', output.getvalue())
        self.assertIn('结论：不通过', output.getvalue())

    def test_publishing_disabled_is_skip_and_bad_account_fails(self):
        v = verification('v0.5')
        data = dict(account=v.account, account_key=v.key, status='DISABLED', enabled=False, outbox={})
        v.http = lambda *a, **k: data
        with self.assertRaises(module.Skip):
            v.publishing()
        data['account_key'] = 'wrong'
        with self.assertRaisesRegex(AssertionError, '账户'):
            v.publishing()

    def test_outbox_checksum_is_checked_and_empty_is_skip(self):
        v = verification('v0.5')
        values = iter([0, []])
        v.sql = lambda *a: next(values)
        with self.assertRaises(module.Skip):
            v.outbox()
        values = iter([0, [dict(event_id='fact:1', wire_body='7b7d', body_sha256='bad', payload={})]])
        v.sql = lambda *a: next(values)
        with self.assertRaisesRegex(AssertionError, '摘要'):
            v.outbox()

    def test_publisher_failure_always_restarts(self):
        v = verification('v0.5', lifecycle=True)
        v.container_ids['trade-parser-publisher'] = 'test-publisher'
        v.http = lambda *a, **k: dict(activated_at='time', activation_epoch=1)
        commands = []
        v.publishing_queue_snapshot = lambda: []
        v.compose = lambda *a, **k: commands.append(a)
        v.command = lambda *a, **k: json.dumps([{'State': {'Running': False, 'ExitCode': 4}}])
        with self.assertRaisesRegex(AssertionError, '正常退出'):
            v.publishing_lifecycle()
        self.assertEqual(commands[-1], ('start', 'trade-parser-publisher'))


if __name__ == '__main__':
    unittest.main()
