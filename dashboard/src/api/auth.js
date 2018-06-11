const HTTP = require("@/api/http");
const ACCESS_TOKEN_KEY = "access_token_key";
const ACCESS_TOKEN_EXPIRATION_TIME = "access_token_expiration_time";

export function authenticate(user, password,cb) {
    return HTTP.POST(HTTP.PathToCall("/oauth/token"), {"username":user, "password":password}, {}).then((response) => {
      
      storeToken(response.body.token, response.body.expiration_time);
      
      cb(response);
    }, (response) => {
      HTTP.handleError(response, cb)
    });
}

export function storeToken(token, expirationTime){
  localStorage.setItem(ACCESS_TOKEN_KEY, token);
  var whenExpires = new Date().getTime() + ((expirationTime - 2) * 1000);
  localStorage.setItem(ACCESS_TOKEN_EXPIRATION_TIME, whenExpires);
}

export function getToken(){
  return localStorage.getItem(ACCESS_TOKEN_KEY);
}

export function getExpirationTime(){
  return localStorage.getItem(ACCESS_TOKEN_EXPIRATION_TIME);
}

export function clearAccessToken() {
  localStorage.removeItem(ACCESS_TOKEN_KEY);
  localStorage.removeItem(ACCESS_TOKEN_EXPIRATION_TIME);
}