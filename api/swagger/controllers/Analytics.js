'use strict';

var utils = require('../utils/writer.js');
var Analytics = require('../service/AnalyticsService');

module.exports.analyticsApiGET = function analyticsApiGET (req, res, next) {
  var authorization = req.swagger.params['Authorization'].value;
  var endpoint = req.swagger.params['endpoint'].value;
  Analytics.analyticsApiGET(authorization,endpoint)
    .then(function (response) {
      utils.writeJson(res, response);
    })
    .catch(function (response) {
      utils.writeJson(res, response);
    });
};

module.exports.analyticsLogsGET = function analyticsLogsGET (req, res, next) {
  var authorization = req.swagger.params['Authorization'].value;
  var endpoint = req.swagger.params['endpoint'].value;
  Analytics.analyticsLogsGET(authorization,endpoint)
    .then(function (response) {
      utils.writeJson(res, response);
    })
    .catch(function (response) {
      utils.writeJson(res, response);
    });
};
