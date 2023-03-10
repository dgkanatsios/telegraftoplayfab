# Telegraf output plugin for PlayFab

## Description

This is an [external Telegraf plugin](https://github.com/influxdata/telegraf/blob/master/docs/EXTERNAL_PLUGINS.md) that sends events to [PlayFab](https://learn.microsoft.com/en-us/gaming/playfab/). Events can be from anything Telegraf lists as [an input plugin](https://github.com/influxdata/telegraf/tree/master/plugins/inputs). You can easily view the events on PlayFab Game Manager Data Explorer with [basic mode](https://learn.microsoft.com/en-us/gaming/playfab/features/insights/data-explorer/getting-started-with-data-explorer-basic) or [advanced mode](https://learn.microsoft.com/en-us/gaming/playfab/features/insights/data-explorer/getting-started-with-data-explorer-advanced).

Telegraf agent uses the [execd output plugin](https://github.com/influxdata/telegraf/blob/master/plugins/outputs/execd/README.md) to launch and emit metrics to this plugin. The plugin itself uses the [Write Telemetry Events](https://learn.microsoft.com/en-us/rest/api/playfab/events/play-stream-events/write-telemetry-events?view=playfab-rest) API to send custom events to PlayFab.

This plugin uses the [Go PlayFab SDK](https://github.com/dgkanatsios/playfabsdk-go). You need to have a [PlayFab account](https://learn.microsoft.com/en-us/gaming/playfab/gamemanager/pfab-account) and a title to use this plugin.

## Usage

This plugin is meant to be used with Telegraf agent, so you should grab the Telegraf binary for your platform and architecture from the [Telegraf releases page](https://github.com/influxdata/telegraf/releases). You should also download the plugin from this [repo's releases page](https://github.com/dgkanatsios/telegraftoplayfab/releases).

On the Telegraf configuration file `telegraf.conf` you should configure the execd output plugin like this:

```toml
[[outputs.execd]]
  command = ["/path/to/telegraftoplayfab", "--config", "/path/to/plugin.conf"]
```

The `plugin.conf` file should contain the following:

```toml
[[outputs.playfab]]
  titleId = "yourPlayFabTitleId"
  developerSecretKey = "yourDeveloperSecretKey" # https://learn.microsoft.com/en-us/gaming/playfab/gamemanager/secret-key-management
  debug = false # enable it only when debugging as it can be quite verbose
```

For full configuration instructions, check the [sample configuration file](plugins/outputs/playfab/sample.conf).

## Building

To build this plugin, run `make build`.