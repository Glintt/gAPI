package controllers

import (
	"encoding/json"
	"gAPIManagement/api/config"
	"gAPIManagement/api/database"
	"gAPIManagement/api/http"
	"gAPIManagement/api/servicediscovery/appgroups"
	"gAPIManagement/api/servicediscovery/constants"
	"gAPIManagement/api/servicediscovery/service"
	"strconv"

	routing "github.com/qiangxue/fasthttp-routing"
	"gopkg.in/mgo.v2/bson"
)

func AppGroupMethods() map[string]interface{} {
	return appgroups.ApplicationGroupMethods[constants.SD_TYPE]
}

func CreateAppGroup(c *routing.Context) error {
	var bodyMap appgroups.ApplicationGroup
	err := json.Unmarshal(c.Request.Body(), &bodyMap)

	if err != nil {
		http.Response(c, err.Error(), 400, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	}

	if bodyMap.Name == "" {
		http.Response(c, `{"error": true, "msg": "Invalid body. Missing body."}`, 400, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	}

	err = AppGroupMethods()["create"].(func(appgroups.ApplicationGroup) error)(bodyMap)

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

	appGroups := AppGroupMethods()["list"].(func(int, string) []appgroups.ApplicationGroup)(page, nameFilter)

	if len(appGroups) == 0 {
		http.Response(c, `[]`, 200, constants.SERVICE_NAME, config.APPLICATION_JSON)
		return nil
	}

	appGroupsString, _ := json.Marshal(appGroups)
	http.Response(c, string(appGroupsString), 200, constants.SERVICE_NAME, config.APPLICATION_JSON)
	return nil
}

func DeleteAppGroup(c *routing.Context) error {
	appGroupId := c.Param("group_id")
	if !bson.IsObjectIdHex(string(appGroupId)) {
		http.Response(c, `{"error" : true, "msg": "Group id not valid."}`, 400, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	}

	err := AppGroupMethods()["delete"].(func(string) error)(appGroupId)

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

	group, err := AppGroupMethods()["getbyid"].(func(string) (appgroups.ApplicationGroup, error))(appGroupId)

	servicesList, err := AppGroupMethods()["getservicesforappgroup"].(func(appgroups.ApplicationGroup) ([]service.Service, error))(group)

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
	var aGroup appgroups.ApplicationGroup
	sgNew := c.Request.Body()
	json.Unmarshal(sgNew, &aGroup)

	err := AppGroupMethods()["update"].(func(string, appgroups.ApplicationGroup) error)(appGroupId, aGroup)

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

	err := AppGroupMethods()["removeappfromgroup"].(func(string, string) error)(appGroupId, serviceId)

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

	err := AppGroupMethods()["addapptogroup"].(func(string, string) error)(appGroupId, serviceId)

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

	appGroup := AppGroupMethods()["getappforservice"].(func(string) appgroups.ApplicationGroup)(serviceId)

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

	var servicesThatMatch []service.Service

	query := []bson.M{
		{"$lookup": bson.M{"from": constants.SERVICE_APPS_GROUP_COLLECTION, "localField": "_id", "foreignField": "services", "as": "service_app_group"}},
		{"$addFields": bson.M{"zeroAppGroups": bson.M{"$not": bson.M{"$size": "$service_app_group"}}}},
		{"$match": bson.M{"zeroAppGroups": true, "matchinguri": bson.RegEx{"/api/(experience|system|process)/(v1/)?" + groupName + "(/\\w*)?", "i"}}},
	}
	// query := bson.M{"matchinguri": bson.RegEx{"/api/(experience|system|process)/" + groupName + "/\\w+", ""}}

	db.C(constants.SERVICES_COLLECTION).Pipe(query).All(&servicesThatMatch)

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
	servicesThatMatch := AppGroupMethods()["ungroupedservices"].(func() []service.Service)()

	servicesThatMatchJson, _ := json.Marshal(servicesThatMatch)
	http.Response(c, string(servicesThatMatchJson), 200, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
	return nil
}
