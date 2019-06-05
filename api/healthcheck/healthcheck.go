package healthcheck

import (
	"github.com/Glintt/gAPI/api/notifications"
	"github.com/Glintt/gAPI/api/servicediscovery/service"
	"github.com/Glintt/gAPI/api/utils"

	"net/http"
	"time"

	"github.com/Glintt/gAPI/api/config"
	"github.com/Glintt/gAPI/api/servicediscovery"
)

var sd *servicediscovery.ServiceDiscovery
var services []service.Service

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

func UpdateServicesList() {
	services, _ = sd.GetAllServices()
}

func CheckServicesHealth() {
	utils.LogMessage("##### HEALTH CHECK #####", utils.DebugLogType)

	UpdateServicesList()

	for _, s := range services {
		// healthcheckURL = "http://" + s.Domain + ":" + s.Port + healthcheckURL
		utils.LogMessage("-----> "+s.HealthcheckUrl, utils.DebugLogType)

		go func(healthcheckURL string, s service.Service) {
			resp, err := http.Get(healthcheckURL)
			isActive := s.IsActive

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

			if isActive != s.IsActive {
				service, _ := sd.FindService(s)
				service.IsActive = s.IsActive
				service.LastActiveTime = s.LastActiveTime
				sd.UpdateService(service)
			}

			if resp != nil && resp.Body != nil {
				resp.Body.Close()
			}
			return
		}(s.HealthcheckUrl, s)
	}
	utils.LogMessage("### HEALTH CHECK ENDED ###", utils.DebugLogType)
}

func NotifyHealthDown(s service.Service) {
	if !config.GApiConfiguration.Healthcheck.Notification || !s.IsActive {
		return
	}

	msg := "*" + s.Name + "* located at *" + s.Domain + ":" + s.Port + s.ToURI + "* is down :thinking_face: :thinking_face:"

	notifications.SendNotification(msg)
}

func NotifyHealthUp(s service.Service) {
	if !config.GApiConfiguration.Healthcheck.Notification || s.IsActive {
		return
	}

	msg := "*" + s.Name + "* located at *" + s.Domain + ":" + s.Port + s.ToURI + "* went up again! :smiley: :smiley:"

	notifications.SendNotification(msg)
}
