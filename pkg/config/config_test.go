package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
	assert.Equal(t, "Unknown app 'bar'", err.Error())
}

func TestResolveApps(t *testing.T) {
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

	apps := conf.ResolveApps()
	assert.Equal(t, 2, len(apps))
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
