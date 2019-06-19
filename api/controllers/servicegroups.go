package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/Glintt/gAPI/api/config"
	"github.com/Glintt/gAPI/api/http"
	"github.com/Glintt/gAPI/api/servicediscovery"
	"github.com/Glintt/gAPI/api/servicediscovery/servicegroup"

	routing "github.com/qiangxue/fasthttp-routing"
)

func getServiceGroupsService(c *routing.Context) (servicegroup.ServiceGroupService, error) {
	return servicegroup.NewServiceGroupService(c)
}

// ListServiceGroupsHandler handle GET /service-groups
func ListServiceGroupsHandler(c *routing.Context) error {
	serviceGroupService, _ := getServiceGroupsService(c)

	sg, err := serviceGroupService.GetServiceGroups()
	if err != nil {
		return http.Error(c, strconv.Quote(err.Error()), 400, ServiceDiscoveryServiceName())
	}

	json, _ := json.Marshal(sg)
	return http.Ok(c, string(json), ServiceDiscoveryServiceName())
}

// RegisterServiceGroupHandler handle POST /service-groups
func RegisterServiceGroupHandler(c *routing.Context) error {
	// Validate post body
	serviceGroup, err := servicediscovery.ValidateServiceGroupBody(c)
	if err != nil {
		return http.Error(c, strconv.Quote(err.Error()), 400, ServiceDiscoveryServiceName())
	}
	
	// get service group service
	serviceGroupService, _ := getServiceGroupsService(c)

	// Create service group
	err = serviceGroupService.CreateServiceGroup(serviceGroup)
	if err != nil {
		return http.Error(c, strconv.Quote(err.Error()), 400, ServiceDiscoveryServiceName())
	}
	return http.Created(c, "Service created successfuly", ServiceDiscoveryServiceName())
}

// AddServiceToGroupHandler handle POST /service-groups
func AddServiceToGroupHandler(c *routing.Context) error {
	serviceGroupID := c.Param("group_id")

	// Try to parse body
	var bodyMap map[string]string
	err := json.Unmarshal(c.Request.Body(), &bodyMap)
	if err != nil {
		return http.Error(c, strconv.Quote(err.Error()), 400, ServiceDiscoveryServiceName())
	}

	// Validate body
	if _, ok := bodyMap["service_id"]; !ok {
		return http.Error(c, "Invalid body. Missing service_id", 400, ServiceDiscoveryServiceName())
	}
	if serviceGroupID == "null" || bodyMap["service_id"] == "null" || bodyMap["service_id"] == "" {
		return http.Error(c, "Invalid body", 400, ServiceDiscoveryServiceName())
	}

	// Create service group service
	serviceGroupService, _ := getServiceGroupsService(c)

	err = serviceGroupService.AddServiceToGroup(serviceGroupID, bodyMap["service_id"])
	if err != nil {
		return http.Error(c, strconv.Quote(err.Error()), 400, ServiceDiscoveryServiceName())
	}
	return http.Created(c, "Service added to group successfuly", ServiceDiscoveryServiceName())
}

// DeassociateServiceFromGroup handle DELETE /
func DeassociateServiceFromGroup(c *routing.Context) error {
	serviceGroupID := c.Param("group_id")
	serviceID := c.Param("service_id")

	if serviceGroupID == "null" || serviceID == "null" || serviceID == "" {
		http.Response(c, `{"error": "Invalid body."}`, 400, ServiceDiscoveryServiceName(), config.APPLICATION_JSON)
		return nil
	}

	// Create service group service
	serviceGroupService, _ := getServiceGroupsService(c)
	
	err := serviceGroupService.RemoveServiceFromGroup(serviceGroupID, serviceID)

	if err != nil {
		return http.Error(c, strconv.Quote(err.Error()), 400, ServiceDiscoveryServiceName())
	}
	return http.Created(c, "Service deassociated from group successfuly", ServiceDiscoveryServiceName())
}

// UpdateServicecGroup handle PUT /
func UpdateServiceGroup(c *routing.Context) error {
	serviceGroupID := c.Param("group_id")

	var sGroup servicegroup.ServiceGroup
	sgNew := c.Request.Body()
	json.Unmarshal(sgNew, &sGroup)

	// Create service group service
	serviceGroupService, _ := getServiceGroupsService(c)
	
	err := serviceGroupService.UpdateServiceGroup(serviceGroupID, sGroup)
	if err != nil {
		return http.Error(c, strconv.Quote(err.Error()), 400, ServiceDiscoveryServiceName())
	}
	return http.Ok(c, "Service group updated successfuly", ServiceDiscoveryServiceName())
}

// RemoveServiceGroup handle DELETE /
func RemoveServiceGroup(c *routing.Context) error {
	serviceGroupID := c.Param("group_id")

	// Create service group service
	serviceGroupService, _ := getServiceGroupsService(c)
	
	err := serviceGroupService.DeleteServiceGroup(serviceGroupID)
	if err != nil {
		return http.Error(c, strconv.Quote(err.Error()), 400, ServiceDiscoveryServiceName())
	}

	return http.Ok(c, "Service group removed successfuly", ServiceDiscoveryServiceName())
}

// GetServiceGroupHandler handle GET /
func GetServiceGroupHandler(c *routing.Context) error {
	serviceGroupID := string(c.Param("group"))

	// Create service group service
	serviceGroupService, _ := getServiceGroupsService(c)
	
	sg, err := serviceGroupService.GetServiceGroupById(serviceGroupID)
	if err != nil {
		return http.Error(c, strconv.Quote(err.Error()), 400, ServiceDiscoveryServiceName())
	}

	sgByte, _ := json.Marshal(sg)
	return http.Ok(c, string(sgByte), ServiceDiscoveryServiceName())
}
