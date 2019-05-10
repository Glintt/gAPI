package controllers

import (
	"encoding/json"
	"gAPIManagement/api/config"
	"gAPIManagement/api/database"
	"gAPIManagement/api/http"
	"gAPIManagement/api/servicediscovery"
	"gAPIManagement/api/servicediscovery/constants"
	"gAPIManagement/api/servicediscovery/servicegroup"
	"strconv"

	routing "github.com/qiangxue/fasthttp-routing"
	"gopkg.in/mgo.v2/bson"
)

func ListServiceGroupsHandler(c *routing.Context) error {
	sg, err := ServiceDiscovery().GetListOfServicesGroup()

	if err != nil {
		http.Response(c, `{"error" : true, "msg": `+strconv.Quote(err.Error())+`}`, 400, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	}

	json, _ := json.Marshal(sg)
	http.Response(c, string(json), 200, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
	return nil
}

func RegisterServiceGroupHandler(c *routing.Context) error {

	serviceGroup, err := servicediscovery.ValidateServiceGroupBody(c)
	if err != nil {
		http.Response(c, err.Error(), 400, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	}

	serviceGroup.Id = bson.NewObjectId()

	err = servicediscovery.ServiceGroupMethods()["create"].(func(servicegroup.ServiceGroup) error)(serviceGroup)

	if err != nil {
		http.Response(c, `{"error" : true, "msg": `+strconv.Quote(err.Error())+`}`, 400, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	}
	http.Response(c, `{"error" : false, "msg": "Service created successfuly."}`, 201, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
	return nil
}

func AddServiceToGroupHandler(c *routing.Context) error {
	serviceGroupId := c.Param("group_id")

	var bodyMap map[string]string
	err := json.Unmarshal(c.Request.Body(), &bodyMap)

	if err != nil {
		http.Response(c, err.Error(), 400, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	}
	if _, ok := bodyMap["service_id"]; !ok {
		http.Response(c, `{"error": "Invalid body. Missing service_id."}`, 400, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	}
	if serviceGroupId == "null" || bodyMap["service_id"] == "null" || bodyMap["service_id"] == "" {
		http.Response(c, `{"error": "Invalid body."}`, 400, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	}

	err = servicediscovery.ServiceGroupMethods()["addservicetogroup"].(func(string, string) error)(serviceGroupId, bodyMap["service_id"])

	if err != nil {
		http.Response(c, `{"error" : true, "msg": `+strconv.Quote(err.Error())+`}`, 400, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	}
	http.Response(c, `{"error" : false, "msg": "Service added to group successfuly."}`, 201, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
	return nil
}

func DeassociateServiceFromGroup(c *routing.Context) error {
	serviceGroupId := c.Param("group_id")
	serviceId := c.Param("service_id")

	if serviceGroupId == "null" || serviceId == "null" || serviceId == "" {
		http.Response(c, `{"error": "Invalid body."}`, 400, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	}

	serviceGroupIdHex := bson.ObjectIdHex(serviceGroupId)
	serviceIdHex := bson.ObjectIdHex(serviceId)

	updateGroup := bson.M{"$pull": bson.M{"services": serviceIdHex}}
	updateService := bson.M{"$set": bson.M{"groupid": nil, "usegroupattributes": false}}

	session, db := database.GetSessionAndDB(database.MONGO_DB)

	err := db.C(constants.SERVICES_COLLECTION).UpdateId(serviceIdHex, updateService)

	if err != nil {
		database.MongoDBPool.Close(session)

		http.Response(c, `{"error" : true, "msg": `+strconv.Quote(err.Error())+`}`, 400, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	}

	err = db.C(constants.SERVICE_GROUP_COLLECTION).UpdateId(serviceGroupIdHex, updateGroup)

	database.MongoDBPool.Close(session)

	if err != nil {
		http.Response(c, `{"error" : true, "msg": "`+strconv.Quote(err.Error())+`}`, 400, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	}
	http.Response(c, `{"error" : false, "msg": "Service deassociated from group successfuly."}`, 201, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
	return nil
}

func UpdateServiceGroup(c *routing.Context) error {
	serviceGroupId := bson.ObjectIdHex(c.Param("group_id"))

	var sGroup servicegroup.ServiceGroup
	sgNew := c.Request.Body()
	json.Unmarshal(sgNew, &sGroup)

	session, db := database.GetSessionAndDB(database.MONGO_DB)

	err := db.C(constants.SERVICE_GROUP_COLLECTION).UpdateId(serviceGroupId, sGroup)

	if err != nil {
		database.MongoDBPool.Close(session)

		http.Response(c, `{"error" : true, "msg": `+strconv.Quote(err.Error())+`}`, 400, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	}

	database.MongoDBPool.Close(session)
	http.Response(c, `{"error" : false, "msg": "Service group update successfuly."}`, 200, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
	return nil
}

func RemoveServiceGroup(c *routing.Context) error {
	serviceGroupId := bson.ObjectIdHex(c.Param("group_id"))

	session, db := database.GetSessionAndDB(database.MONGO_DB)

	err := db.C(constants.SERVICE_GROUP_COLLECTION).RemoveId(serviceGroupId)

	_, err = db.C(constants.SERVICES_COLLECTION).UpdateAll(bson.M{}, bson.M{"$set": bson.M{"groupid": nil}})

	database.MongoDBPool.Close(session)

	if err != nil {
		http.Response(c, `{"error" : true, "msg": `+strconv.Quote(err.Error())+`}`, 400, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	}

	http.Response(c, `{"error" : false, "msg": "Service group removed successfuly."}`, 200, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
	return nil
}

func GetServiceGroupHandler(c *routing.Context) error {
	serviceGroup := string(c.Param("group"))

	session, db := database.GetSessionAndDB(database.MONGO_DB)

	var sg servicegroup.ServiceGroup
	var err error
	if bson.IsObjectIdHex(serviceGroup) {
		err = db.C(constants.SERVICE_GROUP_COLLECTION).FindId(bson.ObjectIdHex(serviceGroup)).One(&sg)
	} else {
		err = db.C(constants.SERVICE_GROUP_COLLECTION).Find(bson.M{"name": serviceGroup}).One(&sg)
	}

	database.MongoDBPool.Close(session)

	if err != nil {
		http.Response(c, `{"error" : true, "msg": `+strconv.Quote(err.Error())+`}`, 400, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	}

	sgByte, _ := json.Marshal(sg)
	http.Response(c, string(sgByte), 200, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
	return nil
}
