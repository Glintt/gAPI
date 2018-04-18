package healthcheck

import (
	"gAPIManagement/api/utils"
	"gAPIManagement/api/notifications"
	"fmt"
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
	services := ServicesList()

	var servicesFinal []servicediscovery.Service

	fmt.Println("##### HEALTH CHECK ##### ")

	for _, s := range services {
		healthcheckURL := s.HealthcheckUrl

		fmt.Println("-----> " + s.Domain + ":" + s.Port + healthcheckURL)
		resp, err := http.Get("http://" + s.Domain+":"+s.Port + healthcheckURL)
		if err != nil || resp.StatusCode != 200 {
			NotifyHealthDown(s)
			if s.IsActive == true {
				s.LastActiveTime = utils.CurrentTimeMilliseconds()
			}
			s.IsActive = false
		} else {
			s.LastActiveTime = 0
			s.IsActive = true
		}

		servicesFinal = append(servicesFinal, s)
	}

	sd.SetRegisteredServices(servicesFinal)
}


func NotifyHealthDown(service servicediscovery.Service){
	if ! config.GApiConfiguration.Healthcheck.Notification ||  ! service.IsActive {
		return
	}

	msg := service.Name + " located at " + service.Domain + ":" + service.Port + service.ToURI + " is down!"

	notifications.SendNotification(msg)
}