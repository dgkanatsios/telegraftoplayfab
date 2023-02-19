//go:generate ../../../tools/readme_config_includer/generator
package playfab_data_insights

import (
	_ "embed"
	"fmt"

	playfab "github.com/dgkanatsios/playfabsdk-go/sdk"
	"github.com/dgkanatsios/playfabsdk-go/sdk/authentication"
	"github.com/dgkanatsios/playfabsdk-go/sdk/events"
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/outputs"
)

//go:embed sample.conf
var sampleConfig string

type PlayFabDataInsights struct {
	Log                telegraf.Logger `toml:"-"`
	TitleId            string          `json:"titleId"`
	DeveloperSecretKey string          `json:"developerSecretKey"`
	EntityToken        string          `json:"entityToken"`
}

func (*PlayFabDataInsights) SampleConfig() string {
	return sampleConfig
}

// Init is for setup, and validating config.
func (p *PlayFabDataInsights) Init() error {
	return nil
}

func (p *PlayFabDataInsights) Connect() error {
	settings := playfab.NewSettingsWithDefaultOptions(p.TitleId)
	postData := &authentication.GetEntityTokenRequestModel{}
	r, err := authentication.GetEntityToken(settings, postData, "", "", p.DeveloperSecretKey)
	if err != nil {
		return err
	}
	p.EntityToken = r.EntityToken
	return nil
}

func (p *PlayFabDataInsights) Write(metrics []telegraf.Metric) error {
	eventsToSend := make([]events.EventContentsModel, 0)
	for _, metric := range metrics {
		for _, field := range metric.FieldList() {
			eventToSend := events.EventContentsModel{
				CustomTags:     metric.Tags(),
				EventNamespace: "custom.telegraf3",
				Name:           fmt.Sprintf("%s_%s", metric.Name(), field.Key),
				Payload:        fmt.Sprintf("%v", field.Value),
			}
			//p.Log.Debugf("%#v\n", eventToSend)
			eventsToSend = append(eventsToSend, eventToSend)
		}
	}
	postData := &events.WriteEventsRequestModel{
		Events: eventsToSend,
	}
	//p.Log.Debugf("%#v\n", postData)
	settings := playfab.NewSettingsWithDefaultOptions(p.TitleId)
	_, err := events.WriteTelemetryEvents(settings, postData, p.EntityToken)
	if err != nil {
		return err
	}
	return nil
}

func (p *PlayFabDataInsights) Close() error {
	// Close any connections here.
	// Write will not be called once Close is called, so there is no need to synchronize.
	return nil
}

func init() {
	outputs.Add("playfab_data_insights", func() telegraf.Output { return &PlayFabDataInsights{} })
}
