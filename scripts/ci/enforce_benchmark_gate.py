#!/usr/bin/env python3
import argparse
import json
import sys
from pathlib import Path


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description='Fail workflow when benchmark regression gate is violated.')
    parser.add_argument('--json', required=True, help='Path to JSON report produced by compare_benchmarks.py')
    parser.add_argument('--override', default='false', help='Whether failing regressions are manually overridden')
    return parser.parse_args()


def as_bool(raw: str) -> bool:
    return str(raw).strip().lower() in {'1', 'true', 'yes', 'y', 'on'}


def main() -> int:
    args = parse_args()
    payload = json.loads(Path(args.json).read_text(encoding='utf-8'))
    regressions = payload.get('regressions', [])
    override = as_bool(args.override)

    if not regressions:
        print('No benchmark regressions detected.')
        return 0

    print(f'Found {len(regressions)} benchmark regression(s) above threshold.')
    for regression in regressions:
        print(
            f" - {regression['benchmark']}: +{regression['slowdown_percent']:.2f}% "
            f"(p={regression['p_value']})"
        )

    if override:
        print('Benchmark override is enabled, workflow will continue despite regressions.')
        return 0

    print('Benchmark regression gate failed. Re-run workflow_dispatch with benchmark_override=true to bypass manually.', file=sys.stderr)
    return 1


if __name__ == '__main__':
    sys.exit(main())
