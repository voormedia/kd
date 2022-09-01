package build

import (
	"os"
	"os/exec"
	"strings"

	"github.com/voormedia/kd/pkg/config"
	"github.com/voormedia/kd/pkg/internal/docker"
	"github.com/voormedia/kd/pkg/util"
)

func Run(log *util.Logger, app *config.ResolvedApp) error {
	if app.PreBuild != "" {
		if strings.Contains(app.PreBuild, "/.ssh") {
			log.Warn("Pre-build command in 'kdeploy.conf' contains reference to '.ssh'.")
			log.Warn("Please use SSH key forwarding: https://github.com/voormedia/kd#ssh-forwarding")
		}

		cmd := exec.Command("sh", "-c", app.PreBuild)
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout

		if err := cmd.Run(); err != nil {
			log.Fatal("Pre-build command failed:", err)
		}
	}

	log.Note("Building", app.Name)
	if err := docker.Build(log, app); err != nil {
		log.Fatal(err)
	}

	log.Note("Pushing", app.Name+":"+app.Tag)
	if err := docker.Push(app); err != nil {
		log.Fatal(err)
	}

	if app.PostBuild != "" {
		cmd := exec.Command("sh", "-c", app.PostBuild)
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout

		if err := cmd.Run(); err != nil {
			log.Fatal("Post-build command failed:", err)
		}
	}

	log.Success("Successfully built", app.Repository())
	return nil
}
