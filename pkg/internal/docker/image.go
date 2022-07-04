package docker

import (
	"context"
	"os"
	"runtime"

	"github.com/docker/cli/cli/config"
	manifesttypes "github.com/docker/cli/cli/manifest/types"
	"github.com/docker/cli/cli/registry/client"
	"github.com/docker/distribution/reference"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/registry"
)

type ImageManifest = manifesttypes.ImageManifest

func GetImage(location string) (ImageManifest, error) {
	ref, err := reference.ParseNormalizedNamed(location)
	if err != nil {
		return ImageManifest{}, err
	}

	ctx := context.Background()
	return newClient().GetManifest(ctx, ref)
}

func TagImage(manifest ImageManifest, location string) error {
	ref, err := reference.ParseNormalizedNamed(location)
	if err != nil {
		return err
	}

	ctx := context.Background()
	_, err = newClient().PutManifest(ctx, ref, manifest)
	return err
}

func newClient() client.RegistryClient {
	return client.NewRegistryClient(resolver, agent, false)
}

func resolver(ctx context.Context, index *registry.IndexInfo) types.AuthConfig {
	conf := config.LoadDefaultConfigFile(os.Stderr)
	authConfig, _ := conf.GetAuthConfig(index.Name)
	return types.AuthConfig(authConfig)
}

const agent = "KD (" + runtime.GOOS + ")"
