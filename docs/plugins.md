# Plugin

Plugins are supported on gAPI.

These plugins must be writen in Golang and are only supportend on Unix.

Currently, it is supported these types of plugins:

1. BeforeRequestPlugin

## Build a plugin

1. A plugin must use _package main_.

2. To build the plugin for usage, run this command:

```
go build -buildmode=plugin -o <plugin_name>.so <plugin_name>.go
```

## Types

#### BeforeRequestPlugin

1. These plugins are run before each request to a microservice on gAPI.

1. When a plugin returns an error, the request will fail and stop.

1. The following methods must be implemented:

```
func (g beforeRequestEntryPlugin) Call(c *routing.Context) error
```

## Activate a plugin on gAPI

Currently, to activate a plugin on gAPI you must add the following configuration to _gAPI.json_ config:

```
"Plugins": {
    "Location": "./custom_plugins",
    "BeforeRequest": ["BeforeRequestPlugin"]
}
```

Params:

1. Location - where plugins are stored
2. BeforeRequest - List of BeforeRequestPlugin type plugins
