package healthcheck

import (
	"gAPIManagement/api/utils"
	"gAPIManagement/api/notifications"
	
	"gAPIManagement/api/config"
	"gAPIManagement/api/servicediscovery"
	"net/http"
	"time"
)

var sd *servicediscovery.ServiceDiscovery

var TickerTime = 30
var TimeoutDuration = 2

func InitHealthCheck() {
	if !config.GApiConfiguration.Healthcheck.Active {
		return
	}

	TickerTime = config.GApiConfiguration.Healthcheck.Frequency

	sd = servicediscovery.GetServiceDiscoveryObject()

	ticker := time.NewTicker(time.Duration(TickerTime) * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				CheckServicesHealth()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func ServicesList() []servicediscovery.Service {
	services, _ := sd.GetAllServices()
	return services
}

func CheckServicesHealth() {
	utils.LogMessage("##### HEALTH CHECK #####")

	for _, s := range ServicesList() {
		healthcheckURL := s.HealthcheckUrl
		// healthcheckURL = "http://" + s.Domain + ":" + s.Port + healthcheckURL

		utils.LogMessage("-----> " + healthcheckURL)
		
		go func(healthcheckURL string, s servicediscovery.Service){
			resp, err := http.Get(healthcheckURL)
			if err != nil || resp.StatusCode != 200 {
				NotifyHealthDown(s)
				if s.IsActive == true {
					s.LastActiveTime = utils.CurrentTimeMilliseconds()
				}
				s.IsActive = false
			} else {
				NotifyHealthUp(s)
				s.LastActiveTime = 0
				s.IsActive = true
			}

			sd.UpdateService(s)
		}(healthcheckURL, s)
	}
	utils.LogMessage("### HEALTH CHECK ENDED ###")
}


func NotifyHealthDown(service servicediscovery.Service){
	if ! config.GApiConfiguration.Healthcheck.Notification ||  ! service.IsActive {
		return
	}

	msg := "*" + service.Name + "* located at *" + service.Domain + ":" + service.Port + service.ToURI + "* is down :thinking_face: :thinking_face:"

	notifications.SendNotification(msg)
}

func NotifyHealthUp(service servicediscovery.Service){
	if ! config.GApiConfiguration.Healthcheck.Notification || service.IsActive {
		return
	}

	msg := "*" + service.Name + "* located at *" + service.Domain + ":" + service.Port + service.ToURI + "* went up again! :smiley: :smiley:"

	notifications.SendNotification(msg)
}