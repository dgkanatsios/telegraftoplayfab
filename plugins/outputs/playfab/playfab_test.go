package playfab

import (
	"fmt"
	"testing"

	"gotest.tools/v3/assert"
)

func TestInit(t *testing.T) {
	p := &PlayFab{}
	err := p.Init()
	assert.Error(t, err, errTitleIdRequired)

	p.TitleId = "test"
	err = p.Init()
	assert.Error(t, err, errDeveloperSecretKeyOrTelemetryKey)

	p.DeveloperSecretKey = "test"
	p.TelemetryKey = "lala"
	err = p.Init()
	assert.Error(t, err, errNotDeveloperSecretKeyAndTelemetryKey)

	p.DeveloperSecretKey = "test"
	p.TelemetryKey = ""
	err = p.Init()
	assert.NilError(t, err)
	assert.Equal(t, p.EventNamespace, defaultNamespace)

	p.EventNamespace = "telegraf"
	err = p.Init()
	assert.Error(t, err, fmt.Sprintf(errEventNamespace, defaultNamespace))

	p.EventNamespace = "custom_telegraf"
	err = p.Init()
	assert.Error(t, err, fmt.Sprintf(errEventNamespace, defaultNamespace))

	p.EventNamespace = "custom.telegraf"
	err = p.Init()
	assert.NilError(t, err)
}
