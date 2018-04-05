const HTTP = require("@/api/http");
const ACCESS_TOKEN_KEY = "access_token_key";

export function authenticate(user, password,cb) {
    return HTTP.POST(HTTP.PathToCall("/oauth/token"), {"username":password, "password":password}, {}).then((response) => {
      cb(response);
    }, (response) => {
      cb(response);
    });
}

export function storeToken(token){
  localStorage.setItem(ACCESS_TOKEN_KEY, token);
}

export function getToken(){
  return localStorage.getItem(ACCESS_TOKEN_KEY);
}

export function clearAccessToken() {
    localStorage.removeItem(ACCESS_TOKEN_KEY);
}