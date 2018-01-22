package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type ResolvedApp struct {
	App
	Registry string
}

type ResolvedTarget struct {
	Target
}

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

func (conf *Config) ResolveApps() []*ResolvedApp {
	apps := make([]*ResolvedApp, len(conf.Apps))

	for i, app := range conf.Apps {
		apps[i] = &ResolvedApp{
			App:      app,
			Registry: conf.Registry,
		}
	}

	return apps
}

func (conf *Config) ResolveApp(name string) (*ResolvedApp, error) {
	for _, app := range conf.Apps {
		if app.Name == name {
			return &ResolvedApp{
				App:      app,
				Registry: conf.Registry,
			}, nil
		}
	}

	return nil, fmt.Errorf("Unknown app '%s'", name)
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

func (app *ResolvedApp) Tag() string {
	return fmt.Sprintf("%s/%s", app.Registry, app.Name)
}
