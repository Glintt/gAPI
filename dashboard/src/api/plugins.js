const HTTP = require("@/api/http");

const Endpoints = {
  all: "/plugins",
  active: "/plugins/active"
};

export function all(cb) {
  return HTTP.GET(HTTP.PathToCall(Endpoints.all), {}).then(
    response => {
      cb(response);
    },
    response => {
      HTTP.handleError(response, cb);
    }
  );
}

export function active(cb) {
  return HTTP.GET(HTTP.PathToCall(Endpoints.active), {}).then(
    response => {
      cb(response);
    },
    response => {
      HTTP.handleError(response, cb);
    }
  );
}
