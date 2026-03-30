param(
  [string]$ProjectRoot = (Resolve-Path (Join-Path $PSScriptRoot "..")).Path
)

$ErrorActionPreference = "Stop"

Set-Location $ProjectRoot

$generatedDir = Join-Path $ProjectRoot "report\generated"
New-Item -ItemType Directory -Force -Path $generatedDir | Out-Null

function Get-CoverageTotal {
  param([string]$CoverageFile)

  if (-not (Test-Path $CoverageFile)) {
    return $null
  }

  $line = & go tool cover "-func=$CoverageFile" | Select-Object -Last 1
  if ($line -match '([0-9]+(\.[0-9]+)?)%') {
    return [double]$matches[1]
  }
  return $null
}

function Count-RgMatches {
  param(
    [string]$Pattern,
    [string[]]$Path
  )

  $result = & rg -n $Pattern @Path 2>$null
  if (-not $result) {
    return 0
  }
  return @($result).Count
}

$unitCoverageFile = Join-Path $ProjectRoot "coverage_unit_all"
$unitCoverage = Get-CoverageTotal -CoverageFile $unitCoverageFile
$unitTopLevelTests = Count-RgMatches -Pattern '^func Test' -Path @('internal', 'pkg')
$unitSubtests = Count-RgMatches -Pattern 't\.Run\(' -Path @('internal', 'pkg')
$unitCases = $unitTopLevelTests + $unitSubtests

$integrationScenarioCount = Count-RgMatches -Pattern '^func TestINT' -Path @('test\integration')
$e2eScenarioCount = Count-RgMatches -Pattern "test\('TC-" -Path @('e2e\tests')

$workflowUnitExists = Test-Path (Join-Path $ProjectRoot ".github\workflows\test.yml")
$workflowSystemExists = Test-Path (Join-Path $ProjectRoot ".github\workflows\system-tests.yml")
$playwrightReportExists = Test-Path (Join-Path $ProjectRoot "e2e\playwright-report\index.html")

$unitStatus = if ($unitCoverage -ge 80) { "PASS" } else { "CHECK" }
$integrationStatus = if ($integrationScenarioCount -ge 10 -and $workflowUnitExists) { "PASS" } else { "CHECK" }
$e2eStatus = if ($e2eScenarioCount -ge 10 -and $workflowSystemExists) { "PASS" } else { "CHECK" }

$generatedAt = Get-Date -Format "yyyy-MM-dd HH:mm:ss"
$reportPath = Join-Path $generatedDir "test-summary.html"

$unitCoverageDisplay = if ($null -ne $unitCoverage) { "{0:N1}%" -f $unitCoverage } else { "n/a" }
$playwrightDisplay = if ($playwrightReportExists) { "Available" } else { "Not generated yet" }
$unitWorkflowDisplay = if ($workflowUnitExists) { "Configured" } else { "Missing" }
$systemWorkflowDisplay = if ($workflowSystemExists) { "Configured" } else { "Missing" }
$coverageWidth = if ($null -ne $unitCoverage) { [math]::Min($unitCoverage, 100) } else { 0 }
$coverageDash = if ($null -ne $unitCoverage) { [math]::Round(5.152 * $unitCoverage, 1) } else { 0 }
$statusClass = if ($unitStatus -eq "PASS") { "pass" } else { "check" }
$integrationStatusClass = if ($integrationStatus -eq "PASS") { "pass" } else { "check" }
$e2eStatusClass = if ($e2eStatus -eq "PASS") { "pass" } else { "check" }

