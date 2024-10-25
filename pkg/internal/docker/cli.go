package docker

import (
	"net"
	"os"
	"path/filepath"

	"github.com/voormedia/kd/pkg/config"
	"github.com/voormedia/kd/pkg/util"
	"golang.org/x/crypto/ssh/agent"
)

func Build(log *util.Logger, app *config.ResolvedApp) error {
	dockerfile := filepath.Join(app.Path, "Dockerfile")

	cmd := []string{
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

		cmd = append(cmd, "--ssh", "default")
	}

	return util.Run(log,
		"docker", append(cmd,
			"--file", dockerfile,
			"--tag", app.Repository(),
			"--platform", app.Platform,
			app.Root,
		)...)
}

func Push(log *util.Logger, app *config.ResolvedApp) error {
	return util.Run(log,
		"docker",
		"push", app.Repository(),
	)
}
