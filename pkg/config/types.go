package config

type App struct {
	Name string `yaml:"name,omitempty"`
	Path string `yaml:"path,omitempty"`
	Root string `yaml:"root,omitempty"`
}

type Target struct {
	Name      string `yaml:"name,omitempty"`
	Context   string `yaml:"context,omitempty"`
	Namespace string `yaml:"namespace,omitempty"`
	Path      string `yaml:"path,omitempty"`
}

type Config struct {
	Registry string   `yaml:"registry,omitempty"`
	Apps     []App    `yaml:"apps,omitempty"`
	Targets  []Target `yaml:"targets,omitempty"`
}
