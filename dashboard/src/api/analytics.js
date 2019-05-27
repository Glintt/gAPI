const HTTP = require("@/api/http");

const APIConfig = require("@/configs/urls").config.API;
const AnalyticsBaseURL = APIConfig.ANALYTICS_BASEPATH;

export function byApi(api, cb) {
  return HTTP.GET(HTTP.PathToCall(AnalyticsBaseURL + "/api"), {
    params: api
  }).then(
    response => {
      cb(response);
    },
    response => {
      HTTP.handleError(response, cb);
    }
  );
}

export function byApplication(app_id, cb) {
  return HTTP.GET(HTTP.PathToCall(AnalyticsBaseURL + "/applications"), {
    params: app_id
  }).then(
    response => {
      cb(response);
    },
    response => {
      HTTP.handleError(response, cb);
    }
  );
}

export function logs(api, cb) {
  return HTTP.GET(HTTP.PathToCall(AnalyticsBaseURL + "/logs"), {
    params: api
  }).then(
    response => {
      cb(response);
    },
    response => {
      HTTP.handleError(response, cb);
    }
  );
}