$html = @"
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>Movies Picker Test Summary</title>
  <style>
    :root {
      --bg: #f4efe6;
      --panel: #fffaf2;
      --ink: #1f2937;
      --muted: #6b7280;
      --line: #d6cfc2;
      --accent: #b45309;
      --good: #15803d;
      --warn: #b45309;
      --navy: #1d3557;
      --teal: #2a9d8f;
    }
    * { box-sizing: border-box; }
    body {
      margin: 0;
      padding: 32px;
      font-family: Georgia, "Times New Roman", serif;
      color: var(--ink);
      background:
        radial-gradient(circle at top left, rgba(180,83,9,.10), transparent 32%),
        radial-gradient(circle at bottom right, rgba(42,157,143,.10), transparent 28%),
        var(--bg);
    }
    .wrap { max-width: 1180px; margin: 0 auto; }
    .hero {
      background: linear-gradient(135deg, rgba(29,53,87,.95), rgba(180,83,9,.92));
      color: #fff;
      border-radius: 24px;
      padding: 28px 32px;
      box-shadow: 0 24px 60px rgba(29,53,87,.18);
    }
    .hero h1 {
      margin: 0 0 8px;
      font-size: 38px;
      line-height: 1.1;
    }
    .hero p {
      margin: 0;
      max-width: 760px;
      font-size: 17px;
      opacity: .92;
    }
    .meta {
      display: flex;
      gap: 12px;
      flex-wrap: wrap;
      margin-top: 18px;
    }
    .pill {
      padding: 8px 12px;
      border-radius: 999px;
      background: rgba(255,255,255,.14);
      border: 1px solid rgba(255,255,255,.18);
      font-size: 14px;
    }
    .grid {
      display: grid;
      grid-template-columns: repeat(12, 1fr);
      gap: 18px;
      margin-top: 22px;
    }
    .card {
      background: var(--panel);
      border: 1px solid var(--line);
      border-radius: 22px;
      padding: 22px;
      box-shadow: 0 10px 30px rgba(31,41,55,.06);
    }
    .span-4 { grid-column: span 4; }
    .span-8 { grid-column: span 8; }
    .span-12 { grid-column: span 12; }
    .eyebrow {
      color: var(--muted);
      font-size: 13px;
      text-transform: uppercase;
      letter-spacing: .08em;
      margin-bottom: 10px;
    }
    .metric {
      font-size: 40px;
      font-weight: 700;
      line-height: 1;
      margin: 6px 0 8px;
    }
    .subtle {
      color: var(--muted);
      font-size: 15px;
    }
    .status {
      display: inline-block;
      margin-top: 14px;
      padding: 6px 10px;
      border-radius: 999px;
      font-size: 13px;
      font-weight: 700;
    }
    .pass { background: rgba(21,128,61,.12); color: var(--good); }
    .check { background: rgba(180,83,9,.14); color: var(--warn); }
    .bar {
      height: 16px;
      background: #ece6da;
      border-radius: 999px;
      overflow: hidden;
      margin-top: 14px;
    }
    .fill {
      height: 100%;
      border-radius: 999px;
      background: linear-gradient(90deg, var(--teal), var(--accent));
    }
    table {
      width: 100%;
      border-collapse: collapse;
      margin-top: 10px;
      font-size: 15px;
    }
    th, td {
      border-bottom: 1px solid var(--line);
      padding: 12px 8px;
      text-align: left;
      vertical-align: top;
    }
    th {
      color: var(--muted);
      font-size: 13px;
      text-transform: uppercase;
      letter-spacing: .05em;
    }
    .chart-row {
      display: grid;
      grid-template-columns: 280px 1fr;
      gap: 18px;
      align-items: center;
    }
    .legend-item {
      display: flex;
      align-items: center;
      gap: 10px;
      margin: 8px 0;
      font-size: 15px;
    }
    .dot {
      width: 12px;
      height: 12px;
      border-radius: 999px;
      display: inline-block;
    }
    .notes ul {
      margin: 12px 0 0 18px;
      padding: 0;
    }
    .notes li {
      margin: 8px 0;
    }
    code {
      background: #f1ebdf;
      padding: 2px 6px;
      border-radius: 6px;
      font-family: Consolas, monospace;
      font-size: 13px;
    }
    @media (max-width: 900px) {
      body { padding: 16px; }
      .span-4, .span-8, .span-12 { grid-column: span 12; }
      .chart-row { grid-template-columns: 1fr; }
      .hero h1 { font-size: 30px; }
    }
  </style>
