package notifications

import (
	"github.com/Glintt/gAPI/api/config"
)

var Methods = map[string]interface{}{
	"Slack": SlackNotification}


func SendNotification(msg string){
	Methods[config.GApiConfiguration.Notifications.Type].(func(string))(msg)
}