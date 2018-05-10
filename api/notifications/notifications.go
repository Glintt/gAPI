package notifications

import (
	"gAPIManagement/api/config"
)

var Methods = map[string]interface{}{
	"Slack": SlackNotification}


func SendNotification(msg string){
	Methods[config.GApiConfiguration.Notifications.Type].(func(string))(msg)
}