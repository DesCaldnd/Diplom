#!/usr/bin/env python3
import argparse
import csv
import json
import math
import re
import subprocess
import sys
from pathlib import Path

BENCH_LINE_RE = re.compile(r'^(Benchmark\S+)\s+\d+\s+(\d+(?:\.\d+)?)\s+ns/op')


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description='Compare Go benchmark results and detect regressions.')
    parser.add_argument('--base', required=True, help='Path to benchmark output for base revision.')
    parser.add_argument('--head', required=True, help='Path to benchmark output for head revision.')
    parser.add_argument('--threshold-percent', type=float, required=True, help='Allowed slowdown threshold in percent.')
    parser.add_argument('--alpha', type=float, default=0.05, help='Alpha value for benchstat statistical significance check.')
    parser.add_argument('--json-output', required=True, help='Path to write machine-readable comparison result.')
    parser.add_argument('--markdown-output', required=True, help='Path to write markdown report.')
    return parser.parse_args()


def parse_benchmark_samples(path: Path) -> dict[str, list[float]]:
    benchmarks: dict[str, list[float]] = {}
    for raw_line in path.read_text(encoding='utf-8').splitlines():
        line = raw_line.strip()
        match = BENCH_LINE_RE.match(line)
        if not match:
            continue
        name = match.group(1)
        ns_per_op = float(match.group(2))
        benchmarks.setdefault(name, []).append(ns_per_op)
    return benchmarks


def geometric_mean(values: list[float]) -> float:
    if not values:
        raise ValueError('geometric_mean requires at least one value')
    return math.exp(sum(math.log(v) for v in values) / len(values))


def run_benchstat_csv(base_path: Path, head_path: Path, alpha: float) -> str:
    benchstat_binary = str(Path.home() / 'go' / 'bin' / 'benchstat')
    command = [
        benchstat_binary,
        '-format',
        'csv',
        '-alpha',
        str(alpha),
        str(base_path),
        str(head_path),
    ]
    completed = subprocess.run(command, check=True, capture_output=True, text=True)
    return completed.stdout + completed.stderr


def parse_change_percent(raw: str) -> float | None:
    value = raw.strip()
    if not value or value == '~':
        return None
    value = value.strip('%')
    if value.startswith('+'):
        value = value[1:]
    try:
        return float(value)
    except ValueError:
        return None


def parse_p_value(raw: str) -> float | None:
    value = raw.strip()
    if not value or value == '~':
        return None
    if value.startswith('p='):
        value = value.split()[0][2:]
    try:
        return float(value)
    except ValueError:
        return None


def extract_significance_map(benchstat_csv: str, alpha: float) -> dict[str, dict[str, object]]:
    significance: dict[str, dict[str, object]] = {}
    current_metric = ''
    reader = csv.reader(benchstat_csv.splitlines())
    for row in reader:
        if not row:
            continue
        first_cell = row[0].strip()
        if first_cell == '':
            if len(row) > 1 and row[1].strip() in {'sec/op', 'B/op', 'allocs/op'}:
                current_metric = row[1].strip()
            continue
        if first_cell in {'.unit', '.config', 'name', 'geomean', 'goos: darwin', 'goarch: arm64'} or first_cell.startswith('pkg: ') or first_cell.startswith('cpu: '):
            continue
        if first_cell.startswith(('B', 'D', 'F')) and len(row) == 1 and ':' in first_cell:
            continue
        if current_metric != 'sec/op':
            continue
        if len(row) < 7:
            continue
        benchmark = first_cell
        if benchmark.startswith('Benchmark'):
            benchmark = benchmark[len('Benchmark'):]
        change_raw = row[5].strip()
        p_value = parse_p_value(row[6])
        significance[benchmark] = {
            'change': change_raw,
            'p': p_value,
            'n1': row[1].strip(),
            'n2': row[3].strip(),
            'significant': bool(p_value is not None and p_value <= alpha),
            'change_percent': parse_change_percent(change_raw),
        }
    return significance


def format_percent(value: float) -> str:
    return f'{value:+.2f}%'


