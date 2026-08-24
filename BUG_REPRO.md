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
?   	warrantyservice/internal/model	[no test files]
ok  	warrantyservice/cmd/warranty-server	0.042s
ok  	warrantyservice/internal/adapters	0.019s
ok  	warrantyservice/internal/config	0.011s
ok  	warrantyservice/internal/httpapi	0.058s
ok  	warrantyservice/internal/integration	0.064s
ok  	warrantyservice/internal/persistence	0.142s
ok  	warrantyservice/internal/support	0.015s
--- FAIL: TestWarrantyMissingRecord (0.01s)
    missing_record_test.go:16: expected missing status, got {Status:valid Record:<nil> Advice:<nil> ServicePoints:[] Message:查询成功}
FAIL
FAIL	warrantyservice/internal/warranty	0.061s
FAIL
```

## Architecture reproduction

### linux/amd64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/warranty-server): exit `0`
### linux/arm64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/warranty-server): exit `0`
