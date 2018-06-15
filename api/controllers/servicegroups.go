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

func AddServiceToGroupHandler(c *routing.Context) error {	
	serviceGroupId := c.Param("group_id")
	
	var bodyMap map[string]string
	err := json.Unmarshal(c.Request.Body(), &bodyMap)

	if err != nil {
		http.Response(c, err.Error(), 400, ServiceDiscoveryServiceName())
		return nil
	}
	if _, ok := bodyMap["service_id"]; !ok {
		http.Response(c, `{"error": "Invalid body. Missing service_id."}`, 400, ServiceDiscoveryServiceName())
		return nil
	}
	if serviceGroupId == "null" || bodyMap["service_id"] == "null" || bodyMap["service_id"] == "" {
		http.Response(c, `{"error": "Invalid body."}`, 400, ServiceDiscoveryServiceName())
		return nil
	}

	serviceGroupIdHex := bson.ObjectIdHex(serviceGroupId)
	serviceId := bson.ObjectIdHex(bodyMap["service_id"])

	updateGroup := bson.M{"$addToSet": bson.M{"services": serviceId }}
	updateService := bson.M{"$set": bson.M{"group_id": serviceGroupIdHex}}
	
	session, db := database.GetSessionAndDB(database.MONGO_DB)
	
	err = db.C(servicediscovery.SERVICES_COLLECTION).UpdateId(serviceId, updateService)
	if err != nil {
		http.Response(c, `{"error" : true, "msg": "` + err.Error() + `"}`, 400, ServiceDiscoveryServiceName())
		return nil
	}

	err = db.C(servicediscovery.SERVICE_GROUP_COLLECTION).UpdateId(serviceGroupIdHex, updateGroup)

	database.MongoDBPool.Close(session)

	if err != nil {
		http.Response(c, `{"error" : true, "msg": "` + err.Error() + `"}`, 400, ServiceDiscoveryServiceName())
		return nil
	}
	http.Response(c, `{"error" : false, "msg": "Service added to group successfuly."}`, 201, ServiceDiscoveryServiceName())
	return nil
}