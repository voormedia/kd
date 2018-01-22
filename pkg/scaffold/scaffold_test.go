package scaffold

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

var fs = &afero.Afero{Fs: afero.NewMemMapFs()}

func TestRequestDetails(t *testing.T) {
	input := bytes.NewBufferString(strings.Join([]string{
		"a Customer Name\n",
		"Main Application\n",
		"project-123456\n",
		"cluster_Context\n",
	}, ""))

	data, err := requestDetails(input, ioutil.Discard)
	assert.Nil(t, err)
	assert.Equal(t, &details{
		Customer: "a-customer-name",
		Name:     "main-application",
		Project:  "project-123456",
		Context:  "cluster_Context",
	}, data)
}

func TestWriteConfig(t *testing.T) {
	details := &details{
		Customer: "a-customer-name",
		Name:     "main-application",
		Project:  "project-123456",
		Context:  "cluster_Context",
	}

	err := writeConfig(fs, details)
	assert.Nil(t, err)

	data, err := fs.ReadFile("kdeploy.conf")
	assert.Nil(t, err)
	assert.Equal(t, strings.Join([]string{
		"# Private docker registry to push images to\n",
		"registry: eu.gcr.io/project-123456/a-customer-name\n",
		"\n",
		"# List of apps to build\n",
		"apps:\n",
		"- name: main-application\n",
		"  path: .\n",
		"\n",
		"# List of available deployment targets\n",
		"targets:\n",
		"- name: acceptance\n",
		"  context: cluster_Context\n",
		"  namespace: a-customer-name-acc\n",
		"  path: config/deploy/acceptance\n",
		"\n",
		"- name: production\n",
		"  context: cluster_Context\n",
		"  namespace: a-customer-name-prd\n",
		"  path: config/deploy/production\n",
	}, ""), string(data))
}
