const HTTP = require("@/api/http");

export function me(token, cb) {
    return HTTP.GET(HTTP.PathToCall("/oauth/me"), {}).then((response) => {
      cb(response);
    }, (response) => {
      cb(response);
    });
  }
  