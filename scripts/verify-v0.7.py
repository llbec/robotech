#!/usr/bin/env python3
"""v0.7 spot acceptance; deployed accounts are read-only, fixtures use isolated resources."""
import argparse
import importlib.util
from pathlib import Path
from decimal import Decimal, localcontext
spec=importlib.util.spec_from_file_location('multi_verify',Path(__file__).with_name('verify-v0.6.py'))
module=importlib.util.module_from_spec(spec);spec.loader.exec_module(module)
common=module.common
class Spot(module.Multi):
    def fixtures(self):
        if not self.args.fixture_tests:
            raise common.Skip('使用 --fixture-tests 启用独立测试库固定现货样本')
        spec=importlib.util.spec_from_file_location('spot_fixtures',Path(__file__).with_name('verify-v0.7-fixtures.py'))
        helper=importlib.util.module_from_spec(spec);spec.loader.exec_module(helper)
        return helper.run()
    def spot_inventory(self):
        found=0
        for row in self.rows:
            account=row['account']
            spot=self.http('/api/v1/trade-events/stored',{'account':account,'instrument_type':'SPOT','limit':5})
            perp=self.http('/api/v1/trade-events/stored',{'account':account,'instrument_type':'PERPETUAL','limit':5})
            default=self.http('/api/v1/trade-events/stored',{'account':account,'limit':5})
            both=self.http('/api/v1/trade-events/stored',{'account':account,'instrument_type':'ALL','limit':5})
            self.require(spot['matched_records']+perp['matched_records']<=both['matched_records'],'分类库存计数矛盾；若持续写入请复验')
            self.require(all(t['payload']['instrument_type']=='PERPETUAL' for t in default['trades']+perp['trades']),'默认或合约查询混入现货')
            for f in spot['trades']:
                p=f['payload'];self.require(p['instrument_type']=='SPOT' and p['position_effect']=='NOT_APPLICABLE','现货字段错误')
                self.require(p['action']==p['side'] and p['side'] in ('BUY','SELL'),'现货买卖语义错误')
                with localcontext() as c:
                    c.prec=100
                    self.require(Decimal(p['price'])*Decimal(p['quantity'])==Decimal(p['notional']),'现货金额不精确')
                self.require(p['extension']['base_token']['tokenId'] and p['extension']['quote_token']['tokenId'],'缺失资产身份')
                self.require(not p['copy_eligible'] or p['fee_asset'] is not None,'未知费用仍可跟单')
                found+=1
            if self.args.live:
                live=self.http('/api/v1/trade-events/live',{'account':account,'instrument_type':'SPOT','limit':5})
                c=live['counts'];self.require(c['counts_version']==2,'计数版本错误')
                self.require(c['source_records']==sum(c[k] for k in ('duplicate_records','perpetual_records','spot_records','unsupported_records','invalid_records')),'分类计数不守恒')
                self.require(all(t['payload']['instrument_type']=='SPOT' for t in live['trades']),'live 现货筛选错误')
        if not found: raise common.Skip('分类接口已检查，但登记账户库存暂无现货；不能判定真实现货采集通过')
        return '抽样 {} 条现货；默认合约、类型筛选、资产身份和精确金额通过'.format(found)
    def validation(self):
        if not self.rows:raise common.Skip('没有登记账户')
        for path in ('/api/v1/trade-events/live','/api/v1/trade-events/stored'):
            for value in ('spot','INVALID',''):
                self.http(path,{'account':self.rows[0]['account'],'instrument_type':value},400)
        return '两类查询的非法 instrument_type 均返回 JSON 400'
    def candidates(self):
        events=self.sql("SELECT COALESCE(json_agg(q),'[]') FROM (SELECT event_id,status,payload FROM trade_log.outbox_events WHERE payload->'fact'->'payload'->>'instrument_type'='SPOT' ORDER BY created_at DESC LIMIT 5) q")
        keys={row['account_key'] for row in self.rows}
        events=[e for e in events if e['payload']['fact']['account_key'] in keys]
        for e in events:
            f=e['payload']['fact']['payload'];self.require(f['copy_eligible'] and f['trigger_type']=='USER' and f['fee_asset'],'队列含不合格现货')
            self.require(e['payload']['observation']['transport']=='WEBSOCKET','历史补录产生候选')
        if not events:raise common.Skip('暂无真实现货候选；固定样本验证准入不等于真实推送通过')
        return '抽样 {} 个现货候选身份及来源；接收方真实消费仍需手动对账'.format(len(events))
    def receipts(self):
        events=self.sql("SELECT COALESCE(json_agg(q),'[]') FROM (SELECT o.event_id,o.partition_key,o.body_sha256 FROM trade_log.outbox_events o JOIN trade_log.publishing_control p ON p.account_key=o.partition_key WHERE o.status='DELIVERED' AND o.payload->'fact'->'payload'->>'instrument_type'='SPOT' AND p.target_url='http://webhook-receiver:8080/events' ORDER BY o.delivered_at DESC LIMIT 5) q")
        keys={row['account_key'] for row in self.rows}
        events=[e for e in events if e['partition_key'] in keys]
        if not events: raise common.Skip('本地验收接收器暂无已送达现货事件，真实投递未验证')
        base=common.urllib.parse.urlsplit(self.base)
        hostname='['+base.hostname+']' if ':' in base.hostname else base.hostname
        url='http://'+hostname+':18080/receipts'
        if base.hostname not in ('127.0.0.1','localhost','::1'): raise common.Skip('请在服务器本机回环地址验收接收器')
        with self.opener.open(url,timeout=self.args.request_timeout) as response:
            receipts=common.json.load(response)['receipts']
        by_id={r['event_id']:r for r in receipts}
        checked=0
        for event in events:
            if event['event_id'] not in by_id:continue
            self.require(by_id[event['event_id']]['sha256']==event['body_sha256'],'接收记录与持久化现货投递正文摘要不一致')
            checked+=1
        if not checked:raise common.Skip('接收器记录可能已清空，暂无可匹配的现货接收证据')
        return '{} 个已送达现货事件与接收器正文摘要一致'.format(checked)
    def run(self):
        print('v0.7 现货验收；PASS/FAIL/SKIP，未验证项目不算通过',flush=True)
        for name,action in [('基础健康和版本',self.foundation),('账户清单',self.accounts),('现货类型参数',self.validation),('现货库存及分类',self.spot_inventory),('现货候选边界',self.candidates),('现货接收记录对账',self.receipts),('隔离现货固定样本',self.fixtures),('应用重启',self.lifecycle),('数据库故障恢复',self.database)]:self.check(name,action)
        print('结果：通过 {PASS} 项，失败 {FAIL} 项，跳过 {SKIP} 项'.format(**self.results))
        return int(bool(self.results['FAIL']))
if __name__=='__main__':
    parser=argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--base-url',default='http://127.0.0.1:8080')
    for name in ('fixture-tests','lifecycle','database-fault','live'):parser.add_argument('--'+name,action='store_true')
    parser.add_argument('--wait-seconds',type=int,default=180);parser.add_argument('--fault-seconds',type=int,default=15);parser.add_argument('--request-timeout',type=int,default=45)
    args=parser.parse_args()
    if not 1<=args.wait_seconds<=3600 or not 1<=args.fault_seconds<=300 or not 1<=args.request_timeout<=120:parser.error('等待时间或超时超出范围')
    url=common.urllib.parse.urlsplit(args.base_url)
    if url.scheme not in ('http','https') or not url.netloc or url.username is not None or url.password is not None or url.query or url.fragment:parser.error('base-url 必须是无凭证及查询参数的 HTTP(S) 地址')
    args.version='v0.7';args.account=None;args.expected_version=None
    raise SystemExit(Spot(args).run())
