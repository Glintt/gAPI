package servicediscovery

import (
	"gopkg.in/mgo.v2/bson"
)


type ApplicationGroup struct {
	Id                         bson.ObjectId `bson:"_id" json:"Id"`
	Name                       string
	Services					[]bson.ObjectId
}