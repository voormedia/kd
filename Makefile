all:
	go build -i -v

release:
	go build -i -ldflags "-s -w" -v
	upx -9 -q kd > /dev/null

test:
	go test -v ./pkg/...

.PHONY: all release test
