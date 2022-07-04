package build

import (
	"os/exec"

	"github.com/voormedia/kd/pkg/config"
	"github.com/voormedia/kd/pkg/internal/docker"
	"github.com/voormedia/kd/pkg/util"
)

func Run(log *util.Logger, app *config.ResolvedApp) error {
	if app.PreBuild != "" {
		cmd := exec.Command("sh", "-c", app.PreBuild)
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}
	}

	log.Note("Building", app.Name)
	if err := docker.Build(app); err != nil {
		log.Fatal(err)
	}

	log.Note("Pushing", app.Name+":"+app.Tag)
	if err := docker.Push(app); err != nil {
		log.Fatal(err)
	}

	if app.PostBuild != "" {
		cmd := exec.Command("sh", "-c", app.PostBuild)
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}
	}

	log.Success("Successfully built", app.Repository())
	return nil
}
