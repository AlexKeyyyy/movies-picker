# Local Testing Guide

## 1. Unit tests

Run:

```powershell
go test ./internal/... ./pkg/... -coverprofile=coverage_unit_all
go tool cover "-func=coverage_unit_all"
```

What this gives:

- all unit tests for `internal` and `pkg`
- updated coverage profile
- percentage for the report

## 2. Integration tests

Start the stack:

```powershell
docker compose down -v
docker compose up -d postgres mock-server api
```

Wait until services are healthy:

- API: `http://localhost:8080/health`
- Mock server: `http://localhost:9090/health`

Then run:

```powershell
go test ./test/integration/... -v -timeout 3m
```

If you want the same sequence as CI, use:

```powershell
.\run.ps1
```

## 3. System / End-to-End tests

The full app must be running:

```powershell
docker compose up -d postgres mock-server api frontend
```

Then inside `e2e`:

```powershell
cd e2e
npm install
npx playwright install chromium webkit
npm test
```

Useful variants:

```powershell
npm run test:headed
npm run test:debug
npm run report
```

## 4. CI workflows

Available workflows:

- `.github/workflows/test.yml`
  runs unit + integration on `push` and `pull_request`
- `.github/workflows/system-tests.yml`
  runs E2E manually (`workflow_dispatch`) and by schedule (`cron`)

## 5. Beautiful summary report

After running tests, generate the local HTML dashboard:

```powershell
powershell -ExecutionPolicy Bypass -File scripts/generate-test-report.ps1
```

Output:

- `report/generated/test-summary.html`

This file is convenient for screenshots and appendix materials.

## 6. GitHub CI reports

Two report modes are now available in GitHub Actions:

- `CI: Unit & Integration Tests`
  artifact: `ci-unit-integration-report`
- `CI: System Tests`
  artifacts:
  `all-tests-report` and `playwright-report`

The most useful artifact for defense is `all-tests-report`, because it contains one HTML page with unit, integration, and E2E sections together.

## 7. What to explain at the defense

### Unit testing

- which packages are covered
- which tools were used: `testing`, `testify`, `sqlmock`, `httptest`
- which test design techniques were used:
  equivalence classes, boundary values, negative scenarios, table-driven tests
- how coverage is measured

### Integration testing

- which modules interact: API, service, repository, PostgreSQL, mock external APIs
- why mock-server is used
- that there are positive and negative scenarios
- that scenarios are documented in `docs/testing/integration-tests.md`
- that CI automatically runs integration tests

### System / E2E testing

- that Playwright checks business scenarios from the user perspective
- browsers used: Chromium and WebKit
- report format: Playwright HTML report
- CI execution mode: manual and scheduled

## 8. Extension procedure example

Example: new feature `Recommendations`

1. Add scenario to the test plan.
2. Add mock endpoint if external data is needed.
3. Add unit tests for service/handler logic.
4. Add integration flow in `test/integration`.
5. Add Playwright scenario for the UI.
6. Re-run tests and regenerate reports.
