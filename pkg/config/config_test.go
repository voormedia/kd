package config

import (
	"strings"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	fs := &afero.Afero{Fs: afero.NewMemMapFs()}

	fs.WriteFile("kdeploy.conf", []byte(strings.Join([]string{
		"# Private docker registry to push images to\n",
		"registry: eu.gcr.io/project-123456/a-customer-name\n",
		"\n",
		"# List of apps to build\n",
		"apps:\n",
		"- name: my-website\n",
		"  path: .\n",
		"- name: other-app\n",
		"  path: apps/other-app\n",
		"  root: apps\n",
		"  preBuild: script/foo.sh\n",
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
		"  context: cluster_Context\n",
		"  namespace: a-customer-name-prd\n",
		"  path: config/deploy/production\n",
	}, "")), 0644)

	conf, err := loadFromFs(fs)
	assert.Nil(t, err)

	expected := &Config{
		Registry: "eu.gcr.io/project-123456/a-customer-name",
		Apps: []App{{
			Name: "my-website",
			Path: ".",
			Root: "",
		}, {
			Name:     "other-app",
			Path:     "apps/other-app",
			Root:     "apps",
			PreBuild: "script/foo.sh",
		}},
		Targets: []Target{{
			Name:      "acceptance",
			Alias:     []string{"acc"},
			Context:   "cluster_Context",
			Namespace: "a-customer-name-acc",
			Path:      "config/deploy/acceptance",
		}, {
			Name:      "production",
			Context:   "cluster_Context",
			Namespace: "a-customer-name-prd",
			Path:      "config/deploy/production",
		}},
	}

	assert.Equal(t, expected, conf)
}

func TestLoadError(t *testing.T) {
	fs := &afero.Afero{Fs: afero.NewMemMapFs()}
	fs.WriteFile("kdeploy.conf", []byte("Bad file format"), 0644)

	conf, err := loadFromFs(fs)
	assert.Nil(t, conf)
	assert.Equal(t, "Config error: yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `Bad fil...` into config.Config", err.Error())
}

func TestLoadMissing(t *testing.T) {
	fs := &afero.Afero{Fs: afero.NewMemMapFs()}
	conf, err := loadFromFs(fs)
	assert.Nil(t, conf)
	assert.Equal(t, "Config error: open kdeploy.conf: file does not exist", err.Error())
}

func TestRepositoryWithTag(t *testing.T) {
	app := &ResolvedApp{
		App: App{
			Name: "foo",
			Path: "apps/foo",
			Root: ".",
		},
		Tag:      "bar",
		Registry: "my.registry.com",
	}

	assert.Equal(t, "my.registry.com/foo:bar", app.Repository())
}

func TestRepositoryWithDefaultTag(t *testing.T) {
	app := &ResolvedApp{
		App: App{
			Name: "foo",
			Path: "apps/foo",
			Root: ".",
		},
		Tag:      "latest",
		Registry: "my.registry.com",
	}

	assert.Equal(t, "my.registry.com/foo", app.Repository())
}

func TestResolveExistingApp(t *testing.T) {
	conf := &Config{
		Registry: "my.registry.com",
		Apps: []App{{
			Name: "foo",
			Path: "apps/foo",
			Root: "",
		}},
	}

	app, err := conf.ResolveApp("foo")
	assert.Nil(t, err)
	assert.Equal(t, "latest", app.Tag)
	assert.Equal(t, "apps/foo", app.Path)
}

func TestResolveExistingAppExplicitTag(t *testing.T) {
	conf := &Config{
		Registry: "my.registry.com",
		Apps: []App{{
			Name: "foo",
			Path: "apps/foo",
			Root: "",
		}},
	}

	app, err := conf.ResolveApp("foo:my-tag")
	assert.Nil(t, err)
	assert.Equal(t, "my-tag", app.Tag)
	assert.Equal(t, "my.registry.com/foo:my-tag", app.Repository())
}

func TestResolveExistingAppExplicitLatest(t *testing.T) {
	conf := &Config{
		Registry: "my.registry.com",
		Apps: []App{{
			Name: "foo",
			Path: "apps/foo",
			Root: "",
		}},
	}

	app, err := conf.ResolveApp("foo:latest")
	assert.Nil(t, err)
	assert.Equal(t, "latest", app.Tag)
	assert.Equal(t, "my.registry.com/foo", app.Repository())
}

func TestResolveMissingApp(t *testing.T) {
	conf := &Config{
		Registry: "my.registry.com",
		Apps: []App{{
			Name: "foo",
			Path: "apps/foo",
			Root: "",
		}},
	}

	app, err := conf.ResolveApp("bar")
	assert.Nil(t, app)
	assert.Equal(t, "Unknown application 'bar'", err.Error())
}

func TestResolveDefaultAppForSingle(t *testing.T) {
	conf := &Config{
		Registry: "my.registry.com",
		Apps: []App{{
			Name: "foo",
			Path: "apps/foo",
			Root: "",
		}},
	}

	app, err := conf.ResolveApp("")
	assert.Nil(t, err)
	assert.Equal(t, "latest", app.Tag)
	assert.Equal(t, "my.registry.com/foo", app.Repository())
}

func TestResolveDefaultAppForMultiple(t *testing.T) {
	conf := &Config{
		Registry: "my.registry.com",
		Apps: []App{{
			Name: "foo",
			Path: "apps/foo",
			Root: "",
		}, {
			Name: "bar",
			Path: "apps/bar",
			Root: "apps",
		}},
	}

	app, err := conf.ResolveApp("")
	assert.Nil(t, app)
	assert.Equal(t, "Selecting default requires exactly 1 application (2 configured)", err.Error())
}

func TestResolveDefaultAppForNone(t *testing.T) {
	conf := &Config{
		Registry: "my.registry.com",
		Apps:     []App{},
	}

	app, err := conf.ResolveApp("")
	assert.Nil(t, app)
	assert.Equal(t, "Selecting default requires exactly 1 application (0 configured)", err.Error())
}

func TestResolveExistingTargetAlias(t *testing.T) {
	conf := &Config{
		Registry: "my.registry.com",
		Targets: []Target{{
			Name:      "acc",
			Alias:     []string{"test", "ac"},
			Context:   "cluster",
			Namespace: "acc",
			Path:      "config/deploy/acc",
		}},
	}

	tgt, err := conf.ResolveTarget("ac")
	assert.Nil(t, err)
	assert.Equal(t, "acc", tgt.Namespace)
}

func TestResolveExistingTarget(t *testing.T) {
	conf := &Config{
		Registry: "my.registry.com",
		Targets: []Target{{
			Name:      "acc",
			Context:   "cluster",
			Namespace: "acc",
			Path:      "config/deploy/acc",
		}},
	}

	tgt, err := conf.ResolveTarget("acc")
	assert.Nil(t, err)
	assert.Equal(t, "acc", tgt.Namespace)
}

func TestResolveMissingTarget(t *testing.T) {
	conf := &Config{
		Registry: "my.registry.com",
		Targets: []Target{{
			Name:      "acc",
			Context:   "cluster",
			Namespace: "acc",
			Path:      "config/deploy/acc",
		}},
	}

	tgt, err := conf.ResolveTarget("prd")
	assert.Nil(t, tgt)
	assert.Equal(t, "Unknown target 'prd'", err.Error())
}
