'use strict';

var utils = require('../utils/writer.js');
var ServiceDiscovery = require('../service/ServiceDiscoveryService');

module.exports.service_discoveryDeleteDELETE = function service_discoveryDeleteDELETE (req, res, next) {
  var authorization = req.swagger.params['Authorization'].value;
  var uri = req.swagger.params['uri'].value;
  ServiceDiscovery.service_discoveryDeleteDELETE(authorization,uri)
    .then(function (response) {
      utils.writeJson(res, response);
    })
    .catch(function (response) {
      utils.writeJson(res, response);
    });
};

module.exports.service_discoveryEndpointGET = function service_discoveryEndpointGET (req, res, next) {
  var uri = req.swagger.params['uri'].value;
  ServiceDiscovery.service_discoveryEndpointGET(uri)
    .then(function (response) {
      utils.writeJson(res, response);
    })
    .catch(function (response) {
      utils.writeJson(res, response);
    });
};

module.exports.service_discoveryRegisterPOST = function service_discoveryRegisterPOST (req, res, next) {
  var authorization = req.swagger.params['Authorization'].value;
  var body = req.swagger.params['body'].value;
  ServiceDiscovery.service_discoveryRegisterPOST(authorization,body)
    .then(function (response) {
      utils.writeJson(res, response);
    })
    .catch(function (response) {
      utils.writeJson(res, response);
    });
};

module.exports.service_discoveryServicesGET = function service_discoveryServicesGET (req, res, next) {
  var authorization = req.swagger.params['Authorization'].value;
  var q = req.swagger.params['q'].value;
  var page = req.swagger.params['page'].value;
  ServiceDiscovery.service_discoveryServicesGET(authorization,q,page)
    .then(function (response) {
      utils.writeJson(res, response);
    })
    .catch(function (response) {
      utils.writeJson(res, response);
    });
};

module.exports.service_discoveryServicesManagePOST = function service_discoveryServicesManagePOST (req, res, next) {
  var authorization = req.swagger.params['Authorization'].value;
  var service = req.swagger.params['service'].value;
  var action = req.swagger.params['action'].value;
  ServiceDiscovery.service_discoveryServicesManagePOST(authorization,service,action)
    .then(function (response) {
      utils.writeJson(res, response);
    })
    .catch(function (response) {
      utils.writeJson(res, response);
    });
};

module.exports.service_discoveryServicesManageTypesGET = function service_discoveryServicesManageTypesGET (req, res, next) {
  var authorization = req.swagger.params['Authorization'].value;
  ServiceDiscovery.service_discoveryServicesManageTypesGET(authorization)
    .then(function (response) {
      utils.writeJson(res, response);
    })
    .catch(function (response) {
      utils.writeJson(res, response);
    });
};

module.exports.service_discoveryUpdatePOST = function service_discoveryUpdatePOST (req, res, next) {
  var authorization = req.swagger.params['Authorization'].value;
  var body = req.swagger.params['body'].value;
  ServiceDiscovery.service_discoveryUpdatePOST(authorization,body)
    .then(function (response) {
      utils.writeJson(res, response);
    })
    .catch(function (response) {
      utils.writeJson(res, response);
    });
};
