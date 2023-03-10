//go:generate ../../../tools/readme_config_includer/generator
package playfab

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
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

// PlayFab is the top level struct for this plugin.
type PlayFab struct {
	TitleId            string `toml:"titleId"`
	DeveloperSecretKey string `toml:"developerSecretKey"`
	EventNamespace     string `toml:"eventNamespace"`
	Debug              bool   `toml:"debug"`
	entityToken        string
}

// SampleConfig returns the sample config for this plugin
func (*PlayFab) SampleConfig() string {
	return sampleConfig
}

// Init is for setup, and validating config
func (p *PlayFab) Init() error {
	if p.TitleId == "" {
		return fmt.Errorf("titleId is a required field for playfab output")
	}

	if p.DeveloperSecretKey == "" {
		return fmt.Errorf("developerSecretKey is a required field for playfab output")
	}

	if p.EventNamespace == "" {
		p.EventNamespace = defaultNamespace
	} else if !strings.HasPrefix(p.EventNamespace, fmt.Sprintf("%s.", defaultNamespace)) {
		return fmt.Errorf("eventNamespace must start with %s followed by a dot (.)", defaultNamespace)
	}

	if p.Debug {
		log.Println("Successfully initialized PlayFab output plugin")
	}

	return nil
}

// Connect tries to connect to PlayFab and get an EntityToken
func (p *PlayFab) Connect() error {
	settings := playfab.NewSettingsWithDefaultOptions(p.TitleId)
	postData := &authentication.GetEntityTokenRequestModel{}
	r, err := authentication.GetEntityToken(settings, postData, "", "", p.DeveloperSecretKey)
	if err != nil {
		return err
	}
	p.entityToken = r.EntityToken
	if p.Debug {
		log.Println("Successfully connected to PlayFab")
	}
	return nil
}

// Write writes metrics to PlayFab
func (p *PlayFab) Write(metrics []telegraf.Metric) error {
	eventsToSend := make([]events.EventContentsModel, 0)
	for _, metric := range metrics {
		// marshal the entire payload (fields names and keys into JSON)
		payloadBytes, err := json.Marshal(metric.Fields())
		if err != nil {
			log.Printf("E! Error marshalling metric fields to JSON: %v", metric.Fields())
			continue
		}

		// create the event to send to PlayFab
		eventToSend := events.EventContentsModel{
			CustomTags:        metric.Tags(),
			EventNamespace:    p.EventNamespace,
			Name:              metric.Name(),
			PayloadJSON:       string(payloadBytes),
			OriginalTimestamp: time.Now().UTC(),
		}
		if p.Debug {
			log.Printf("Gathering eventToSend %#v\n", eventToSend)
		}

		eventsToSend = append(eventsToSend, eventToSend)
	}

	postData := &events.WriteEventsRequestModel{
		Events: eventsToSend,
	}

	settings := playfab.NewSettingsWithDefaultOptions(p.TitleId)
	_, err := events.WriteTelemetryEvents(settings, postData, p.entityToken)

	if err != nil {
		return err
	}
	if p.Debug {
		log.Println("Successfully sent events to PlayFab")
	}
	return nil
}

// Close is a no-op for this plugin
func (p *PlayFab) Close() error {
	return nil
}

func init() {
	outputs.Add("playfab", func() telegraf.Output { return &PlayFab{} })
}
