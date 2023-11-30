package config

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
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
	return LoadFromFs(&afero.Afero{Fs: afero.NewOsFs()})
}

func LoadFromFs(afs *afero.Afero) (*Config, error) {
	conf, err := GetRawFromFS(afs)
	if err != nil {
		return nil, err
	}

	if conf.ApiVersion > LatestVersion {
		return nil, errors.Errorf("Unsupported configuration version %d, please get the latest version of kd", conf.ApiVersion)
	}

	if conf.ApiVersion == 1 {
		return nil, errors.Errorf("Unsupported configuration version %d, please run 'kd upgrade' to upgrade to version %d", conf.ApiVersion, LatestVersion)
	}

	defaults := 0
	for i := range conf.Apps {
		app := &conf.Apps[i]

		if app.Default {
			defaults += 1
		}

		if app.Name == "" {
			parts := strings.Split(app.Path, "/")
			app.Name = parts[len(parts)-1]
		}

		if app.Root == "" {
			app.Root = app.Path
		}

		if app.Platform == "" {
			app.Platform = "linux/amd64"
		}
	}

	if defaults > 1 {
		return nil, fmt.Errorf("Only one application may be marked with 'default: true'")
	}

	return conf, nil
}

func GetRawFromFS(afs *afero.Afero) (*Config, error) {
	bytes, err := afs.ReadFile(ConfigName)
	if err != nil {
		return nil, errors.Wrap(err, "Config error")
	}

	var conf Config
	if err := yaml.Unmarshal(bytes, &conf); err != nil {
		return nil, errors.Wrap(err, "Config error")
	}

	return &conf, nil
}

func (conf *Config) AppNames() (names []string) {
	for _, app := range conf.Apps {
		names = append(names, app.Name)
	}
	return
}

func (conf *Config) ResolveApp(name string, tag string) (*ResolvedApp, error) {
	/* TODO: No naming conflicts are checked yet. Returns the first match. */
	parts := strings.Split(name, ":")
	name = parts[0]

	if tag == "" {
		tag = DefaultTag
		if len(parts) > 1 {
			tag = parts[1]
		}
	} else if len(parts) > 1 {
		return nil, fmt.Errorf("Specify a tag either with 'app:%s' or with '--tag %s', but not both", parts[1], tag)
	}

	for _, app := range conf.Apps {
		if (name == "" && (app.Default || len(conf.Apps) == 1)) || name == app.Name {
			return &ResolvedApp{
				App:      app,
				Tag:      tag,
				Registry: conf.Registry,
			}, nil
		}
	}

	if len(conf.Apps) == 0 {
		return nil, fmt.Errorf("No applications configured")
	}

	if name == "" {
		return nil, fmt.Errorf("Selecting default requires one application to be marked with 'default: true'")
	}

	return nil, fmt.Errorf("Unknown application '%s'", name)
}

func (conf *Config) TargetNames() (names []string) {
	for _, tgt := range conf.Targets {
		names = append(names, tgt.Name)
	}
	return
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
