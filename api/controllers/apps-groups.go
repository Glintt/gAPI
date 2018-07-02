package controllers

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"gopkg.in/mgo.v2"
	"gAPIManagement/api/servicediscovery"
	"gAPIManagement/api/database"
	"encoding/json"
	"gAPIManagement/api/config"
	"github.com/qiangxue/fasthttp-routing"
	"gAPIManagement/api/http"
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
		http.Response(c, `{"error" : true, "msg": ` + strconv.Quote(err.Error()) + `}`, 400, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
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

	session, db := database.GetSessionAndDB(database.MONGO_DB)

	// Get list of application groups
	var appGroups []servicediscovery.ApplicationGroup
	db.C(servicediscovery.SERVICE_APPS_GROUP_COLLECTION).Find(bson.M{"name": bson.RegEx{nameFilter+".*", ""}}).Sort("name").Skip(skips).Limit(servicediscovery.PAGE_LENGTH).All(&appGroups)

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
	appGroupId := bson.ObjectIdHex(c.Param("group_id"))

	session, db := database.GetSessionAndDB(database.MONGO_DB)

	err := db.C(servicediscovery.SERVICE_APPS_GROUP_COLLECTION).RemoveId(appGroupId)

	database.MongoDBPool.Close(session)

	if err != nil {
		http.Response(c, `{"error" : true, "msg": ` + strconv.Quote(err.Error()) + `}`, 400, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	}

	http.Response(c, `{"error" : false, "msg": "Applications group removed successfuly."}`, 200, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
	return nil
}

func GetAppGroupById(c *routing.Context) error {
	appGroupId := bson.ObjectIdHex(c.Param("group_id"))

	session, db := database.GetSessionAndDB(database.MONGO_DB)

	var group servicediscovery.ApplicationGroup
	err := db.C(servicediscovery.SERVICE_APPS_GROUP_COLLECTION).FindId(appGroupId).One(&group)

	database.MongoDBPool.Close(session)

	if err != nil {
		http.Response(c, `{"error" : true, "msg": ` + strconv.Quote(err.Error()) + `}`, 400, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	}

	gjson,_ := json.Marshal(group)
	http.Response(c, string(gjson), 200, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
	return nil
}

func UpdateAppGroup(c *routing.Context) error {	
	appGroupId := bson.ObjectIdHex(c.Param("group_id"))

	var aGroup servicediscovery.ApplicationGroup
	sgNew := c.Request.Body()
	json.Unmarshal(sgNew, &aGroup)

	session, db := database.GetSessionAndDB(database.MONGO_DB)

	fmt.Println(appGroupId.Hex())
	updateGroupQuery := bson.M{"$set": bson.M{"name": aGroup.Name }}
	err := db.C(servicediscovery.SERVICE_APPS_GROUP_COLLECTION).UpdateId(appGroupId, updateGroupQuery)

	database.MongoDBPool.Close(session)

	if err != nil {
		http.Response(c, `{"error" : true, "msg": ` + strconv.Quote(err.Error()) + `}`, 400, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	}

	http.Response(c, `{"error" : false, "msg": "Applications group updated successfuly."}`, 200, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
	return nil
}

func DeassociateServiceFromApplicationGroup(c *routing.Context) error {	
	appGroupId := c.Param("group_id")
	serviceId := c.Param("service_id")

	serviceGroupIdHex := bson.ObjectIdHex(appGroupId)
	serviceIdHx := bson.ObjectIdHex(serviceId)

	removeFromAllGroups := bson.M{"$pull": bson.M{"services": serviceIdHx }}
	
	session, db := database.GetSessionAndDB(database.MONGO_DB)
	
	err := db.C(servicediscovery.SERVICE_APPS_GROUP_COLLECTION).UpdateId(serviceGroupIdHex, removeFromAllGroups)

	if err != nil {
		database.MongoDBPool.Close(session)

		http.Response(c, `{"error" : true, "msg": ` + strconv.Quote(err.Error()) + `}`, 400, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	}

	database.MongoDBPool.Close(session)

	if err != nil {
		http.Response(c, `{"error" : true, "msg": ` + strconv.Quote(err.Error()) + `}`, 400, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	}
	http.Response(c, `{"error" : false, "msg": "Service deassociated from group successfuly."}`, 201, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
	return nil
}

func AssociateServiceToAppGroup(c *routing.Context) error {	
	appGroupId := c.Param("group_id")
	serviceId := c.Param("service_id")
	
	serviceGroupIdHex := bson.ObjectIdHex(appGroupId)
	serviceIdHx := bson.ObjectIdHex(serviceId)

	removeFromAllGroups := bson.M{"$pull": bson.M{"services": serviceIdHx }}
	updateGroup := bson.M{"$addToSet": bson.M{"services": serviceIdHx }}

	session, db := database.GetSessionAndDB(database.MONGO_DB)
	
	_,err := db.C(servicediscovery.SERVICE_APPS_GROUP_COLLECTION).UpdateAll(bson.M{}, removeFromAllGroups)
	err = db.C(servicediscovery.SERVICE_APPS_GROUP_COLLECTION).UpdateId(serviceGroupIdHex, updateGroup)

	if err != nil {
		database.MongoDBPool.Close(session)

		http.Response(c, `{"error" : true, "msg": ` + strconv.Quote(err.Error()) + `}`, 400, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	}

	database.MongoDBPool.Close(session)

	if err != nil {
		http.Response(c, `{"error" : true, "msg": ` + strconv.Quote(err.Error()) + `}`, 400, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	}
	http.Response(c, `{"error" : false, "msg": "Service added to group successfuly."}`, 201, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
	return nil
}
