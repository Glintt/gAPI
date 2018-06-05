package servicediscovery

import (
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"gAPIManagement/api/http"
	"github.com/qiangxue/fasthttp-routing"
)


func (sd *ServiceDiscovery) GetListOfServicesGroup() ([]ServiceGroup, error) {
	session, db := GetSessionAndDB(MONGO_DB)

	var servicesGroup []ServiceGroup
	err := db.C(SERVICE_GROUP_COLLECTION).Find(bson.M{}).All(&servicesGroup)

	mongoPool.Close(session)
	
	return servicesGroup, err
}

func ListServiceGroupsHandler(c *routing.Context) error {
	sg, err := sd.GetListOfServicesGroup()

	if err != nil {
		http.Response(c, `{"error" : true, "msg": "` + err.Error() + `"}`, 400, SERVICE_NAME)
		return nil
	}

	json, _ := json.Marshal(sg)
	http.Response(c, string(json), 200, SERVICE_NAME)
	return nil
}

func RegisterServiceGroupHandler(c *routing.Context) error {
	session, db := GetSessionAndDB(MONGO_DB)

	servicegroup, err := ValidateServiceGroupBody(c)
	if err != nil {
		http.Response(c, err.Error(), 400, SERVICE_NAME)
		return nil
	}

	servicegroup.Id = bson.NewObjectId()

	err = db.C(SERVICE_GROUP_COLLECTION).Insert(&servicegroup)

	mongoPool.Close(session)

	if err != nil {
		http.Response(c, `{"error" : true, "msg": "` + err.Error() + `"}`, 400, SERVICE_NAME)
		return nil
	}
	http.Response(c, `{"error" : false, "msg": "Service created successfuly."}`, 201, SERVICE_NAME)
	return nil
}