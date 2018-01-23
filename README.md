# Build and deploy apps to Kubernetes engine

## Installing

### Option 1 – prebuilt
1. Ensure you have `~/.bin` directory or similar that is in your `$PATH`
1. `curl -L https://github.com/voormedia/kd/releases/download/v1.0.0/darwin_amd64_kd -o ~/.bin/kd && chmod +x ~/.bin/kd`

### Option 2 – from source
1. Make sure you have a working `go` installation
2. `go install github.com/voormedia/kd`
