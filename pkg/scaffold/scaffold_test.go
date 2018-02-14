package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

var fs = &afero.Afero{Fs: afero.NewMemMapFs()}

func TestFindApps(t *testing.T) {
	root, _ := os.Getwd()
	root, _ = filepath.EvalSymlinks(root)

	fs.WriteFile(root+"/Dockerfile", []byte{}, 0644)
	fs.MkdirAll(root+"/foobar/baz", 0644)
	fs.WriteFile(root+"/foobar/baz/Dockerfile", []byte{}, 0644)

	apps, err := findApps(fs)
	assert.Nil(t, err)
	assert.Equal(t, []app{{
		Name: "scaffold", // PWD
		Path: ".",
	}, {
		Name: "baz",
		Path: "foobar/baz",
	}}, apps)
}

func TestWriteConfig(t *testing.T) {
	details := &details{
		Apps: []app{{
			Name: "my-website",
			Path: ".",
		}, {
			Name: "other-app",
			Path: "apps/other-app",
		}},
		Customer: "a-customer-name",
		Project:  "project-123456",
		Context:  "cluster_Context",
	}

	err := writeConfig(fs, details)
	assert.Nil(t, err)

	kdeploy, err := fs.ReadFile("kdeploy.conf")
	assert.Nil(t, err)
	assert.Equal(t, strings.Join([]string{
		"# Private docker registry to push images to\n",
		"registry: eu.gcr.io/project-123456/a-customer-name\n",
		"\n",
		"# List of apps to build\n",
		"apps:\n",
		"- name: my-website\n",
		"  path: .\n",
		"- name: other-app\n",
		"  path: apps/other-app\n",
		"\n",
		"# List of available deployment targets\n",
		"targets:\n",
		"- name: acceptance\n",
		"  alias: acc\n",
		"  context: cluster_Context\n",
		"  namespace: a-customer-name-acc\n",
		"  path: config/deploy/acceptance\n",
		"\n",
		"- name: production\n",
		"  alias: prd\n",
		"  context: cluster_Context\n",
		"  namespace: a-customer-name-prd\n",
		"  path: config/deploy/production\n",
	}, ""), string(kdeploy))

	bseManifest, err := fs.ReadFile("config/deploy/kube-manifest.yaml")
	assert.Nil(t, err)
	assert.Equal(t, strings.Join([]string{
		"# List of base resources\n",
		"resources:\n",
		"- deployment.yaml\n",
		"- service.yaml\n",
		"- ingress.yaml\n",
	}, ""), string(bseManifest))

	accManifest, err := fs.ReadFile("config/deploy/acceptance/kube-manifest.yaml")
	assert.Nil(t, err)
	assert.Equal(t, strings.Join([]string{
		"# List of patches to apply (in order) for this environment\n",
		"patches:\n",
		"- namespace.yaml\n",
		"- deployment.yaml\n",
		"- ingress.yaml\n",
		"\n",
		"# Patches are applied to base resources\n",
		"resources:\n",
		"- ..\n",
	}, ""), string(accManifest))
}