def build_report(base_samples: dict[str, list[float]], head_samples: dict[str, list[float]], significance_map: dict[str, dict[str, object]], threshold_percent: float) -> tuple[dict, str]:
    benchmark_names = sorted(set(base_samples) & set(head_samples))
    regressions: list[dict] = []
    rows: list[str] = []
    comparisons: list[dict] = []

    for benchmark_name in benchmark_names:
        base_values = base_samples[benchmark_name]
        head_values = head_samples[benchmark_name]
        base_mean = geometric_mean(base_values)
        head_mean = geometric_mean(head_values)
        ratio = head_mean / base_mean
        slowdown_percent = (ratio - 1.0) * 100.0
        significance = significance_map.get(benchmark_name.replace('Benchmark', ''), {})
        is_regression = slowdown_percent > threshold_percent and bool(significance.get('significant', False))
        comparison = {
            'benchmark': benchmark_name,
            'base_ns_per_op_geomean': base_mean,
            'head_ns_per_op_geomean': head_mean,
            'slowdown_percent': slowdown_percent,
            'ratio': ratio,
            'base_samples': base_values,
            'head_samples': head_values,
            'p_value': significance.get('p'),
            'statistically_significant': bool(significance.get('significant', False)),
            'regression': is_regression,
        }
        comparisons.append(comparison)
        status = '❌ regression' if is_regression else ('⚠️ slower' if slowdown_percent > 0 else '✅ faster/equal')
        rows.append(
            f'| `{benchmark_name}` | {base_mean:.0f} | {head_mean:.0f} | {format_percent(slowdown_percent)} | {comparison["p_value"] if comparison["p_value"] is not None else "n/a"} | {status} |'
        )
        if is_regression:
            regressions.append(comparison)

    missing_in_head = sorted(set(base_samples) - set(head_samples))
    missing_in_base = sorted(set(head_samples) - set(base_samples))

    payload = {
        'threshold_percent': threshold_percent,
        'benchmarks': comparisons,
        'regressions': regressions,
        'missing_in_head': missing_in_head,
        'missing_in_base': missing_in_base,
        'has_regressions': bool(regressions),
    }

    header = [
        '# Go benchmark comparison',
        '',
        f'- Threshold for failure: slower by more than **{threshold_percent:.2f}%**',
        '- Regression is reported only when slowdown is statistically significant according to `benchstat`.',
        '',
    ]

    if regressions:
        header.extend([
            '## Status',
            '',
            f'Found **{len(regressions)}** benchmark regression(s) above threshold.',
            '',
        ])
    else:
        header.extend([
            '## Status',
            '',
            'No benchmark regressions above threshold were found.',
            '',
        ])

    table = [
        '## Details',
        '',
        '| Benchmark | Base ns/op | Head ns/op | Delta | p-value | Status |',
        '| --- | ---: | ---: | ---: | ---: | --- |',
        *rows,
        '',
    ]

    extras: list[str] = []
    if missing_in_head:
        extras.extend([
            '## Missing in head run',
            '',
            *[f'- `{name}`' for name in missing_in_head],
            '',
        ])
    if missing_in_base:
        extras.extend([
            '## Missing in base run',
            '',
            *[f'- `{name}`' for name in missing_in_base],
            '',
        ])

    markdown = '\n'.join(header + table + extras)
    return payload, markdown


def main() -> int:
    args = parse_args()
    base_path = Path(args.base)
    head_path = Path(args.head)
    json_output = Path(args.json_output)
    markdown_output = Path(args.markdown_output)

    base_samples = parse_benchmark_samples(base_path)
    head_samples = parse_benchmark_samples(head_path)
    if not base_samples:
        raise SystemExit(f'No benchmark samples were parsed from {base_path}.')
    if not head_samples:
        raise SystemExit(f'No benchmark samples were parsed from {head_path}.')

    benchstat_csv = run_benchstat_csv(base_path, head_path, args.alpha)
    significance_map = extract_significance_map(benchstat_csv, args.alpha)
    payload, markdown = build_report(base_samples, head_samples, significance_map, args.threshold_percent)

    json_output.write_text(json.dumps(payload, indent=2, ensure_ascii=False) + '\n', encoding='utf-8')
    markdown_output.write_text(markdown, encoding='utf-8')

    return 0


if __name__ == '__main__':
    sys.exit(main())
