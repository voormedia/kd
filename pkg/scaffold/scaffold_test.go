package scaffold

import (
	"bytes"
	// "fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

var fs = &afero.Afero{Fs: afero.NewMemMapFs()}

func TestRequestDetails(t *testing.T) {
	input := bytes.NewBufferString(strings.Join([]string{
		"a Customer Name\n",
		"project-123456\n",
		"cluster_Context\n",
	}, ""))

	data, err := requestDetails(fs, input, ioutil.Discard)
	assert.Nil(t, err)
	assert.Equal(t, &details{
		Apps:     nil,
		Customer: "a-customer-name",
		Project:  "project-123456",
		Context:  "cluster_Context",
	}, data)
}

func TestRequestDetailsWithApps(t *testing.T) {
	root, _ := os.Getwd()
	fs.WriteFile(root+"/Dockerfile", []byte{}, 0644)
	fs.MkdirAll(root+"/foobar/baz", 0644)
	fs.WriteFile(root+"/foobar/baz/Dockerfile", []byte{}, 0644)

	input := bytes.NewBufferString(strings.Join([]string{
		"a Customer Name\n",
		"project-123456\n",
		"cluster_Context\n",
	}, ""))

	data, err := requestDetails(fs, input, ioutil.Discard)
	assert.Nil(t, err)
	assert.Equal(t, &details{
		Apps: []app{{
			Name: "scaffold", // PWD
			Path: ".",
		}, {
			Name: "baz",
			Path: "foobar/baz",
		}},
		Customer: "a-customer-name",
		Project:  "project-123456",
		Context:  "cluster_Context",
	}, data)
}

func TestWriteConfig(t *testing.T) {
	details := &details{
		Apps: []app{{
			Name: "my-website",
			Path: ".",
		}, {
			Name: "other-app",
			Path: "apps/other-app",
		}},
		Customer: "a-customer-name",
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
		"- name: my-website\n",
		"  path: .\n",
		"- name: other-app\n",
		"  path: apps/other-app\n",
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
