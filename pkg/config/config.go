package config

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
)

type ResolvedApp struct {
	App
	Tag      string
	Registry string
}

type ResolvedTarget struct {
	Target
}

const DefaultTag = "latest"
const ConfigName = "kdeploy.conf"

func Load() (*Config, error) {
	return loadFromFs(&afero.Afero{Fs: afero.NewOsFs()})
}

func loadFromFs(afs *afero.Afero) (*Config, error) {
	bytes, err := afs.ReadFile(ConfigName)
	if err != nil {
		return nil, errors.Wrap(err, "Config error")
	}

	var conf Config
	if err := yaml.UnmarshalStrict(bytes, &conf); err != nil {
		return nil, errors.Wrap(err, "Config error")
	}

	if conf.ApiVersion > LatestVersion {
		return nil, errors.Errorf("Unsupported configuration version %d, please update kd!", conf.ApiVersion)
	}

	return &conf, nil
}

func (conf *Config) AppNames() []string {
	names := make([]string, len(conf.Apps))

	for i, app := range conf.Apps {
		names[i] = app.Name
	}

	return names
}

func (conf *Config) ResolveApp(name string) (*ResolvedApp, error) {
	/* TODO: No naming conflicts are checked yet. Returns the first match. */
	parts := strings.Split(name, ":")
	tag := DefaultTag
	name = parts[0]
	if len(parts) > 1 {
		tag = parts[1]
	}

	if name == "" && len(conf.Apps) != 1 {
		return nil, fmt.Errorf("Selecting default requires exactly 1 application (%d configured)", len(conf.Apps))
	} else {
		for _, app := range conf.Apps {
			if app.Name == "" {
				parts := strings.Split(app.Path, "/")
				app.Name = parts[len(parts)-1]
			}

			if name == "" || name == app.Name {
				return &ResolvedApp{
					App:      app,
					Tag:      tag,
					Registry: conf.Registry,
				}, nil
			}
		}

		return nil, fmt.Errorf("Unknown application '%s'", name)
	}
}

func (conf *Config) TargetNames() []string {
	names := make([]string, len(conf.Targets))

	for i, tgt := range conf.Targets {
		names[i] = tgt.Name
	}

	return names
}

func (conf *Config) ResolveTarget(name string) (*ResolvedTarget, error) {
	/* TODO: No naming conflicts are checked yet. Returns the first match. */
	for _, tgt := range conf.Targets {
		if tgt.Name == name || stupidContains(tgt.Alias, name) {
			return &ResolvedTarget{
				Target: tgt,
			}, nil
		}
	}

	return nil, fmt.Errorf("Unknown target '%s'", name)
}

func (app *ResolvedApp) Repository() string {
	return app.RepositoryWithTag(app.Tag)
}

func (app *ResolvedApp) RepositoryWithTag(tag string) string {
	return app.Registry + "/" + app.Name + ":" + tag
}

func (app *ResolvedApp) RepositoryWithDigest(digest string) string {
	return app.Registry + "/" + app.Name + "@" + digest
}

func (a *StringArray) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var multi []string
	err := unmarshal(&multi)
	if err != nil {
		var single string
		err := unmarshal(&single)
		if err != nil {
			return err
		}
		*a = []string{single}
	} else {
		*a = multi
	}
	return nil
}

func stupidContains(slice []string, search string) bool {
	for _, item := range slice {
		if item == search {
			return true
		}
	}
	return false
}
