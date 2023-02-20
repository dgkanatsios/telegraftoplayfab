//go:generate ../../../tools/readme_config_includer/generator
package playfab_insights

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	playfab "github.com/dgkanatsios/playfabsdk-go/sdk"
	"github.com/dgkanatsios/playfabsdk-go/sdk/authentication"
	"github.com/dgkanatsios/playfabsdk-go/sdk/events"
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/outputs"
)

const (
	defaultNamespace = "custom"
)

//go:embed sample.conf
var sampleConfig string

// PlayFabInsights is the top level struct for this plugin.
type PlayFabInsights struct {
	Log                telegraf.Logger `toml:"-"`
	TitleId            string          `json:"titleId"`
	DeveloperSecretKey string          `json:"developerSecretKey"`
	EventNamespace     string          `json:"eventNamespace"`
	entityToken        string
}

// SampleConfig returns the sample config for this plugin
func (*PlayFabInsights) SampleConfig() string {
	return sampleConfig
}

// Init is for setup, and validating config
func (p *PlayFabInsights) Init() error {
	if p.TitleId == "" {
		return fmt.Errorf("titleId is a required field for playfab_insights output")
	}

	if p.DeveloperSecretKey == "" {
		return fmt.Errorf("developerSecretKey is a required field for playfab_insights output")
	}

	if p.EventNamespace == "" {
		p.EventNamespace = defaultNamespace
	} else if !strings.HasPrefix(p.EventNamespace, fmt.Sprintf("%s.", defaultNamespace)) {
		return fmt.Errorf("eventNamespace must start with %s followed by a dot (.)", defaultNamespace)
	}

	return nil
}

// Connect tries to connect to PlayFab and get an EntityToken
func (p *PlayFabInsights) Connect() error {
	settings := playfab.NewSettingsWithDefaultOptions(p.TitleId)
	postData := &authentication.GetEntityTokenRequestModel{}
	r, err := authentication.GetEntityToken(settings, postData, "", "", p.DeveloperSecretKey)
	if err != nil {
		return err
	}
	p.entityToken = r.EntityToken
	return nil
}

// Write writes metrics to PlayFab Insights
func (p *PlayFabInsights) Write(metrics []telegraf.Metric) error {
	eventsToSend := make([]events.EventContentsModel, 0)
	for _, metric := range metrics {
		// marshal the entire payload (fields names and keys into JSON)
		payloadBytes, err := json.Marshal(metric.Fields())
		if err != nil {
			return err
		}

		// create the event to send to PlayFab Insights
		eventToSend := events.EventContentsModel{
			CustomTags:        metric.Tags(),
			EventNamespace:    p.EventNamespace,
			Name:              metric.Name(),
			PayloadJSON:       string(payloadBytes),
			OriginalTimestamp: time.Now().UTC(),
		}
		//p.Log.Debugf("Gathering eventToSend %#v\n", eventToSend)
		eventsToSend = append(eventsToSend, eventToSend)
	}
	postData := &events.WriteEventsRequestModel{
		Events: eventsToSend,
	}
	settings := playfab.NewSettingsWithDefaultOptions(p.TitleId)
	_, err := events.WriteTelemetryEvents(settings, postData, p.entityToken)
	return err
}

// Close is a no-op for this plugin
func (p *PlayFabInsights) Close() error {
	// Close any connections here.
	// Write will not be called once Close is called, so there is no need to synchronize.
	return nil
}

func init() {
	outputs.Add("playfab_insights", func() telegraf.Output { return &PlayFabInsights{} })
}
