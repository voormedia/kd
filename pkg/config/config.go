package config

import (
	"fmt"
	"io/ioutil"
	"strings"

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
	bytes, err := ioutil.ReadFile(ConfigName)
	if err != nil {
		return nil, err
	}

	var pkg Config
	err = yaml.UnmarshalStrict(bytes, &pkg)
	return &pkg, err
}

func AppNames() []string {
	config, err := Load()
	if err != nil {
		return []string{}
	}
	return config.AppNames()
}

func TargetNames() []string {
	config, err := Load()
	if err != nil {
		return []string{}
	}
	return config.TargetNames()
}

func ResolveAppNames(names []string) ([]*ResolvedApp, error) {
	config, err := Load()
	if err != nil {
		return nil, err
	}

	var apps []*ResolvedApp
	if len(names) == 0 {
		apps = config.ResolveApps()
	} else {
		for _, name := range names {
			app, err := config.ResolveApp(name)
			if err != nil {
				return nil, err
			}

			apps = append(apps, app)
		}
	}

	return apps, nil
}

func ResolveTargetName(name string) (*ResolvedTarget, error) {
	config, err := Load()
	if err != nil {
		return nil, err
	}

	return config.ResolveTarget(name)
}

func (conf *Config) AppNames() []string {
	names := make([]string, len(conf.Apps))

	for i, app := range conf.Apps {
		names[i] = app.Name
	}

	return names
}

func (conf *Config) ResolveApps() []*ResolvedApp {
	apps := make([]*ResolvedApp, len(conf.Apps))

	for i, app := range conf.Apps {
		apps[i] = &ResolvedApp{
			App:      app,
			Tag:      DefaultTag,
			Registry: conf.Registry,
		}
	}

	return apps
}

func (conf *Config) ResolveApp(name string) (*ResolvedApp, error) {
	parts := strings.Split(name, ":")
	tag := DefaultTag
	name = parts[0]
	if len(parts) > 1 {
		tag = parts[1]
	}

	for _, app := range conf.Apps {
		if app.Name == name {
			return &ResolvedApp{
				App:      app,
				Tag:      tag,
				Registry: conf.Registry,
			}, nil
		}
	}

	return nil, fmt.Errorf("Unknown app '%s'", name)
}

func (conf *Config) TargetNames() []string {
	names := make([]string, len(conf.Targets))

	for i, tgt := range conf.Targets {
		names[i] = tgt.Name
	}

	return names
}

func (conf *Config) ResolveTarget(name string) (*ResolvedTarget, error) {
	for _, tgt := range conf.Targets {
		if tgt.Name == name {
			return &ResolvedTarget{
				Target: tgt,
			}, nil
		}
	}

	return nil, fmt.Errorf("Unknown target '%s'", name)
}

func (app *ResolvedApp) Repository() string {
	return app.TaggedRepository(app.Tag)
}

func (app *ResolvedApp) TaggedRepository(tag string) string {
	reg := app.Registry + "/" + app.Name
	if tag != DefaultTag {
		reg = reg + ":" + tag
	}
	return reg
}
