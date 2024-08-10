
test:
	rm -rf testdata/TestIntegration/
	go build -v -cover -ldflags '-s -w -X main.version=test' -o ace .
	ACE_TESTBIN=./ace go test -v .
.PHONY: test

coverage.txt: export GOCOVERDIR:=${shell mktemp -d}
coverage.txt: test
	go tool covdata textfmt -i=${GOCOVERDIR} -o $@

coverage.html: coverage.txt
	go tool cover -html=$< -o $@
