package scaffold

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/afero"
	"github.com/voormedia/kd/pkg/config"
	"github.com/voormedia/kd/pkg/util"
)

func Run(log *util.Logger) error {
	log.Note("Please enter a few project details")

	fs := &afero.Afero{Fs: afero.NewOsFs()}

	details, err := requestDetails(fs, os.Stdin, os.Stdout)
	if err != nil {
		return err
	}

	err = writeConfig(fs, details)
	if err != nil {
		return err
	}

	return nil
}

type app struct {
	Name string
	Path string
}

type details struct {
	Apps     []app
	Customer string
	Project  string
	Context  string
}

func requestDetails(afs *afero.Afero, in io.Reader, out io.Writer) (data *details, err error) {
	reader := bufio.NewReader(in)

	fmt.Fprint(out, "Name of customer: ")
	customer, err := reader.ReadString('\n')
	if err != nil {
		return
	}

	fmt.Fprint(out, "Google Cloud project id: ")
	project, err := reader.ReadString('\n')
	if err != nil {
		return
	}

	fmt.Fprint(out, "Kubernetes cluster context: ")
	context, err := reader.ReadString('\n')
	if err != nil {
		return
	}

	apps, err := findApps(afs)
	if err != nil {
		return
	}

	data = &details{
		Apps:     apps,
		Customer: util.Slugify(customer),
		Project:  strings.TrimSpace(project),
		Context:  strings.TrimSpace(context),
	}

	return
}

func findApps(afs *afero.Afero) ([]app, error) {
	root, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	var apps []app
	afs.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.Name() == "Dockerfile" {
			path := filepath.Dir(path)
			name := filepath.Base(path)

			path = strings.Replace(path, root+"", "", 1)
			path = strings.Replace(path, "/", "", 1)
			if path == "" {
				path = "."
			}

			apps = append(apps, app{
				Path: path,
				Name: name,
			})
		}

		return nil
	})

	return apps, nil
}

var kdeploy = template.Must(template.New(config.ConfigName).Parse(
	`# Private docker registry to push images to
registry: eu.gcr.io/{{.Project}}/{{.Customer}}

# List of apps to build
apps:
{{- range .Apps}}
- name: {{.Name}}
  path: {{.Path}}
{{- end}}

# List of available deployment targets
targets:
- name: acceptance
  context: {{.Context}}
  namespace: {{.Customer}}-acc
  path: config/deploy/acceptance

- name: production
  context: {{.Context}}
  namespace: {{.Customer}}-prd
  path: config/deploy/production
`))

func writeConfig(afs *afero.Afero, details *details) error {
	var buf bytes.Buffer
	if err := kdeploy.Execute(&buf, details); err != nil {
		return err
	}

	return afs.WriteFile(config.ConfigName, buf.Bytes(), 0644)
}
