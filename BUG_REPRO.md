# BUG_REPRO

The following failures were observed while validating the initial project state.
Each section records what failed, how to reproduce it, and the complete command output.
They are preserved intentionally; only failing build gates are omitted from the generated Dockerfile.

## Failure 1: Go test (.)

- Observed problem: `Go test (.)` failed in the initial project state.
- Working directory: `.`
- Command: `cd /app && GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off go test -count=1 ./...`
- Exit status: `1`

```text
?   	telemetry.local/drone/cmd/telemetryd	[no test files]
?   	telemetry.local/drone/internal/crypto	[no test files]
?   	telemetry.local/drone/internal/domain	[no test files]
?   	telemetry.local/drone/internal/service	[no test files]
ok  	telemetry.local/drone/internal/api	0.010s
ok  	telemetry.local/drone/internal/flow001	0.010s
ok  	telemetry.local/drone/internal/flow002	0.012s
ok  	telemetry.local/drone/internal/flow003	0.009s
--- FAIL: Test2200BusinessRegression (0.00s)
    regression_test.go:24: first record title="Expired Flight Plan"
FAIL
FAIL	telemetry.local/drone/internal/flow012	0.001s
ok  	telemetry.local/drone/internal/importer	0.002s
ok  	telemetry.local/drone/internal/report	0.001s
ok  	telemetry.local/drone/internal/review	0.001s
ok  	telemetry.local/drone/internal/store	0.002s
FAIL
```

## Architecture reproduction

### linux/amd64
- Go toolchain version: exit `0`
- Node.js version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/telemetryd): exit `0`
- Frontend build (web): exit `0`
### linux/arm64
- Go toolchain version: exit `0`
- Node.js version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/telemetryd): exit `0`
- Frontend build (web): exit `0`
