package docker

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/voormedia/kd/pkg/config"
	"github.com/voormedia/kd/pkg/util"
)

func Build(log *util.Logger, app *config.ResolvedApp) error {
	dockerfile := filepath.Join(app.Path, "Dockerfile")

	args := []string{
		"buildx", "build",
	}

	if _, ok := os.LookupEnv("SSH_AUTH_SOCK"); ok {
		log.Note("Enabled SSH agent key forwarding")
		args = append(args, "--ssh", "default")
	}

	return run(append(args,
		"--file", dockerfile,
		"--tag", app.Repository(),
		"--platform", app.Platform,
		app.Root,
	)...)
}

func Push(app *config.ResolvedApp) error {
	return run(
		"push", app.Repository(),
	)
}

func run(args ...string) error {
	cmd := exec.Command("docker", args...)

	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return cmd.Run()
}
