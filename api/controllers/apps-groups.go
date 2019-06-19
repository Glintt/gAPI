package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/Glintt/gAPI/api/config"
	"github.com/Glintt/gAPI/api/database"
	"github.com/Glintt/gAPI/api/http"
	"github.com/Glintt/gAPI/api/servicediscovery/appgroups"
	"github.com/Glintt/gAPI/api/servicediscovery/constants"
	"github.com/Glintt/gAPI/api/servicediscovery/service"

	routing "github.com/qiangxue/fasthttp-routing"
	"gopkg.in/mgo.v2/bson"
)

func AppGroupMethods() map[string]interface{} {
	return appgroups.ApplicationGroupMethods[database.SD_TYPE]
}

func CreateAppGroup(c *routing.Context) error {
	var bodyMap appgroups.ApplicationGroup
	err := json.Unmarshal(c.Request.Body(), &bodyMap)

	if err != nil {
		return http.Error(c, err.Error(), 400, ServiceDiscoveryServiceName())
	}

	if bodyMap.Name == "" {
		return http.Error(c, `Invalid body. Missing body parameter`, 400, ServiceDiscoveryServiceName())
	}

	err = AppGroupMethods()["create"].(func(appgroups.ApplicationGroup) error)(bodyMap)

	if err != nil {
		return http.Error(c, strconv.Quote(err.Error()), 400, ServiceDiscoveryServiceName())
	}
	return http.Created(c, `Service created successfuly`, ServiceDiscoveryServiceName())
}

func GetAppGroups(c *routing.Context) error {
	nameFilter := ""
	if c.QueryArgs().Has("name") {
		nameFilter = string(c.QueryArgs().Peek("name"))
	}

	// Get page
	page, err := http.ParsePageParam(c)
	if err != nil {
		return http.Error(c, err.Error(), 400, ServiceDiscoveryServiceName())
	}

	appGroups := AppGroupMethods()["list"].(func(int, string) []appgroups.ApplicationGroup)(page, nameFilter)

	if len(appGroups) == 0 {
		return http.Ok(c, `[]`, constants.SERVICE_NAME)
	}

	appGroupsString, _ := json.Marshal(appGroups)
	return http.Ok(c, string(appGroupsString), constants.SERVICE_NAME)
}

