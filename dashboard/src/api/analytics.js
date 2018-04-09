const HTTP = require("@/api/http");

const APIConfig = require("@/configs/urls").config.API;
const AnalyticsBaseURL = APIConfig.ANALYTICS_BASEPATH;

export function byApi(api, cb){
  return HTTP.GET(HTTP.PathToCall(AnalyticsBaseURL + "/api"), {params:api}).then(response => {
      cb(response);
    }, response => {
      cb(response);
    });
}
