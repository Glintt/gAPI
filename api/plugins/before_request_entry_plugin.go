package plugins

import (
	"plugin"

	routing "github.com/qiangxue/fasthttp-routing"
)

type BeforeRequestEntryPlugin interface {
	Call(c *routing.Context) error
}

func CallBeforeRequestEntryPlugins(c *routing.Context) error {
	var hasStoppingError error

	configuration := Configurations()

	for _, pluginToLoad := range configuration.BeforeRequest {
		mod := configuration.Location + BEFORE_REQUEST_PLUGINS_INNER_FOLDER_NAME + "/" + pluginToLoad + PLUGIN_EXTENSION
		plug, err := plugin.Open(mod)
		if err != nil {
			continue
		}

		// load symbol
		symReqEntryPlugin, err := plug.Lookup(BeforeRequestEntryPluginLookup)
		if err != nil {
			continue
		}

		// check if loaded symbol is of desired type
		var pluginBefore BeforeRequestEntryPlugin
		pluginBefore, ok := symReqEntryPlugin.(BeforeRequestEntryPlugin)
		if !ok {
			continue
		}

		// call the plugin
		err = pluginBefore.Call(c)

		if err != nil {
			hasStoppingError = err
			break
		}
	}

	return hasStoppingError
}
