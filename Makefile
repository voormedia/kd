all:
	go build -v

test:
	go test -v ./pkg/...

cross:
	gox -ldflags "-s -w -X github.com/voormedia/kd/cmd.version=${version}" -arch arm64 -os darwin -output "dist/{{.OS}}_{{.Arch}}_{{.Dir}}"
	gox -ldflags "-s -w -X github.com/voormedia/kd/cmd.version=${version}" -arch amd64 -os darwin -output "dist/{{.OS}}_{{.Arch}}_{{.Dir}}"
	gox -ldflags "-s -w -X github.com/voormedia/kd/cmd.version=${version}" -arch amd64 -os linux -output "dist/{{.OS}}_{{.Arch}}_{{.Dir}}"

.PHONY: all test
