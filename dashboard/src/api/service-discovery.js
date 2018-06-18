const HTTP = require("@/api/http");
const APIConfig = require("@/configs/urls").config.API;
const ServiceDiscoveryBaseURL = APIConfig.SERVICE_DISCOVERY_BASEPATH;

const Endpoints = {
  "list" : "/services",
  "get" : "/endpoint",
  "store" : "/register",
  "delete":"/delete",
  "manage":"/services/manage",
  "manage_types":"/services/manage/types",
  "update" :"/update",
  "add_to_group" : "/service-groups/<group_id>/services",
  "deassociate_from_group" : "/service-groups/<group_id>/services/<service_id>",
  "list_groups": "/service-groups",
  "store_group": "/service-groups",
  "update_group": "/service-groups/<group_id>",
  "remove_group": "/service-groups/<group_id>"

};

export const CustomManagementActions = ["logs"]

export function listServices(page, searchQuery, cb){
  HTTP.GET(HTTP.PathToCall(ServiceDiscoveryBaseURL + Endpoints.list + "?page=" + page + "&q=" + searchQuery), {}).then(response => {
      cb(response);
    }, response => {
      HTTP.handleError(response, cb)
  });
}

export function getServices(serviceEndpoint, cb){
  HTTP.GET(HTTP.PathToCall(ServiceDiscoveryBaseURL + Endpoints.get), {params: {uri: serviceEndpoint}}).then(response => {
      cb(response);
    }, response => {
      HTTP.handleError(response, cb)
    });
}

export function storeService(service, cb){
  HTTP.POST(HTTP.PathToCall(ServiceDiscoveryBaseURL + Endpoints.store), service, {}).then(response => {
    cb(response);
  }, response => {
    HTTP.handleError(response, cb)
  });
}

export function storeServiceGroup(group, cb){
  HTTP.POST(HTTP.PathToCall(ServiceDiscoveryBaseURL + Endpoints.store_group), group, {}).then(response => {
    cb(response);
  }, response => {
    HTTP.handleError(response, cb)
  });
}

export function addServiceToServiceGroup(groupId, serviceId, cb){
  let obj = {service_id: serviceId}
  HTTP.POST(HTTP.PathToCall(ServiceDiscoveryBaseURL + Endpoints.add_to_group.replace('<group_id>', groupId)), obj , {}).then(response => {
    cb(response);
  }, response => {
    HTTP.handleError(response, cb)
  });
}

export function deassociateServiceFromServiceGroup(groupId, serviceId, cb){
  HTTP.DELETE(HTTP.PathToCall(ServiceDiscoveryBaseURL + Endpoints.deassociate_from_group.replace('<group_id>', groupId).replace('<service_id>', serviceId)),  {}).then(response => {
    cb(response);
  }, response => {
    HTTP.handleError(response, cb)
  });
}

export function listServiceGroups(cb){
  HTTP.GET(HTTP.PathToCall(ServiceDiscoveryBaseURL + Endpoints.list_groups), {}).then(response => {
    cb(response);
  }, response => {
    HTTP.handleError(response, cb)
  });
}

export function updateServiceGroup(group, cb){
  HTTP.PUT(HTTP.PathToCall(ServiceDiscoveryBaseURL + Endpoints.update_group.replace('<group_id>', group.Id)), group, {}).then(response => {
    cb(response);
  }, response => {
    HTTP.handleError(response, cb)
  });
}
export function deleteServiceGroup(groupId, cb){
  HTTP.DELETE(HTTP.PathToCall(ServiceDiscoveryBaseURL + Endpoints.remove_group.replace('<group_id>', groupId)), {}).then(response => {
    cb(response);
  }, response => {
    HTTP.handleError(response, cb)
  });
}



export function deleteService(serviceEndpoint, cb){
  HTTP.DELETE(HTTP.PathToCall(ServiceDiscoveryBaseURL + Endpoints.delete), {params: {uri: serviceEndpoint}}).then(response => {
    cb(response);
  }, response => {
    HTTP.handleError(response, cb)
  });
}

export function updateService(service, cb){
  HTTP.POST(HTTP.PathToCall(ServiceDiscoveryBaseURL + Endpoints.update), service, {}).then(response => {
    cb(response);
  }, response => {
    HTTP.handleError(response, cb)
  });
}

export function manageService(service, action, cb){
  HTTP.POST(HTTP.PathToCall(ServiceDiscoveryBaseURL + Endpoints.manage +"?service=" + service + "&action="+action), {}, {}).then(response => {
    cb(response);
  }, response => {
    HTTP.handleError(response, cb)
  });
}
export function manageServiceTypes(cb){
  HTTP.GET(HTTP.PathToCall(ServiceDiscoveryBaseURL + Endpoints.manage_types), {}, {}).then(response => {
    cb(response);
  }, response => {
    HTTP.handleError(response, cb)
  });
}