// DeleteAppGroup handles DELETE /app-groups/<group_id> api endpoint
func DeleteAppGroup(c *routing.Context) error {
	appGroupID := c.Param("group_id")
	if !bson.IsObjectIdHex(string(appGroupID)) {
		return http.Error(c, "Group id not valid.", 400, ServiceDiscoveryServiceName())
	}

	err := AppGroupMethods()["delete"].(func(string) error)(appGroupID)

	if err != nil {
		http.Response(c, `{"error" : true, "msg": `+strconv.Quote(err.Error())+`}`, 400, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	}

	http.Response(c, `{"error" : false, "msg": "Applications group removed successfuly."}`, 200, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
	return nil
}

// GetAppGroupByID handles GET /app-groups/<group_id> api endpoint
func GetAppGroupByID(c *routing.Context) error {
	appGroupID := c.Param("group_id")
	if !bson.IsObjectIdHex(string(appGroupID)) {
		return http.Error(c, "Group id not valid.", 400, ServiceDiscoveryServiceName())
	}

	// Get application group by id
	group, err := AppGroupMethods()["getbyid"].(func(string) (appgroups.ApplicationGroup, error))(appGroupID)
	if err != nil {
		return http.Error(c, strconv.Quote(err.Error()), 400, ServiceDiscoveryServiceName())
	}

	// Get all services associated to applicaiton group
	servicesList, err := AppGroupMethods()["getservicesforappgroup"].(func(appgroups.ApplicationGroup) ([]service.Service, error))(group)
	if err != nil {
		return http.Error(c, strconv.Quote(err.Error()), 400, ServiceDiscoveryServiceName())
	}

	//var responseMap map[string]interface{}
	responseMap := make(map[string]interface{})
	responseMap["Name"] = group.Name
	responseMap["Id"] = group.Id
	responseMap["Services"] = servicesList
	if len(servicesList) == 0 {
		responseMap["Services"] = []string{}
	}

	gjson, _ := json.Marshal(responseMap)

	return http.Ok(c, string(gjson), ServiceDiscoveryServiceName())
}

// UpdateAppGroup handles PUT /app-groups/<group_id> api endpoint
func UpdateAppGroup(c *routing.Context) error {
	appGroupID := c.Param("group_id")
	if !bson.IsObjectIdHex(string(appGroupID)) {
		return http.Error(c, "Group id not valid.", 400, ServiceDiscoveryServiceName())
	}

	var aGroup appgroups.ApplicationGroup
	json.Unmarshal(c.Request.Body(), &aGroup)

	err := AppGroupMethods()["update"].(func(string, appgroups.ApplicationGroup) error)(appGroupID, aGroup)

	if err != nil {
		return http.Error(c, strconv.Quote(err.Error()), 400, ServiceDiscoveryServiceName())
	}

	return http.OkFormated(c, "Applications group updated successfuly.", ServiceDiscoveryServiceName())
}

// DeassociateServiceFromApplicationGroup handles DELETE /apps-groups/<group_id>/<service_id> api endpoint
func DeassociateServiceFromApplicationGroup(c *routing.Context) error {
	appGroupID := c.Param("group_id")
	serviceID := c.Param("service_id")
	if !bson.IsObjectIdHex(string(serviceID)) || !bson.IsObjectIdHex(string(appGroupID)) {
		return http.Error(c, "Service/Group id not valid.", 400, ServiceDiscoveryServiceName())
	}

	err := AppGroupMethods()["removeappfromgroup"].(func(string, string) error)(appGroupID, serviceID)

	if err != nil {
		return http.Error(c, strconv.Quote(err.Error()), 400, ServiceDiscoveryServiceName())
	}

	return http.OkFormated(c, "Service deassociated from group successfuly.", ServiceDiscoveryServiceName())
}

// AssociateServiceToAppGroup handles POST /apps-groups/<group_id>/<service_id> api endpoint
func AssociateServiceToAppGroup(c *routing.Context) error {
	appGroupID := c.Param("group_id")
	serviceID := c.Param("service_id")
	if !bson.IsObjectIdHex(string(serviceID)) || !bson.IsObjectIdHex(string(appGroupID)) {
		return http.Error(c, "Service/Group id not valid.", 400, ServiceDiscoveryServiceName())
	}

	err := AppGroupMethods()["addapptogroup"].(func(string, string) error)(appGroupID, serviceID)

	if err != nil {
		return http.Error(c, strconv.Quote(err.Error()), 400, ServiceDiscoveryServiceName())
	}
	return http.OkFormated(c, "Service added to group successfuly.", ServiceDiscoveryServiceName())
}

// FindAppGroupForService handles GET /apps-groups/search/<service_id> api endpoint. Finds the group for the given service id
func FindAppGroupForService(c *routing.Context) error {
	serviceID := c.Param("service_id")
	if !bson.IsObjectIdHex(string(serviceID)) {
		return http.Error(c, "Service id not valid.", 400, ServiceDiscoveryServiceName())
	}

	appGroup := AppGroupMethods()["getappforservice"].(func(string) appgroups.ApplicationGroup)(serviceID)

	if appGroup.Name == "" {
		return http.Error(c, "Service is not associated to an application group.", 400, ServiceDiscoveryServiceName())
	}
	jsonAppGroup, _ := json.Marshal(appGroup)

	return http.Ok(c, string(jsonAppGroup), ServiceDiscoveryServiceName())
}

// AppGroupsMatches handles GET /apps-groups/matches api endpoint. Finds possible matches for the group. Only available on mongo
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
	servicesThatMatchJSON, _ := json.Marshal(servicesThatMatch)
	return http.Ok(c, string(servicesThatMatchJSON), ServiceDiscoveryServiceName())
}

// UngroupedApps handles GET /apps-groups/ungrouped api endpoint. Gets all services that don't have an application group
func UngroupedApps(c *routing.Context) error {
	servicesThatMatch := AppGroupMethods()["ungroupedservices"].(func() []service.Service)()

	servicesThatMatchJSON, _ := json.Marshal(servicesThatMatch)
	return http.Ok(c, string(servicesThatMatchJSON), ServiceDiscoveryServiceName())
}
