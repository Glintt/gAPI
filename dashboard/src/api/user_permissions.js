const HTTP = require("@/api/http");

export function get(username, cb) {
  return HTTP.GET(
    HTTP.PathToCall(`/user-permissions/${username}/groups`),
    {}
  ).then(
    response => {
      cb(response);
    },
    response => {
      HTTP.handleError(response, cb);
    }
  );
}

export function add(username, applicationId, cb) {
  return HTTP.POST(
    HTTP.PathToCall(`/user-permissions/${username}/${applicationId}`),
    {},
    {}
  ).then(
    response => {
      cb(response);
    },
    response => {
      HTTP.handleError(response, cb);
    }
  );
}

export function remove(username, applicationId, cb) {
  return HTTP.DELETE(
    HTTP.PathToCall(`/user-permissions/${username}/${applicationId}`),
    {},
    {}
  ).then(
    response => {
      cb(response);
    },
    response => {
      HTTP.handleError(response, cb);
    }
  );
}
