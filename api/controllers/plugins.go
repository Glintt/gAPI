package controllers

import (
	"encoding/json"
	"fmt"
	"gAPIManagement/api/config"
	"gAPIManagement/api/http"
	"gAPIManagement/api/plugins"

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
