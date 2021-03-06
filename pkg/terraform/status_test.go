package terraform

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/deislabs/porter/pkg/printer"
	"github.com/deislabs/porter/pkg/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	yaml "gopkg.in/yaml.v2"
)

type statusTest struct {
	format                printer.Format
	expectedCommandSuffix string
}

func TestMixin_UnmarshalStatusStep(t *testing.T) {
	b, err := ioutil.ReadFile("testdata/status-input.yaml")
	require.NoError(t, err)

	var action StatusAction
	err = yaml.Unmarshal(b, &action)
	require.NoError(t, err)
	require.Len(t, action.Steps, 1)
	step := action.Steps[0]

	assert.Equal(t, "Status MySQL", step.Description)
}

func TestMixin_Status(t *testing.T) {
	os.Setenv(test.ExpectedCommandEnv, strings.Join([]string{
		"terraform init",
		"terraform show",
	}, "\n"))

	statusStep := StatusStep{
		StatusArguments: StatusArguments{},
	}

	action := StatusAction{Steps: []StatusStep{statusStep}}
	b, err := yaml.Marshal(action)
	require.NoError(t, err)

	h := NewTestMixin(t)
	h.In = bytes.NewReader(b)

	// Set up working dir as current dir
	h.WorkingDir, err = os.Getwd()
	require.NoError(t, err)

	err = h.Status()
	require.NoError(t, err)

	wd, err := os.Getwd()
	require.NoError(t, err)
	assert.Equal(t, wd, h.WorkingDir)
}
