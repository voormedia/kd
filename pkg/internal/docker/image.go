package docker

import (
	"context"
	"os"
	"runtime"

	"github.com/distribution/reference"
	"github.com/docker/cli/cli/config"
	manifesttypes "github.com/docker/cli/cli/manifest/types"
	"github.com/docker/cli/cli/registry/client"
	"github.com/docker/docker/api/types/registry"
	"github.com/voormedia/kd/pkg/util"
)

type ImageManifest = manifesttypes.ImageManifest

func GetImage(log *util.Logger, location string) (ImageManifest, error) {
	ref, err := reference.ParseNormalizedNamed(location)
	if err != nil {
		return ImageManifest{}, err
	}

	ctx := context.Background()
	log.Debug("Retrieving remote manifest", ref)
	return newClient().GetManifest(ctx, ref)
}

func TagImage(log *util.Logger, manifest ImageManifest, location string) error {
	ref, err := reference.ParseNormalizedNamed(location)
	if err != nil {
		return err
	}

	ctx := context.Background()
	log.Debug("Storing remote manifest", ref)
	_, err = newClient().PutManifest(ctx, ref, manifest)
	return err
}

func newClient() client.RegistryClient {
	return client.NewRegistryClient(resolver, userAgent, false)
}

func resolver(ctx context.Context, index *registry.IndexInfo) registry.AuthConfig {
	conf := config.LoadDefaultConfigFile(os.Stderr)
	authConfig, _ := conf.GetAuthConfig(index.Name)
	return registry.AuthConfig(authConfig)
}

const userAgent = "KD (" + runtime.GOOS + ")"
