'use strict';

var utils = require('../utils/writer.js');
var Admin = require('../service/AdminService');

module.exports.invalidate_cacheGET = function invalidate_cacheGET (req, res, next) {
  var authorization = req.swagger.params['Authorization'].value;
  Admin.invalidate_cacheGET(authorization)
    .then(function (response) {
      utils.writeJson(res, response);
    })
    .catch(function (response) {
      utils.writeJson(res, response);
    });
};

module.exports.reloadGET = function reloadGET (req, res, next) {
  var authorization = req.swagger.params['Authorization'].value;
  Admin.reloadGET(authorization)
    .then(function (response) {
      utils.writeJson(res, response);
    })
    .catch(function (response) {
      utils.writeJson(res, response);
    });
};

module.exports.service_discoveryAdminNormalizePOST = function service_discoveryAdminNormalizePOST (req, res, next) {
  var authorization = req.swagger.params['Authorization'].value;
  Admin.service_discoveryAdminNormalizePOST(authorization)
    .then(function (response) {
      utils.writeJson(res, response);
    })
    .catch(function (response) {
      utils.writeJson(res, response);
    });
};
