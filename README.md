# PlayFab Insights output plugin for Telegraf

## Description

Azure PlayFab Insights https://learn.microsoft.com/en-us/gaming/playfab/features/insights/

Based on https://github.com/dgkanatsios/playfabsdk-go

## Usage

```toml
[[outputs.execd]]
  command = ["/path/to/telegraftoplayfabinsights", "--config", "/path/to/plugin.conf"]
```