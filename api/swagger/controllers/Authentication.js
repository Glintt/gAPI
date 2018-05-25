'use strict';

var utils = require('../utils/writer.js');
var Authentication = require('../service/AuthenticationService');

module.exports.oauthAuthorizeGET = function oauthAuthorizeGET (req, res, next) {
  var authorization = req.swagger.params['Authorization'].value;
  Authentication.oauthAuthorizeGET(authorization)
    .then(function (response) {
      utils.writeJson(res, response);
    })
    .catch(function (response) {
      utils.writeJson(res, response);
    });
};

module.exports.oauthTokenPOST = function oauthTokenPOST (req, res, next) {
  var body = req.swagger.params['body'].value;
  Authentication.oauthTokenPOST(body)
    .then(function (response) {
      utils.writeJson(res, response);
    })
    .catch(function (response) {
      utils.writeJson(res, response);
    });
};
