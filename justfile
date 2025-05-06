fmt:
    gci write . --skip-generated
    gofumpt -w .

lint *args:
    golangci-lint run {{ args }}

test:
    go run gotest.tools/gotestsum@latest --format testname ./...
    cd omitzero && go run gotest.tools/gotestsum@latest --format testname ./...
