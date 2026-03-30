#!/usr/bin/env python3
from __future__ import annotations

import argparse
import html
import json
import re
import xml.etree.ElementTree as ET
from pathlib import Path


def parse_go_json(path: Path | None) -> dict:
    if not path or not path.exists():
        return {
            "cases_total": 0,
            "passed": 0,
            "failed": 0,
            "skipped": 0,
            "packages_total": 0,
            "packages_passed": 0,
            "packages_failed": 0,
            "failed_tests": [],
        }

    tests: dict[tuple[str, str], str] = {}
    packages: dict[str, str] = {}

    with path.open("r", encoding="utf-8", errors="ignore") as fh:
        for line in fh:
            line = line.strip()
            if not line.startswith("{"):
                continue
            try:
                event = json.loads(line)
            except json.JSONDecodeError:
                continue

            pkg = event.get("Package", "unknown")
            action = event.get("Action")
            test = event.get("Test")

            if test and action in {"pass", "fail", "skip"}:
                tests[(pkg, test)] = action
            elif not test and action in {"pass", "fail"}:
                packages[pkg] = action

    failed_tests = [f"{pkg} :: {test}" for (pkg, test), status in tests.items() if status == "fail"]

    return {
        "cases_total": len(tests),
        "passed": sum(1 for status in tests.values() if status == "pass"),
        "failed": sum(1 for status in tests.values() if status == "fail"),
        "skipped": sum(1 for status in tests.values() if status == "skip"),
        "packages_total": len(packages),
        "packages_passed": sum(1 for status in packages.values() if status == "pass"),
        "packages_failed": sum(1 for status in packages.values() if status == "fail"),
        "failed_tests": failed_tests[:10],
    }


def parse_coverage(path: Path | None) -> str:
    if not path or not path.exists():
        return "n/a"
    text = path.read_text(encoding="utf-8", errors="ignore")
    match = re.search(r"([0-9]+(?:\.[0-9]+)?)%", text)
    return f"{match.group(1)}%" if match else "n/a"


def parse_junit(path: Path | None) -> dict:
    if not path or not path.exists():
        return {
            "cases_total": 0,
            "passed": 0,
            "failed": 0,
            "skipped": 0,
            "suites_total": 0,
            "failed_tests": [],
        }

    root = ET.parse(path).getroot()
    testcases = root.findall(".//testcase")

    total = len(testcases)
    failed = 0
    skipped = 0
    failed_tests: list[str] = []

    for case in testcases:
      # intentional 2-space continuation for XML traversal readability
        classname = case.attrib.get("classname", "playwright")
        name = case.attrib.get("name", "unnamed")
        if case.find("failure") is not None:
            failed += 1
            failed_tests.append(f"{classname} :: {name}")
        elif case.find("skipped") is not None:
            skipped += 1

    passed = total - failed - skipped
    suites = root.findall(".//testsuite")

    return {
        "cases_total": total,
        "passed": passed,
        "failed": failed,
        "skipped": skipped,
        "suites_total": len(suites),
        "failed_tests": failed_tests[:10],
    }


def status_class(failed: int, total: int) -> str:
    if total == 0:
        return "neutral"
    return "pass" if failed == 0 else "fail"


def bar_width(done: int, total: int) -> float:
    if total <= 0:
        return 0.0
    return round(done / total * 100, 1)


def render_failed_list(items: list[str]) -> str:
    if not items:
        return "<div class='empty'>No failed tests.</div>"
    rows = "".join(f"<li>{html.escape(item)}</li>" for item in items)
    return f"<ul>{rows}</ul>"


