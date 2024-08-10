test:
	rm -rf testdata/TestIntegration/
	go build -v -cover -ldflags '-s -w -X main.version=test' -o ace .
	ACE_TESTBIN=./ace go test -v .
.PHONY: test
