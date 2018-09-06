package plugins

import (
	"errors"
	"gAPIManagement/api/config"
	"gAPIManagement/api/http"
	"gAPIManagement/api/utils"
	"path/filepath"

	"github.com/manucorporat/try"
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
		utils.LogMessage("BeforeRequestPlugin Call() : "+err.Error(), utils.DebugLogType)
		return err
	}

	// check if loaded symbol is of desired type
	var pluginBefore BeforeRequestEntryPlugin
	pluginBefore, ok := symReqEntryPlugin.(BeforeRequestEntryPlugin)
	if !ok {
		utils.LogMessage("BeforeRequestPlugin Call() : "+err.Error(), utils.DebugLogType)
		return err
	}

	// call the plugin
	err = pluginBefore.Call(c)

	if err != nil {
		utils.LogMessage("BeforeRequestPlugin Call() : "+err.Error(), utils.DebugLogType)
		return err
	}

	return nil
}

func CallBeforeRequestPlugins(c *routing.Context) error {
	allPlugins := Configurations().BeforeRequest

	pluginsLocation := filepath.Join(PLUGINS_LOCATION, BEFORE_REQUEST_PLUGINS_NAME)

	utils.LogMessage("CallBeforeRequestPlugins() - pluginsLocation: "+pluginsLocation, utils.DebugLogType)

	for _, pluginToLoad := range allPlugins {
		pluginToCall := BeforeRequestPlugin{Location: pluginsLocation, Filename: pluginToLoad}

		utils.LogMessage("CallBeforeRequestPlugins() - pluginToLoad: "+pluginToLoad, utils.DebugLogType)
		var hasStoppingError error
		try.This(func() {
			hasStoppingError = pluginToCall.Call(c)
		}).Catch(func(e try.E) {
			hasStoppingError = errors.New("Plugin " + pluginToLoad + " failed.")
			utils.LogMessage("CallBeforeRequestPlugins() - pluginCall: error on plugin", utils.DebugLogType)
			http.Response(c, `{"error": true, "msg": "`+hasStoppingError.Error()+`"}`, 500, SERVICE_NAME, config.APPLICATION_JSON)
		})

		if hasStoppingError != nil {
			return hasStoppingError
		}
	}

	return nil
}
