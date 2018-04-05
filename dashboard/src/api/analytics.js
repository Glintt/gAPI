const HTTP = require("@/api/http");

export function byApi(api, cb){
  return HTTP.GET(HTTP.PathToCall("/api"), {params:api}).then(response => {
      cb(response);
    }, response => {
      cb(response);
    });
}
