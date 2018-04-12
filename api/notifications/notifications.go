package notifications

import (
	"gAPIManagement/api/config"
)

var funcMap = map[string]interface{}{
	"Slack": SlackNotification}


func SendNotification(msg string){
	funcMap[config.GApiConfiguration.Notifications.Type].(func(string))(msg)
}