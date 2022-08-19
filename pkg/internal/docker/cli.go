package docker

import (
	"net"
	"os"
	"os/exec"
	"path/filepath"

	"golang.org/x/crypto/ssh/agent"
	"github.com/voormedia/kd/pkg/config"
	"github.com/voormedia/kd/pkg/util"
)

func Build(log *util.Logger, app *config.ResolvedApp) error {
	dockerfile := filepath.Join(app.Path, "Dockerfile")

	args := []string{
		"buildx", "build",
	}

	if sock, ok := os.LookupEnv("SSH_AUTH_SOCK"); ok {
		conn, err := net.Dial("unix", sock)

		if err != nil {
			return err
		}

		signers, err := agent.NewClient(conn).Signers()
		if err != nil {
			return err
		}

		if len(signers) == 0 {
			log.Warn("Enabled SSH agent key forwarding, but no SSH keys are exposed")
		} else {
			log.Note("Enabled SSH agent key forwarding")
		}

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
