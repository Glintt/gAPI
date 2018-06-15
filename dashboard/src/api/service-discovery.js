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
  "list_groups": "/service-groups"
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
  console.log(obj)
  HTTP.POST(HTTP.PathToCall(ServiceDiscoveryBaseURL + Endpoints.add_to_group.replace('<group_id>', groupId)), obj , {}).then(response => {
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