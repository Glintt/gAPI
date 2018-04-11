package healthcheck

import (
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

/* 	timeout := time.Duration(time.Duration(TimeoutDuration) * time.Second) */
	var servicesFinal []servicediscovery.Service

	fmt.Println("##### HEALTH CHECK ##### ")

	for _, s := range services {


		healthcheckURL := s.HealthcheckUrl

		fmt.Println("-----> " + s.Domain + ":" + s.Port + healthcheckURL)
		resp, err := http.Get("http://" + s.Domain+":"+s.Port + healthcheckURL)
		if err != nil || resp.StatusCode != 200 {
			s.IsActive = false
		} else {
			s.IsActive = true
		}

		servicesFinal = append(servicesFinal, s)
	}

	sd.SetRegisteredServices(servicesFinal)
}
