'use strict';


/**
 * Delete service from service discovery
 * 
 *
 * authorization String Authorization token
 * uri String Matching uri to search for
 * no response value expected for this operation
 **/
exports.service_discoveryDeleteDELETE = function(authorization,uri) {
  return new Promise(function(resolve, reject) {
    resolve();
  });
}


/**
 * Search service by matching uri
 * 
 *
 * uri String Matching uri to search for
 * no response value expected for this operation
 **/
exports.service_discoveryEndpointGET = function(uri) {
  return new Promise(function(resolve, reject) {
    resolve();
  });
}


/**
 * Register new service
 * 
 *
 * authorization String Authorization token
 * body Service Service information
 * no response value expected for this operation
 **/
exports.service_discoveryRegisterPOST = function(authorization,body) {
  return new Promise(function(resolve, reject) {
    resolve();
  });
}


/**
 * List all services registered
 * 
 *
 * authorization String Authorization token
 * q String Search query (optional)
 * page Integer Page to fetch (optional)
 * no response value expected for this operation
 **/
exports.service_discoveryServicesGET = function(authorization,q,page) {
  return new Promise(function(resolve, reject) {
    resolve();
  });
}


/**
 * Call action over a service
 * 
 *
 * authorization String Authorization token
 * service String Service matching uri to which action will be applied
 * action String Action to apply. Available actions: call /service-discovery/services/manage/types endpoint
 * no response value expected for this operation
 **/
exports.service_discoveryServicesManagePOST = function(authorization,service,action) {
  return new Promise(function(resolve, reject) {
    resolve();
  });
}


/**
 * Get management available types
 * 
 *
 * authorization String Authorization token
 * no response value expected for this operation
 **/
exports.service_discoveryServicesManageTypesGET = function(authorization) {
  return new Promise(function(resolve, reject) {
    resolve();
  });
}


/**
 * Update service on service discovery
 * 
 *
 * authorization String Authorization token
 * body Service Service update information
 * no response value expected for this operation
 **/
exports.service_discoveryUpdatePOST = function(authorization,body) {
  return new Promise(function(resolve, reject) {
    resolve();
  });
}

