package scaffold

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestDetails(t *testing.T) {
	input := bytes.NewBufferString(strings.Join([]string{
		"a Customer Name\n",
		"Main Application\n",
		"project-123456\n",
		"cluster_Context\n",
	}, ""))

	data, err := requestDetails(input, ioutil.Discard)
	assert.Nil(t, err)
	assert.Equal(t, data, &details{
		Customer: "a-customer-name",
		Name:     "main-application",
		Project:  "project-123456",
		Context:  "cluster_Context",
	})
}
