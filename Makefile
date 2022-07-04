all:
	go build -v

test:
	go test -v ./pkg/...

.PHONY: all test
