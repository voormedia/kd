package build

import (
	"strings"

	"github.com/voormedia/kd/pkg/config"
	"github.com/voormedia/kd/pkg/internal/docker"
	"github.com/voormedia/kd/pkg/util"
)

func Run(log *util.Logger, app *config.ResolvedApp, buildCacheTag string, secrets []string, producer string) error {
	if app.SkipBuild {
		log.Fatal("Build is skipped for", app.Name)
	}

	if app.PreBuild != "" {
		if strings.Contains(app.PreBuild, "/.ssh") {
			log.Warn("Pre-build command in 'kdeploy.conf' contains reference to '.ssh'.")
			log.Warn("Please use SSH key forwarding: https://github.com/voormedia/kd#ssh-forwarding")
		}

		err := util.Run(log, "sh", "-c", app.PreBuild)
		if err != nil {
			log.Fatal("Pre-build command failed:", err)
		}
	}

	log.Note("Building", app.Name)

	if err := docker.Build(log, app, buildCacheTag, secrets, producer); err != nil {
		log.Fatal(err)
	}

	log.Note("Pushed to", app.Repository())

	if app.PostBuild != "" {
		err := util.Run(log, "sh", "-c", app.PostBuild)
		if err != nil {
			log.Fatal("Post-build command failed:", err)
		}
	}

	log.Success("Successfully built", app.Name+":"+app.Tag)
	return nil
}
