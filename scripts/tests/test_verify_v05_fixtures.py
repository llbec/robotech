import importlib.util
import json
from pathlib import Path
import unittest
from unittest.mock import patch
spec = importlib.util.spec_from_file_location('v05_fixtures', Path(__file__).parents[1] / 'verify-v0.5-fixtures.py')
module = importlib.util.module_from_spec(spec)
spec.loader.exec_module(module)

class FixtureIsolationTests(unittest.TestCase):
    def test_failure_cleans_only_created_database_and_container(self):
        calls = []
        def execute(argv, **kwargs):
            calls.append((argv, kwargs.get('input')))
            if argv[:3] == ['docker', 'compose', 'ps']:
                return 'postgres-id'
            if argv[:2] == ['docker', 'inspect']:
                return json.dumps([{'Config': {'Labels': {'com.docker.compose.project': 'robotech'}}, 'NetworkSettings': {'Networks': {'robotech_default': {}}}}])
            if argv[:2] == ['docker', 'run']:
                raise RuntimeError('fixture failure')
            return ''
        with patch.object(module, 'execute', execute), patch.object(module.Path, 'read_text', return_value='a'*64), patch.object(module.subprocess, 'run') as cleanup:
            with self.assertRaisesRegex(RuntimeError, 'fixture failure'):
                module.run()
        created = next(text for argv, text in calls if text and text.startswith('CREATE DATABASE '))
        dropped = next(text for argv, text in calls if text and text.startswith('DROP DATABASE '))
        name = created.split()[2].strip(';')
        self.assertTrue(name.startswith('robotech_verify_v05_'))
        self.assertIn(name, dropped)
        self.assertEqual(cleanup.call_args.args[0][:3], ['docker', 'rm', '-f'])
        self.assertTrue(cleanup.call_args.args[0][3].startswith('robotech-verify-v05-'))
    def test_failed_build_creates_no_database(self):
        with patch.object(module, 'execute', side_effect=RuntimeError('build failure')) as execute:
            with self.assertRaisesRegex(RuntimeError, 'build failure'):
                module.run()
        self.assertEqual(execute.call_count, 1)

if __name__ == '__main__':
    unittest.main()
