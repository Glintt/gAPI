const HTTP = require("@/api/http");

export function get(user, cb) {
  return HTTP.GET(
    HTTP.PathToCall(`/user-permissions/${user.Username}`),
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

export function update(user, newPermissions, cb) {
  return HTTP.PUT(
    HTTP.PathToCall(`/user-permissions/${user.Username}`),
    newPermissions,
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
