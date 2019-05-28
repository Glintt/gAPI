package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/Glintt/gAPI/api/config"
	"github.com/Glintt/gAPI/api/http"
	"github.com/Glintt/gAPI/api/plugins"

	"github.com/fatih/structs"

	routing "github.com/qiangxue/fasthttp-routing"
)

func PluginsServiceName() string {
	return plugins.SERVICE_NAME
}

func ListPluginsAvailable(c *routing.Context) error {
	plugins, err := plugins.ListAll()

	if err != nil {
		http.Response(c, fmt.Sprintf(`{"error_msg": "%v"}`, err.Error()), 404, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	}
	pluginsList, _ := json.Marshal(plugins)

	http.Response(c, string(pluginsList), 200, PluginsServiceName(), config.APPLICATION_JSON)
	return nil
}

func ActivePlugins(c *routing.Context) error {
	m := structs.Map(config.GApiConfiguration.Plugins)

	delete(m, "Location")
	activePlugins, _ := json.Marshal(m)

	http.Response(c, string(activePlugins), 200, PluginsServiceName(), config.APPLICATION_JSON)
	return nil
}
