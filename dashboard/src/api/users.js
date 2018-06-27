const HTTP = require("@/api/http");

export function me(token, cb) {
  return HTTP.GET(HTTP.PathToCall("/oauth/me"), {}).then((response) => {
    cb(response);
  }, (response) => {
    HTTP.handleError(response, cb)
  });
}

export function find(query, page, cb) {
  if (query == undefined) {
    query = ""
  }
  if (page == undefined || page < 1) page = 1
  return HTTP.GET(HTTP.PathToCall("/users?q="+query+"&page="+page), {}).then((response) => {
    cb(response);
  }, (response) => {
    HTTP.handleError(response, cb)
  });
}
  
export function update(user, cb) {
  return HTTP.PUT(HTTP.PathToCall(`/users/${user.Username}`), user, {}).then((response) => {
    cb(response);
  }, (response) => {
    HTTP.handleError(response, cb)
  });
}
  
export function updateByAdmin(user, cb) {
  return HTTP.PUT(HTTP.PathToCall(`/users/admin/${user.Username}`), user, {}).then((response) => {
    cb(response);
  }, (response) => {
    HTTP.handleError(response, cb)
  });
}

export function create(user, cb) {
  return HTTP.POST(HTTP.PathToCall(`/users`), user, {}).then((response) => {
    cb(response);
  }, (response) => {
    HTTP.handleError(response, cb)
  });
}