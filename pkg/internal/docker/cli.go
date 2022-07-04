package docker

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/voormedia/kd/pkg/config"
)

func Build(app *config.ResolvedApp) error {
	dockerfile := filepath.Join(app.Path, "Dockerfile")

	return run(
		"buildx",
		"build", "--ssh", "default",
		"--file", dockerfile,
		"--tag", app.Repository(),
		app.Root,
	)
}

func Push(app *config.ResolvedApp) error {
	return run(
		"push", app.Repository(),
	)
}

func run(args ...string) error {
	cmd := exec.Command("docker", args...)

	cmd.Stdin = bytes.NewReader([]byte{})
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return cmd.Run()
}
