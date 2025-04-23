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

func Build(log *util.Logger, app *config.ResolvedApp, buildCacheTag string, secrets []string, producer string) error {
	dockerfile := filepath.Join(app.Path, "Dockerfile")

	cmd := []string{
		"buildx", "build",
	}

	buildCacheTagParts := []string{}

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
			if buildCacheTag == "" {
				buildCacheTagParts = append(buildCacheTagParts, "ssh")
			}
		}

		cmd = append(cmd, "--ssh", "default")
	}

	for _, secret := range secrets {
		cmd = append(cmd, "--secret", secret)
		parts := strings.Split(secret, "=")
		if len(parts) > 1 {
			name := parts[1]
			buildCacheTagParts = append(buildCacheTagParts, strings.ToLower(strings.ReplaceAll(name, "_", "-")))
		}
	}

	buildCacheFallbackTag := ""

	if buildCacheTag == "" {
		currentBranch, err := util.GetCurrentBranch(log, app.Path)
		if err != nil {
			log.Warn("Could not determine current branch:", err)
			currentBranch = "unknown"
		}

		buildCacheTag = strings.Join(append([]string{currentBranch}, buildCacheTagParts...), "-")

		if currentBranch != "main" {
			buildCacheFallbackTag = strings.Join(append([]string{"main"}, buildCacheTagParts...), "-")
		}
	}

	if supportsCacheExport(log) {
		targetBuildCache := app.RepositoryBuildCache(buildCacheTag)
		cmd = append(cmd,
			"--provenance=false",
			"--cache-to", "type=registry,ref="+targetBuildCache+",mode=max",
			"--cache-from", "type=registry,ref="+targetBuildCache,
		)

		if buildCacheFallbackTag != "" {
			targetFallbackBuildCache := app.RepositoryBuildCache(buildCacheFallbackTag)
			cmd = append(cmd,
				"--cache-from", "type=registry,ref="+targetFallbackBuildCache,
			)
		}
	} else {
		log.Warn("Builder does not support remote cache, using local cache only")
	}

	return util.Run(log,
		"docker", append(cmd,
			"--output=type=image,name="+app.Repository()+",push=true",
			"--file", dockerfile,
			"--tag", app.Repository(),
			"--platform", app.Platform,
			"--label", "producer="+producer,
			app.Root,
		)...)
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
