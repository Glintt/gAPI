package notifications

import (
	"gAPIManagement/api/config"
	"bytes"
	"net/http"
)


func SlackNotification(msg string){
	json := []byte(`{"text":"` + msg + `"}`)

	webHookApi := config.GApiConfiguration.Notifications.Slack.WebhookUrl
	
    req, _ := http.NewRequest("POST", webHookApi, bytes.NewBuffer(json))
	req.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{}
	
    resp, _ := client.Do(req)
	
	defer resp.Body.Close()
}