def metric_card(title: str, subtitle: str, data: dict, extra: str = "") -> str:
    cls = status_class(data["failed"], data["cases_total"])
    return f"""
    <article class="card">
      <div class="eyebrow">{html.escape(title)}</div>
      <div class="metric">{data['cases_total']}</div>
      <div class="subtitle">{html.escape(subtitle)}</div>
      <div class="pill {cls}">{'PASS' if cls == 'pass' else 'CHECK' if cls == 'neutral' else 'FAIL'}</div>
      <div class="stats">
        <span>passed: <b>{data['passed']}</b></span>
        <span>failed: <b>{data['failed']}</b></span>
        <span>skipped: <b>{data['skipped']}</b></span>
      </div>
      {extra}
    </article>
    """


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("--unit-json")
    parser.add_argument("--unit-coverage")
    parser.add_argument("--integration-json")
    parser.add_argument("--e2e-junit")
    parser.add_argument("--output", required=True)
    parser.add_argument("--title", default="Movies Picker Test Report")
    args = parser.parse_args()

    unit = parse_go_json(Path(args.unit_json)) if args.unit_json else parse_go_json(None)
    integration = parse_go_json(Path(args.integration_json)) if args.integration_json else parse_go_json(None)
    e2e = parse_junit(Path(args.e2e_junit)) if args.e2e_junit else parse_junit(None)
    unit_coverage = parse_coverage(Path(args.unit_coverage)) if args.unit_coverage else "n/a"

    total_cases = unit["cases_total"] + integration["cases_total"] + e2e["cases_total"]
    total_failed = unit["failed"] + integration["failed"] + e2e["failed"]

    output = Path(args.output)
    output.parent.mkdir(parents=True, exist_ok=True)

    html_doc = f"""<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>{html.escape(args.title)}</title>
  <style>
    :root {{
      --bg: #f3efe7;
      --panel: #fffaf2;
      --ink: #1f2937;
      --muted: #6b7280;
      --line: #ddd4c5;
      --navy: #1d3557;
      --amber: #b45309;
      --teal: #2a9d8f;
      --green: #15803d;
      --red: #b91c1c;
      --shadow: 0 18px 40px rgba(17,24,39,.08);
    }}
    * {{ box-sizing: border-box; }}
    body {{
      margin: 0;
      padding: 28px;
      color: var(--ink);
      font-family: Georgia, "Times New Roman", serif;
      background:
        radial-gradient(circle at top left, rgba(42,157,143,.10), transparent 28%),
        radial-gradient(circle at right bottom, rgba(180,83,9,.10), transparent 32%),
        var(--bg);
    }}
    .wrap {{ max-width: 1240px; margin: 0 auto; }}
    .hero {{
      padding: 30px 34px;
      border-radius: 28px;
      color: white;
      background: linear-gradient(135deg, rgba(29,53,87,.98), rgba(180,83,9,.92));
      box-shadow: var(--shadow);
    }}
    .hero h1 {{ margin: 0 0 8px; font-size: 38px; line-height: 1.1; }}
    .hero p {{ margin: 0; max-width: 780px; font-size: 17px; opacity: .92; }}
    .meta {{
      display: flex;
      flex-wrap: wrap;
      gap: 10px;
      margin-top: 18px;
    }}
    .meta span {{
      padding: 8px 12px;
      border-radius: 999px;
      background: rgba(255,255,255,.14);
      border: 1px solid rgba(255,255,255,.18);
      font-size: 14px;
    }}
    .grid {{
      margin-top: 20px;
      display: grid;
      grid-template-columns: repeat(12, 1fr);
      gap: 18px;
    }}
    .card, .wide {{
      background: var(--panel);
      border: 1px solid var(--line);
      border-radius: 22px;
      padding: 22px;
      box-shadow: var(--shadow);
    }}
    .card {{ grid-column: span 4; }}
    .wide {{ grid-column: span 12; }}
    .eyebrow {{
      color: var(--muted);
      font-size: 13px;
      letter-spacing: .08em;
      text-transform: uppercase;
      margin-bottom: 10px;
    }}
    .metric {{ font-size: 42px; font-weight: 700; line-height: 1; margin: 0 0 8px; }}
    .subtitle {{ color: var(--muted); font-size: 15px; }}
    .pill {{
      display: inline-block;
      margin-top: 14px;
      padding: 6px 10px;
      border-radius: 999px;
      font-size: 13px;
      font-weight: 700;
    }}
    .pass {{ background: rgba(21,128,61,.12); color: var(--green); }}
    .fail {{ background: rgba(185,28,28,.12); color: var(--red); }}
    .neutral {{ background: rgba(107,114,128,.12); color: var(--muted); }}
    .stats {{
      margin-top: 16px;
      display: flex;
      gap: 14px;
      flex-wrap: wrap;
      color: var(--muted);
      font-size: 14px;
    }}
    .layout {{
      display: grid;
      grid-template-columns: 1.1fr .9fr;
      gap: 18px;
    }}
    .bar {{
      margin-top: 12px;
      height: 16px;
      background: #ece5d8;
      border-radius: 999px;
      overflow: hidden;
    }}
    .bar > span {{
      display: block;
      height: 100%;
      border-radius: 999px;
      background: linear-gradient(90deg, var(--teal), var(--amber));
    }}
    table {{
      width: 100%;
      border-collapse: collapse;
      margin-top: 10px;
      font-size: 15px;
    }}
    th, td {{
      padding: 12px 8px;
      border-bottom: 1px solid var(--line);
      text-align: left;
      vertical-align: top;
    }}
    th {{
      color: var(--muted);
      text-transform: uppercase;
      font-size: 12px;
      letter-spacing: .06em;
    }}
    ul {{ margin: 10px 0 0 18px; }}
    li {{ margin: 8px 0; }}
    .empty {{
      margin-top: 12px;
      color: var(--muted);
      font-style: italic;
    }}
    code {{
      background: #f0e8db;
      padding: 2px 6px;
      border-radius: 6px;
      font-family: Consolas, monospace;
      font-size: 13px;
    }}
    @media (max-width: 920px) {{
      body {{ padding: 16px; }}
      .card {{ grid-column: span 12; }}
      .layout {{ grid-template-columns: 1fr; }}
      .hero h1 {{ font-size: 31px; }}
    }}
  </style>
</head>
<body>
  <div class="wrap">
    <section class="hero">
      <h1>{html.escape(args.title)}</h1>
      <p>Unified CI dashboard for unit, integration, and end-to-end testing. This artifact is meant for quick visual inspection in GitHub Actions and for sharing with reviewers or teachers.</p>
      <div class="meta">
        <span>Total cases: {total_cases}</span>
        <span>Total failed: {total_failed}</span>
        <span>Unit coverage: {unit_coverage}</span>
        <span>Layers included: Unit / Integration / E2E</span>
      </div>
    </section>

    <section class="grid">
      {metric_card("Unit Tests", "Go unit tests for internal/pkg", unit, f"<div class='stats'><span>packages passed: <b>{unit['packages_passed']}</b></span><span>packages failed: <b>{unit['packages_failed']}</b></span><span>coverage: <b>{unit_coverage}</b></span></div>")}
      {metric_card("Integration Tests", "API flows over real app modules", integration, f"<div class='stats'><span>packages passed: <b>{integration['packages_passed']}</b></span><span>packages failed: <b>{integration['packages_failed']}</b></span></div>")}
      {metric_card("System / E2E", "Playwright business scenarios", e2e, f"<div class='stats'><span>suites: <b>{e2e['suites_total']}</b></span></div>")}

      <section class="wide layout">
        <article>
          <div class="eyebrow">Quality Snapshot</div>
          <table>
            <tr><th>Layer</th><th>Passed</th><th>Failed</th><th>Skipped</th><th>Progress</th></tr>
            <tr><td>Unit</td><td>{unit['passed']}</td><td>{unit['failed']}</td><td>{unit['skipped']}</td><td><div class="bar"><span style="width:{bar_width(unit['passed'], unit['cases_total'])}%"></span></div></td></tr>
            <tr><td>Integration</td><td>{integration['passed']}</td><td>{integration['failed']}</td><td>{integration['skipped']}</td><td><div class="bar"><span style="width:{bar_width(integration['passed'], integration['cases_total'])}%"></span></div></td></tr>
            <tr><td>E2E</td><td>{e2e['passed']}</td><td>{e2e['failed']}</td><td>{e2e['skipped']}</td><td><div class="bar"><span style="width:{bar_width(e2e['passed'], e2e['cases_total'])}%"></span></div></td></tr>
          </table>
        </article>
        <article>
          <div class="eyebrow">CI Notes</div>
          <ul>
            <li>Unit layer is validated with <code>go test -json</code> and coverage extraction.</li>
            <li>Integration layer validates full backend flow with Dockerized PostgreSQL and mock services.</li>
            <li>E2E layer uses Playwright in Chromium and WebKit and exports JUnit + HTML artifacts.</li>
            <li>If a section shows zero tests, that layer was not included in the current workflow run.</li>
          </ul>
        </article>
      </section>

      <section class="wide layout">
        <article>
          <div class="eyebrow">Failed Unit Tests</div>
          {render_failed_list(unit['failed_tests'])}
        </article>
        <article>
          <div class="eyebrow">Failed Integration Tests</div>
          {render_failed_list(integration['failed_tests'])}
        </article>
      </section>

      <section class="wide">
        <div class="eyebrow">Failed E2E Tests</div>
        {render_failed_list(e2e['failed_tests'])}
      </section>
    </section>
  </div>
</body>
</html>
"""

    output.write_text(html_doc, encoding="utf-8")
    print(f"Generated CI report: {output}")


if __name__ == "__main__":
    main()
