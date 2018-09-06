package controllers

import (
	"encoding/json"
	"gAPIManagement/api/config"
	"gAPIManagement/api/database"
	"gAPIManagement/api/http"
	"gAPIManagement/api/servicediscovery"
	"strconv"

	"github.com/qiangxue/fasthttp-routing"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func CreateAppGroup(c *routing.Context) error {
	var bodyMap servicediscovery.ApplicationGroup
	err := json.Unmarshal(c.Request.Body(), &bodyMap)

	if err != nil {
		http.Response(c, err.Error(), 400, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	}

	if bodyMap.Name == "" {
		http.Response(c, `{"error": true, "msg": "Invalid body. Missing body."}`, 400, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	}

	session, db := database.GetSessionAndDB(database.MONGO_DB)
	collection := db.C(servicediscovery.SERVICE_APPS_GROUP_COLLECTION)
	index := mgo.Index{
		Key:        []string{"name"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err = collection.EnsureIndex(index)

	bodyMap.Id = bson.NewObjectId()
	err = collection.Insert(&bodyMap)

	database.MongoDBPool.Close(session)

	if err != nil {
		http.Response(c, `{"error" : true, "msg": `+strconv.Quote(err.Error())+`}`, 400, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	}
	http.Response(c, `{"error" : false, "msg": "Service created successfuly."}`, 201, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
	return nil
}

func GetAppGroups(c *routing.Context) error {
	nameFilter := ""
	if c.QueryArgs().Has("name") {
		nameFilter = string(c.QueryArgs().Peek("name"))
	}

	// Get page
	page := 1
	if c.QueryArgs().Has("page") {
		var err error
		page, err = strconv.Atoi(string(c.QueryArgs().Peek("page")))

		if err != nil {
			http.Response(c, `{"error" : true, "msg": "Invalid page provided."}`, 404, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
			return nil
		}
	}
	skips := servicediscovery.PAGE_LENGTH * (page - 1)

	// Get list of application groups
	var appGroups []servicediscovery.ApplicationGroup

	session, db := database.GetSessionAndDB(database.MONGO_DB)

	db.C(servicediscovery.SERVICE_APPS_GROUP_COLLECTION).Find(bson.M{"name": bson.RegEx{nameFilter + ".*", ""}}).Sort("name").Skip(skips).Limit(servicediscovery.PAGE_LENGTH).All(&appGroups)

	database.MongoDBPool.Close(session)

	if len(appGroups) == 0 {
		http.Response(c, `[]`, 200, servicediscovery.SERVICE_NAME, config.APPLICATION_JSON)
		return nil
	}
	appGroupsString, _ := json.Marshal(appGroups)
	http.Response(c, string(appGroupsString), 200, servicediscovery.SERVICE_NAME, config.APPLICATION_JSON)
	return nil
}

func DeleteAppGroup(c *routing.Context) error {
	appGroupId := c.Param("group_id")
	if !bson.IsObjectIdHex(string(appGroupId)) {
		http.Response(c, `{"error" : true, "msg": "Group id not valid."}`, 400, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	}
	appGroupIdHex := bson.ObjectIdHex(appGroupId)

	session, db := database.GetSessionAndDB(database.MONGO_DB)

	err := db.C(servicediscovery.SERVICE_APPS_GROUP_COLLECTION).RemoveId(appGroupIdHex)

	database.MongoDBPool.Close(session)

	if err != nil {
		http.Response(c, `{"error" : true, "msg": `+strconv.Quote(err.Error())+`}`, 400, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	}

	http.Response(c, `{"error" : false, "msg": "Applications group removed successfuly."}`, 200, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
	return nil
}

func GetAppGroupById(c *routing.Context) error {
	appGroupId := c.Param("group_id")
	if !bson.IsObjectIdHex(string(appGroupId)) {
		http.Response(c, `{"error" : true, "msg": "Group id not valid."}`, 400, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	}
	appGroupIdHex := bson.ObjectIdHex(appGroupId)
	session, db := database.GetSessionAndDB(database.MONGO_DB)

	var group servicediscovery.ApplicationGroup
	err := db.C(servicediscovery.SERVICE_APPS_GROUP_COLLECTION).FindId(appGroupIdHex).One(&group)

	var servicesList []servicediscovery.Service
	err = db.C(servicediscovery.SERVICES_COLLECTION).Find(bson.M{"_id": bson.M{"$in": group.Services}}).All(&servicesList)

	database.MongoDBPool.Close(session)

	if err != nil {
		http.Response(c, `{"error" : true, "msg": `+strconv.Quote(err.Error())+`}`, 400, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	}

	var responseMap map[string]interface{}
	responseMap = make(map[string]interface{})
	responseMap["Name"] = group.Name
	responseMap["Id"] = group.Id
	responseMap["Services"] = servicesList

	if len(servicesList) == 0 {
		responseMap["Services"] = []string{}
	}

	gjson, _ := json.Marshal(responseMap)

	http.Response(c, string(gjson), 200, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
	return nil
}

func UpdateAppGroup(c *routing.Context) error {
	appGroupId := c.Param("group_id")
	if !bson.IsObjectIdHex(string(appGroupId)) {
		http.Response(c, `{"error" : true, "msg": "Group id not valid."}`, 400, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	}
	appGroupIdHex := bson.ObjectIdHex(appGroupId)

	var aGroup servicediscovery.ApplicationGroup
	sgNew := c.Request.Body()
	json.Unmarshal(sgNew, &aGroup)

	session, db := database.GetSessionAndDB(database.MONGO_DB)

	updateGroupQuery := bson.M{"$set": bson.M{"name": aGroup.Name}}
	err := db.C(servicediscovery.SERVICE_APPS_GROUP_COLLECTION).UpdateId(appGroupIdHex, updateGroupQuery)

	database.MongoDBPool.Close(session)

	if err != nil {
		http.Response(c, `{"error" : true, "msg": `+strconv.Quote(err.Error())+`}`, 400, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	}

	http.Response(c, `{"error" : false, "msg": "Applications group updated successfuly."}`, 200, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
	return nil
}

func DeassociateServiceFromApplicationGroup(c *routing.Context) error {
	appGroupId := c.Param("group_id")
	serviceId := c.Param("service_id")
	if !bson.IsObjectIdHex(string(serviceId)) || !bson.IsObjectIdHex(string(appGroupId)) {
		http.Response(c, `{"error" : true, "msg": "Service/Group id not valid."}`, 400, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	}
	serviceGroupIdHex := bson.ObjectIdHex(appGroupId)
	serviceIdHx := bson.ObjectIdHex(serviceId)

	removeFromAllGroups := bson.M{"$pull": bson.M{"services": serviceIdHx}}

	session, db := database.GetSessionAndDB(database.MONGO_DB)

	err := db.C(servicediscovery.SERVICE_APPS_GROUP_COLLECTION).UpdateId(serviceGroupIdHex, removeFromAllGroups)

	if err != nil {
		database.MongoDBPool.Close(session)

		http.Response(c, `{"error" : true, "msg": `+strconv.Quote(err.Error())+`}`, 400, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	}

	database.MongoDBPool.Close(session)

	if err != nil {
		http.Response(c, `{"error" : true, "msg": `+strconv.Quote(err.Error())+`}`, 400, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	}
	http.Response(c, `{"error" : false, "msg": "Service deassociated from group successfuly."}`, 201, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
	return nil
}

func AssociateServiceToAppGroup(c *routing.Context) error {
	appGroupId := c.Param("group_id")
	serviceId := c.Param("service_id")
	if !bson.IsObjectIdHex(string(serviceId)) || !bson.IsObjectIdHex(string(appGroupId)) {
		http.Response(c, `{"error" : true, "msg": "Service/Group id not valid."}`, 400, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	}
	serviceGroupIdHex := bson.ObjectIdHex(appGroupId)
	serviceIdHx := bson.ObjectIdHex(serviceId)

	removeFromAllGroups := bson.M{"$pull": bson.M{"services": serviceIdHx}}
	updateGroup := bson.M{"$addToSet": bson.M{"services": serviceIdHx}}

	session, db := database.GetSessionAndDB(database.MONGO_DB)

	_, err := db.C(servicediscovery.SERVICE_APPS_GROUP_COLLECTION).UpdateAll(bson.M{}, removeFromAllGroups)
	err = db.C(servicediscovery.SERVICE_APPS_GROUP_COLLECTION).UpdateId(serviceGroupIdHex, updateGroup)

	if err != nil {
		database.MongoDBPool.Close(session)

		http.Response(c, `{"error" : true, "msg": `+strconv.Quote(err.Error())+`}`, 400, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	}

	database.MongoDBPool.Close(session)

	if err != nil {
		http.Response(c, `{"error" : true, "msg": `+strconv.Quote(err.Error())+`}`, 400, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	}
	http.Response(c, `{"error" : false, "msg": "Service added to group successfuly."}`, 201, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
	return nil
}

func FindAppGroupForService(c *routing.Context) error {
	serviceId := c.Param("service_id")
	if !bson.IsObjectIdHex(string(serviceId)) {
		http.Response(c, `{"error" : true, "msg": "Service id not valid."}`, 400, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	}
	serviceIdHx := bson.ObjectIdHex(serviceId)

	session, db := database.GetSessionAndDB(database.MONGO_DB)

	var appGroup servicediscovery.ApplicationGroup

	query := bson.M{"services": serviceIdHx}
	db.C(servicediscovery.SERVICE_APPS_GROUP_COLLECTION).Find(query).One(&appGroup)

	database.MongoDBPool.Close(session)

	if appGroup.Name == "" {
		http.Response(c, `{"error" : true, "msg": "Service is not associated to an application group."}`, 404, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	}
	jsonAppGroup, _ := json.Marshal(appGroup)

	http.Response(c, string(jsonAppGroup), 200, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
	return nil
}

func AppGroupsMatches(c *routing.Context) error {
	groupName := string(c.QueryArgs().Peek("group_name"))
	if groupName == "" {
		http.Response(c, `{"error":true, "msg": "Invalid group name"}`, 400, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	}

	session, db := database.GetSessionAndDB(database.MONGO_DB)

	var servicesThatMatch []servicediscovery.Service

	query := []bson.M{
		{"$lookup": bson.M{"from": servicediscovery.SERVICE_APPS_GROUP_COLLECTION, "localField": "_id", "foreignField": "services", "as": "service_app_group"}},
		{"$addFields": bson.M{"zeroAppGroups": bson.M{"$not": bson.M{"$size": "$service_app_group"}}}},
		{"$match": bson.M{"zeroAppGroups": true, "matchinguri": bson.RegEx{"/api/(experience|system|process)/(v1/)?" + groupName + "(/\\w*)?", "i"}}},
	}
	// query := bson.M{"matchinguri": bson.RegEx{"/api/(experience|system|process)/" + groupName + "/\\w+", ""}}

	db.C(servicediscovery.SERVICES_COLLECTION).Pipe(query).All(&servicesThatMatch)

	database.MongoDBPool.Close(session)

	if len(servicesThatMatch) == 0 {
		http.Response(c, `[]`, 200, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	}
	servicesThatMatchJson, _ := json.Marshal(servicesThatMatch)
	http.Response(c, string(servicesThatMatchJson), 200, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
	return nil
}

func UngroupedApps(c *routing.Context) error {
	session, db := database.GetSessionAndDB(database.MONGO_DB)

	var servicesThatMatch []servicediscovery.Service

	query := []bson.M{
		{"$lookup": bson.M{"from": servicediscovery.SERVICE_APPS_GROUP_COLLECTION, "localField": "_id", "foreignField": "services", "as": "service_app_group"}},
		{"$addFields": bson.M{"zeroAppGroups": bson.M{"$not": bson.M{"$size": "$service_app_group"}}}},
		{"$match": bson.M{"zeroAppGroups": true}},
	}

	db.C(servicediscovery.SERVICES_COLLECTION).Pipe(query).All(&servicesThatMatch)

	database.MongoDBPool.Close(session)

	servicesThatMatchJson, _ := json.Marshal(servicesThatMatch)
	http.Response(c, string(servicesThatMatchJson), 200, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
	return nil
}