</head>
<body>
  <div class="wrap">
    <section class="hero">
      <h1>Movies Picker Testing Dashboard</h1>
      <p>Unified local report for unit, integration, and system testing. Use it for screenshots, appendix materials, and quick review before the defense.</p>
      <div class="meta">
        <div class="pill">Generated: $generatedAt</div>
        <div class="pill">Unit coverage: $unitCoverageDisplay</div>
        <div class="pill">Integration scenarios: $integrationScenarioCount</div>
        <div class="pill">E2E scenarios: $e2eScenarioCount</div>
      </div>
    </section>

    <section class="grid">
      <article class="card span-4">
        <div class="eyebrow">Unit Testing</div>
        <div class="metric">$unitCoverageDisplay</div>
        <div class="subtle">Coverage of packages <code>internal</code> and <code>pkg</code>.</div>
        <div class="bar"><div class="fill" style="width: $coverageWidth%"></div></div>
        <div class="status $statusClass">$unitStatus</div>
      </article>

      <article class="card span-4">
        <div class="eyebrow">Integration Testing</div>
        <div class="metric">$integrationScenarioCount</div>
        <div class="subtle">Scenarios <code>TestINT..</code> in <code>test/integration</code>.</div>
        <div class="status $integrationStatusClass">$integrationStatus</div>
      </article>

      <article class="card span-4">
        <div class="eyebrow">System / E2E</div>
        <div class="metric">$e2eScenarioCount</div>
        <div class="subtle">Playwright scenarios for user-facing flows.</div>
        <div class="status $e2eStatusClass">$e2eStatus</div>
      </article>

      <article class="card span-8">
        <div class="eyebrow">Coverage And Scenario Balance</div>
        <div class="chart-row">
          <svg viewBox="0 0 220 220" width="220" height="220" aria-label="Coverage chart">
            <circle cx="110" cy="110" r="82" fill="none" stroke="#e7dfd1" stroke-width="24"></circle>
            <circle cx="110" cy="110" r="82" fill="none" stroke="#2a9d8f" stroke-width="24"
              stroke-linecap="round"
              stroke-dasharray="$coverageDash 999"
              transform="rotate(-90 110 110)"></circle>
            <text x="110" y="104" text-anchor="middle" font-size="18" fill="#6b7280">Unit</text>
            <text x="110" y="128" text-anchor="middle" font-size="28" font-weight="700" fill="#1f2937">$unitCoverageDisplay</text>
          </svg>
          <div>
            <div class="legend-item"><span class="dot" style="background:#2a9d8f"></span>Unit coverage over source packages</div>
            <div class="legend-item"><span class="dot" style="background:#1d3557"></span>Integration scenarios cover API, service, DB, and mocks</div>
            <div class="legend-item"><span class="dot" style="background:#b45309"></span>E2E scenarios cover full user journeys in browser</div>
            <table>
              <tr><th>Metric</th><th>Value</th></tr>
              <tr><td>Top-level unit tests</td><td>$unitTopLevelTests</td></tr>
              <tr><td>Unit subtests</td><td>$unitSubtests</td></tr>
              <tr><td>Total unit cases</td><td>$unitCases</td></tr>
              <tr><td>Integration cases</td><td>$integrationScenarioCount</td></tr>
              <tr><td>E2E cases</td><td>$e2eScenarioCount</td></tr>
            </table>
          </div>
        </div>
      </article>

      <article class="card span-4">
        <div class="eyebrow">CI Status</div>
        <table>
          <tr><th>Workflow</th><th>Status</th></tr>
          <tr><td>Unit + Integration</td><td>$unitWorkflowDisplay</td></tr>
          <tr><td>System E2E</td><td>$systemWorkflowDisplay</td></tr>
          <tr><td>Playwright HTML report</td><td>$playwrightDisplay</td></tr>
        </table>
      </article>

      <article class="card span-12 notes">
        <div class="eyebrow">Local Commands</div>
        <table>
          <tr><th>Layer</th><th>Command</th><th>Purpose</th></tr>
          <tr><td>Unit</td><td><code>go test ./internal/... ./pkg/... -coverprofile=coverage_unit_all</code></td><td>Run unit tests and refresh coverage profile.</td></tr>
          <tr><td>Integration</td><td><code>go test ./test/integration/... -v -timeout 3m</code></td><td>Run API integration flows against local stack.</td></tr>
          <tr><td>E2E</td><td><code>cd e2e; npm test</code></td><td>Run Playwright system tests.</td></tr>
          <tr><td>Playwright report</td><td><code>cd e2e; npm run report</code></td><td>Open the detailed Playwright browser report.</td></tr>
          <tr><td>Summary</td><td><code>powershell -ExecutionPolicy Bypass -File scripts/generate-test-report.ps1</code></td><td>Generate this dashboard.</td></tr>
        </table>
      </article>

      <article class="card span-12 notes">
        <div class="eyebrow">What To Say At The Defense</div>
        <ul>
          <li>Unit tests are written in Go with <code>testing</code>, <code>testify</code>, <code>sqlmock</code>, and <code>httptest</code>; current coverage is $unitCoverageDisplay.</li>
          <li>Integration tests cover $integrationScenarioCount complete API and backend scenarios and use a mock server to isolate Kinopoisk and YouTube.</li>
          <li>System tests cover $e2eScenarioCount major user flows through Playwright in Chromium and WebKit.</li>
          <li>CI is split into two workflows: <code>.github/workflows/test.yml</code> for unit/integration and <code>.github/workflows/system-tests.yml</code> for E2E.</li>
          <li>To extend the test set, add a scenario to the test plan, then add the test itself, then refresh CI and reports.</li>
        </ul>
      </article>
    </section>
  </div>
</body>
</html>
"@

$utf8NoBom = New-Object System.Text.UTF8Encoding($false)
[System.IO.File]::WriteAllText($reportPath, $html, $utf8NoBom)
Write-Output "Generated report: $reportPath"
