package servicediscovery

import (
	"gopkg.in/mgo.v2/bson"
)


type ServiceGroup struct {
	Id                    bson.ObjectId `bson:"_id" json:"Id"`
	Name string
	IsReachable bool
	HostsEnabled []string
	Services []bson.ObjectId
}

func (sg *ServiceGroup) Contains(s Service) bool {
	for _, v := range sg.Services {
		if s.Id == v {
			return true
		}
	}
	return false
}