package deploy

import (
	"strings"

	"github.com/voormedia/kd/pkg/config"
	"github.com/voormedia/kd/pkg/internal/docker"
	"github.com/voormedia/kd/pkg/internal/gcloud"
	"github.com/voormedia/kd/pkg/internal/kubectl"
	"github.com/voormedia/kd/pkg/internal/kustomize"
	"github.com/voormedia/kd/pkg/util"
)

func Run(log *util.Logger, app *config.ResolvedApp, target *config.ResolvedTarget) error {
	log.Note("Retrieving image", app.Name+":"+app.Tag)
	img, err := docker.GetImage(log, app.Repository())
	if err != nil {
		return err
	}

	res, err := kustomize.GetResources(log, app, target, img.Descriptor.Digest.String())
	if err != nil {
		return err
	}

	vrs, err := kubectl.Version(log)
	if err != nil {
		return err
	}

	log.Note("Applying configuration with kubectl", vrs)
	err = kubectl.ApplyFromStdin(log, target, res)
	if err != nil {
		return err
	}

	log.Note("Tagging image", app.Name+":"+target.Name)
	err = docker.TagImage(log, img, app.RepositoryWithTag(target.Name))
	if err != nil {
		return err
	}

	ingresses, err := kubectl.GetGCEIngresses(log, target)
	if err != nil {
		return err
	}

	var names []string

	for _, ingress := range ingresses {
		flushed, err := gcloud.ScheduleCDNCacheFlush(log, target, ingress.Annotations)
		if err != nil {
			return err
		}

		if flushed {
			names = append(names, ingress.Name)
		}
	}

	if len(names) > 0 {
		log.Note("Requesting cache flush for", strings.Join(names, ", "))
	}

	log.Success("Successfully deployed", app.Repository(), "to", target.Name)
	return nil
}
