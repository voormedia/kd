package deploy

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"runtime"

	cliconfig "github.com/docker/cli/cli/config"
	manifesttypes "github.com/docker/cli/cli/manifest/types"
	"github.com/docker/cli/cli/registry/client"
	"github.com/docker/distribution/reference"
	"github.com/docker/docker/api/types"
	registrytypes "github.com/docker/docker/api/types/registry"
	"github.com/voormedia/kd/pkg/config"
	"github.com/voormedia/kd/pkg/kubectl"
	"github.com/voormedia/kd/pkg/util"
	outil "k8s.io/kubectl/pkg/kinflate"
	kutil "k8s.io/kubectl/pkg/kinflate/util"
)

func Run(verbose bool, log *util.Logger, app *config.ResolvedApp, target *config.ResolvedTarget) error {
	log.Note("Retrieving image", app.Name+":"+app.Tag)
	img, err := getImage(app.Repository())
	if err != nil {
		return err
	}

	log.Note("Applying configuration")
	err = apply(app, target, img)
	if err != nil {
		return err
	}

	log.Note("Tagging image", app.Name+":"+target.Name)
	err = tagImage(img, app.RepositoryWithTag(target.Name))
	if err != nil {
		return err
	}

	log.Success("Successfully deployed", app.Repository(), "to", target.Name)
	return nil
}

func getImage(location string) (manifesttypes.ImageManifest, error) {
	ref, err := reference.ParseNormalizedNamed(location)
	if err != nil {
		return manifesttypes.ImageManifest{}, err
	}

	ctx := context.Background()
	return newClient().GetManifest(ctx, ref)
}

func apply(app *config.ResolvedApp, target *config.ResolvedTarget, img manifesttypes.ImageManifest) error {
	manifest, err := outil.LoadFromManifestPath(filepath.Join(app.Path, target.Path))
	if err != nil {
		return err
	}

	res, err := kutil.Encode(manifest)
	if err != nil {
		return err
	}

	/* HACK to set deployment image. */
	url := app.RepositoryWithDigest(img.Digest.String())
	buf := bytes.NewBuffer(bytes.Replace(res, []byte(" image: "+app.Name), []byte(" image: "+url), -1))

	// os.Stdout.Write(buf.Bytes())
	return kubectl.Apply(target.Context, target.Namespace, buf, os.Stdout, os.Stderr, &kubectl.ApplyOptions{})
}

func tagImage(manifest manifesttypes.ImageManifest, location string) error {
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

func resolver(ctx context.Context, index *registrytypes.IndexInfo) types.AuthConfig {
	conf := cliconfig.LoadDefaultConfigFile(os.Stderr)
	authConfig, _ := conf.GetAuthConfig(index.Name)
	return authConfig
}

const agent = "KD (" + runtime.GOOS + ")"
