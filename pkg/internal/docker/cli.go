package docker

import (
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/voormedia/kd/pkg/config"
	"github.com/voormedia/kd/pkg/util"
	"golang.org/x/crypto/ssh/agent"
)

func Build(log *util.Logger, app *config.ResolvedApp, buildCacheTag string) error {
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

	targetBuildCache := app.RepositoryBuildCache(buildCacheTag)

	if supportsCacheExport(log) {
		cmd = append(cmd,
			"--provenance=false",
			"--cache-to", "type=registry,ref="+targetBuildCache+",mode=max",
			"--cache-from", "type=registry,ref="+targetBuildCache,
		)
	} else {
		log.Warn("Builder does not support remote cache, using local cache only")
	}

	if buildCacheTag != "main" {
		cmd = append(cmd,
			"--cache-from", "type=registry,ref="+app.RepositoryBuildCache("main"),
		)
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
	// Swallow stderr to ignore nag about pushing single-platform image.
	return util.RunWithoutStdErr(log,
		"docker",
		"push", app.Repository(),
		"--platform", app.Platform,
	)
}

func supportsCacheExport(log *util.Logger) bool {
	output, err := util.Capture(log, "docker", "buildx", "inspect", "--debug")
	if err != nil {
		return false
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "Cache export:") {
			return strings.Contains(line, "true")
		}
	}

	return false
}
