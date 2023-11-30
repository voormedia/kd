package config

/* Maximum configuration version accepted by this version of kd.
   Increment this version on API incompatible changes. */
const LatestVersion uint = 2

type StringArray []string

type App struct {
	Name      string `yaml:"name,omitempty"`
	Path      string `yaml:"path,omitempty"`
	Root      string `yaml:"root,omitempty"`
	Default   bool   `yaml:"default,omitempty"`
	Platform  string `yaml:"platform,omitempty"`
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
	ApiVersion uint     `yaml:"version,omitempty"`
	Registry   string   `yaml:"registry,omitempty"`
	Apps       []App    `yaml:"apps,omitempty"`
	Targets    []Target `yaml:"targets,omitempty"`
}
