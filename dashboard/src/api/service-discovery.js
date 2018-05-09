const HTTP = require("@/api/http");
const APIConfig = require("@/configs/urls").config.API;
const ServiceDiscoveryBaseURL = APIConfig.SERVICE_DISCOVERY_BASEPATH;

export function listServices(page, searchQuery, cb){
  HTTP.GET(HTTP.PathToCall(ServiceDiscoveryBaseURL + "/services?page=" + page + "&q=" + searchQuery), {}).then(response => {
      cb(response);
    }, response => {
      cb(response)
  });
}

export function getServices(serviceEndpoint, cb){
  HTTP.GET(HTTP.PathToCall(ServiceDiscoveryBaseURL + "/endpoint"), {params: {uri: serviceEndpoint}}).then(response => {
      cb(response);
    }, response => {
      cb(response);
    });
}

export function storeService(service, cb){
  HTTP.POST(HTTP.PathToCall(ServiceDiscoveryBaseURL + "/register"), service, {}).then(response => {
    cb(response);
  }, response => {
    cb(response);
  });
}

export function deleteService(serviceEndpoint, cb){
  HTTP.DELETE(HTTP.PathToCall(ServiceDiscoveryBaseURL + "/delete"), {params: {uri: serviceEndpoint}}).then(response => {
    cb(response);
  }, response => {
    cb(response);
  });
}

export function updateService(service, cb){
  HTTP.POST(HTTP.PathToCall(ServiceDiscoveryBaseURL + "/update"), service, {}).then(response => {
    cb(response);
  }, response => {
    HTTP.handleError(response, cb)
  });
}

export function refreshService(service, cb){
  HTTP.POST(HTTP.PathToCall(ServiceDiscoveryBaseURL + "/service/restart?service=" + service), {}, {}).then(response => {
    cb(response);
  }, response => {
    HTTP.handleError(response, cb)
  });
}