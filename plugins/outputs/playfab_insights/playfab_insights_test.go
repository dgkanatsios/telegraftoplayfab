package playfab_insights

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestInit(t *testing.T) {
	p := &PlayFabInsights{}
	err := p.Init()
	assert.Error(t, err, "titleId is a required field for playfab_insights output")

	p.TitleId = "test"
	err = p.Init()
	assert.Error(t, err, "developerSecretKey is a required field for playfab_insights output")

	p.DeveloperSecretKey = "test"
	err = p.Init()
	assert.NilError(t, err)
	assert.Equal(t, p.EventNamespace, defaultNamespace)

	p.EventNamespace = "telegraf"
	err = p.Init()
	assert.Error(t, err, "eventNamespace must start with custom followed by a dot (.)")

	p.EventNamespace = "custom_telegraf"
	err = p.Init()
	assert.Error(t, err, "eventNamespace must start with custom followed by a dot (.)")

	p.EventNamespace = "custom.telegraf"
	err = p.Init()
	assert.NilError(t, err)
}
