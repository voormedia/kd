package build

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/command/image/build"
	cliconfig "github.com/docker/cli/cli/config"
	"github.com/docker/distribution/reference"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/idtools"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/registry"
	"github.com/pkg/errors"
	"github.com/voormedia/kd/pkg/config"
	"github.com/voormedia/kd/pkg/util"
)

func buildImage(verbose bool, app *config.ResolvedApp) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cli, err := client.NewClient(client.DefaultDockerHost, "1.35", nil, nil)
	if err != nil {
		return err
	}

	if app.Root == "" {
		app.Root = app.Path
	}

	dockerfile := filepath.Join(app.Path, build.DefaultDockerfileName)
	dir, dockerfile, err := build.GetContextFromLocalDir(app.Root, dockerfile)
	excludes, err := build.ReadDockerignore(dir)
	if err != nil {
		return err
	}

	if err := build.ValidateContextDirectory(dir, excludes); err != nil {
		return errors.Errorf("Error checking context: '%s'", err)
	}

	dockerfile, err = archive.CanonicalTarNameForPath(dockerfile)
	if err != nil {
		return errors.Errorf("Cannot canonicalize dockerfile path %s: %v", dockerfile, err)
	}

	excludes = build.TrimBuildFilesFromExcludes(excludes, dockerfile, false)
	build, err := archive.TarWithOptions(dir, &archive.TarOptions{
		ChownOpts:       &idtools.IDPair{UID: 0, GID: 0},
		ExcludePatterns: excludes,
	})

	if err != nil {
		return err
	}

	opt := types.ImageBuildOptions{
		Dockerfile:  dockerfile,
		ForceRemove: true,
		PullParent:  true,
		Tags:        []string{app.Tag()},
	}

	res, err := cli.ImageBuild(ctx, build, opt)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if verbose {
		return jsonmessage.DisplayJSONMessagesStream(res.Body, os.Stderr, os.Stderr.Fd(), true, nil)
	} else {
		var buf bytes.Buffer
		err := jsonmessage.DisplayJSONMessagesStream(res.Body, &buf, os.Stderr.Fd(), true, nil)
		if err != nil {
			os.Stderr.Write(buf.Bytes())
			return err
		}

		return nil
	}
}

func pushImage(verbose bool, app *config.ResolvedApp) error {
	/* TODO: Figure out how to directly integrate with gcloud docker tool? */
	cmd := exec.Command("gcloud", "docker", "--authorize-only")
	var out bytes.Buffer
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return errors.Errorf("Failed to authorize: %s", out.String())
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cli, err := client.NewClient(client.DefaultDockerHost, "1.35", nil, nil)
	if err != nil {
		return err
	}

	ref, err := reference.ParseNormalizedNamed(app.Tag())
	if err != nil {
		return err
	}

	repoInfo, err := registry.ParseRepositoryInfo(ref)
	if err != nil {
		return err
	}

	conf := cliconfig.LoadDefaultConfigFile(os.Stderr)
	authConfig, err := conf.GetAuthConfig(repoInfo.Index.Name)
	if err != nil {
		return err
	}

	encodedAuth, err := command.EncodeAuthToBase64(authConfig)
	if err != nil {
		return err
	}

	opt := types.ImagePushOptions{
		RegistryAuth: encodedAuth,
	}

	res, err := cli.ImagePush(ctx, reference.FamiliarString(ref), opt)
	if err != nil {
		return errors.Errorf("Failed to push image: %s", err)
	}

	defer res.Close()

	return jsonmessage.DisplayJSONMessagesStream(res, os.Stderr, os.Stderr.Fd(), true, nil)
}

func Run(verbose bool, log *util.Logger, apps []*config.ResolvedApp) error {
	for _, app := range apps {
		log.Note("Building", app.Name)
		if err := buildImage(verbose, app); err != nil {
			log.Fatal(err)
		}

		log.Note("Pushing", app.Name)
		if err := pushImage(verbose, app); err != nil {
			log.Fatal(err)
		}

		log.Success("Successfully built", app.Tag())
	}

	return nil
}
