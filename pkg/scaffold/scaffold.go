package scaffold

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"

	"github.com/spf13/afero"
	"github.com/voormedia/kd/pkg/config"
	"github.com/voormedia/kd/pkg/util"
)

func Run(log *util.Logger) error {
	log.Note("Please enter a few project details")

	fs := &afero.Afero{Fs: afero.NewOsFs()}
	details, err := requestDetails(os.Stdin, os.Stdout)
	if err != nil {
		return err
	}

	err = writeConfig(fs, details)
	if err != nil {
		return err
	}

	return nil
}

type details struct {
	Customer string
	Name     string
	Project  string
	Context  string
}

func requestDetails(in io.Reader, out io.Writer) (data *details, err error) {
	reader := bufio.NewReader(in)

	fmt.Fprint(out, "Name of customer: ")
	customer, err := reader.ReadString('\n')
	if err != nil {
		return
	}

	fmt.Fprint(out, "Name of main application: ")
	name, err := reader.ReadString('\n')
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

	data = &details{
		Customer: util.Slugify(customer),
		Name:     util.Slugify(name),
		Project:  strings.TrimSpace(project),
		Context:  strings.TrimSpace(context),
	}

	return
}

var kdeploy = template.Must(template.New(config.ConfigName).Parse(
	`# Private docker registry to push images to
registry: eu.gcr.io/{{.Project}}/{{.Customer}}

# List of apps to build
apps:
- name: {{.Name}}
  path: .

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
