package controllers

import (
	"gAPIManagement/api/servicediscovery"
	"gAPIManagement/api/database"
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"gAPIManagement/api/http"
	"github.com/qiangxue/fasthttp-routing"
)

func ListServiceGroupsHandler(c *routing.Context) error {
	sg, err := ServiceDiscovery().GetListOfServicesGroup()

	if err != nil {
		http.Response(c, `{"error" : true, "msg": "` + err.Error() + `"}`, 400, ServiceDiscoveryServiceName())
		return nil
	}

	json, _ := json.Marshal(sg)
	http.Response(c, string(json), 200, ServiceDiscoveryServiceName())
	return nil
}

func RegisterServiceGroupHandler(c *routing.Context) error {
	session, db := database.GetSessionAndDB(database.MONGO_DB)

	servicegroup, err := servicediscovery.ValidateServiceGroupBody(c)
	if err != nil {
		http.Response(c, err.Error(), 400, ServiceDiscoveryServiceName())
		return nil
	}

	servicegroup.Id = bson.NewObjectId()

	err = db.C(servicediscovery.SERVICE_GROUP_COLLECTION).Insert(&servicegroup)

	database.MongoDBPool.Close(session)

	if err != nil {
		http.Response(c, `{"error" : true, "msg": "` + err.Error() + `"}`, 400, ServiceDiscoveryServiceName())
		return nil
	}
	http.Response(c, `{"error" : false, "msg": "Service created successfuly."}`, 201, ServiceDiscoveryServiceName())
	return nil
}