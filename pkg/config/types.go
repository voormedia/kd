package config

type StringArray []string

type App struct {
	Name      string `yaml:"name,omitempty"`
	Path      string `yaml:"path,omitempty"`
	Root      string `yaml:"root,omitempty"`
	PreBuild  string `yaml:"preBuild,omitempty"`
	PostBuild string `yaml:"postBuild,omitempty"`
}

type Target struct {
	Name      string      `yaml:"name,omitempty"`
	Alias     StringArray `yaml:"alias,omitempty"`
	Context   string      `yaml:"context,omitempty"`
	Namespace string      `yaml:"namespace,omitempty"`
	Path      string      `yaml:"path,omitempty"`
}

type Config struct {
	Registry string   `yaml:"registry,omitempty"`
	Apps     []App    `yaml:"apps,omitempty"`
	Targets  []Target `yaml:"targets,omitempty"`
}
