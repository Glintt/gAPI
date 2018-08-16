package plugins

import (
	"path/filepath"

	routing "github.com/qiangxue/fasthttp-routing"
)

type BeforeRequestPlugin struct {
	Location string
	Filename string
}

func (p *BeforeRequestPlugin) Call(c *routing.Context) error {
	generalPlugin := GeneralPlugin{Location: p.Location, Filename: p.Filename}
	symReqEntryPlugin, err := generalPlugin.Load(BeforeRequestEntryPluginLookup)

	if err != nil {
		return err
	}

	// check if loaded symbol is of desired type
	var pluginBefore BeforeRequestEntryPlugin
	pluginBefore, ok := symReqEntryPlugin.(BeforeRequestEntryPlugin)
	if !ok {
		return err
	}

	// call the plugin
	err = pluginBefore.Call(c)

	if err != nil {
		return err
	}

	return nil
}

func CallBeforeRequestPlugins(c *routing.Context) error {
	allPlugins, err := ListAll()

	if err != nil {
		return nil
	}

	pluginsLocation := filepath.Join(PLUGINS_LOCATION, BEFORE_REQUEST_PLUGINS_NAME)

	for _, pluginToLoad := range allPlugins[BEFORE_REQUEST_PLUGINS_NAME] {
		pluginToCall := BeforeRequestPlugin{Location: pluginsLocation, Filename: pluginToLoad}

		hasStoppingError := pluginToCall.Call(c)

		if hasStoppingError != nil {
			return hasStoppingError
		}
	}

	return nil
}
