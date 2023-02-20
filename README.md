# PlayFab Insights output plugin for Telegraf

## Description

This is an external plugin for Telegraf that sends metrics to [PlayFab Insights](https://learn.microsoft.com/en-us/gaming/playfab/features/insights/). The telegraf agent uses the [execd output plugin](https://github.com/influxdata/telegraf/blob/master/plugins/outputs/execd/README.md) to launch this that uses the [Write Telemetry Events](https://learn.microsoft.com/en-us/rest/api/playfab/events/play-stream-events/write-telemetry-events?view=playfab-rest) API to send custom events to PlayFab Insights.

This plugin uses the [Go PlayFab SDK](https://github.com/dgkanatsios/playfabsdk-go)

## Usage

On `telegraf.conf` you should configure the execd output plugin in this way:

```toml
[[outputs.execd]]
  command = ["/path/to/telegraftoplayfabinsights", "--config", "/path/to/plugin.conf"]
```

The `plugin.conf` file should contain the following:

```toml
[[outputs.playfab_insights]]
  titleId = "yourPlayFabTitleId"
  developerSecretKey = "yourDeveloperSecretKey" # https://learn.microsoft.com/en-us/gaming/playfab/gamemanager/secret-key-management
```

For full configuration instructions, check the [sample configuration file](plugins/outputs/playfab_insights/sample.conf).