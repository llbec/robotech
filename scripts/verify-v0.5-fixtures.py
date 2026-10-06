#!/usr/bin/env python3
"""Build and run fixed samples in a database owned only by this invocation."""
import json
import os
from pathlib import Path
import re
import subprocess
import tempfile
import uuid

def execute(argv, *, input=None, capture=False, timeout=3600):
    result = subprocess.run(argv, input=input, text=True, stdout=subprocess.PIPE if capture else None,
                            stderr=subprocess.PIPE if capture else None, timeout=timeout)
    if result.returncode:
        raise RuntimeError('隔离测试命令失败，退出码 {}{}'.format(result.returncode, '：' + result.stderr.strip()[:300] if capture and result.stderr else ''))
    return result.stdout if capture else None

def run():
    identity = uuid.uuid4().hex
    database = 'robotech_verify_v05_' + identity
    container = 'robotech-verify-v05-' + identity
    created = False
    psql = ['docker', 'compose', 'exec', '-T', 'postgres', 'psql', '-qAt', '-v', 'ON_ERROR_STOP=1', '-U', 'robotech_admin', '-d', 'robotech']
    # Build before creating a database, so lengthy downloads leave no test resources.
    execute(['docker', 'build', '--target', 'verify-v05', '-f', 'gateway/query-api/Dockerfile', '-t', 'robotech/verify-v05:0.5.0', '.'])
    postgres = execute(['docker', 'compose', 'ps', '-q', 'postgres'], capture=True).strip()
    if not postgres:
        raise RuntimeError('PostgreSQL 容器未运行')
    info = json.loads(execute(['docker', 'inspect', postgres], capture=True))[0]
    project = info['Config']['Labels']['com.docker.compose.project']
    network = project + '_default'
    if network not in info['NetworkSettings']['Networks']:
        raise RuntimeError('找不到当前 Compose 的默认网络')
    password = Path('secrets/postgres-password').read_text().strip()
    if not re.fullmatch(r'[0-9a-fA-F]+', password):
        raise RuntimeError('管理员凭证格式不符合初始化脚本约定')
    try:
        created = True
        execute(psql, input='CREATE DATABASE {};\n'.format(database), capture=True)
        with tempfile.TemporaryDirectory(prefix='robotech-v05-') as temporary:
            secret = Path(temporary) / 'database-url'
            secret.write_text('postgresql://robotech_admin:{}@postgres:5432/{}\n'.format(password, database))
            os.chmod(secret, 0o600)
            execute(['docker', 'run', '--rm', '--name', container, '--network', network, '--read-only', '--tmpfs', '/tmp:rw,size=64m',
                     '--mount', 'type=bind,source={},target=/run/secrets/fixture-database-url,readonly'.format(secret.resolve()),
                     'robotech/verify-v05:0.5.0'])
    finally:
        subprocess.run(['docker', 'rm', '-f', container], stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)
        if created:
            execute(psql, input='DROP DATABASE IF EXISTS {} WITH (FORCE);\n'.format(database), capture=True)
    return '独立数据库和接收器：候选、抑制、500/429/401、响应丢失、重复、租约恢复；测试资源已清理'
if __name__ == '__main__':
    print(run())
