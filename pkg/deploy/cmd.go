package deploy

import (
	"github.com/voormedia/kd/pkg/config"
	"github.com/voormedia/kd/pkg/internal/docker"
	"github.com/voormedia/kd/pkg/internal/kubectl"
	"github.com/voormedia/kd/pkg/internal/kustomize"
	"github.com/voormedia/kd/pkg/util"
)

func Run(log *util.Logger, app *config.ResolvedApp, target *config.ResolvedTarget) error {
	log.Note("Retrieving image", app.Name+":"+app.Tag)
	img, err := docker.GetImage(app.Repository())
	if err != nil {
		return err
	}

	res, err := kustomize.GetResources(app, target, img.Descriptor.Digest.String())
	if err != nil {
		return err
	}

	vrs, err := kubectl.Version()
	if err != nil {
		return err
	}

	log.Note("Applying configuration with kubectl", vrs)
	err = kubectl.ApplyFromStdin(target, res)
	if err != nil {
		return err
	}

	log.Note("Tagging image", app.Name+":"+target.Name)
	err = docker.TagImage(img, app.RepositoryWithTag(target.Name))
	if err != nil {
		return err
	}

	log.Success("Successfully deployed", app.Repository(), "to", target.Name)
	return nil
}
