all:
	go build -i -v -o kd

release:
	go build -i -ldflags "-s -w" -v -o kd
	upx -9 -q kd > /dev/null

test:
	go test -v ./pkg/...

.PHONY: all release test
