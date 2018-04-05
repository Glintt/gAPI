
const OAUTH_API = require("@/api/auth");
import Vue from "vue";

  
export function logout() {
    OAUTH_API.clearAccessToken();
}
  
export function isLoggedIn() {
    return OAUTH_API.getToken();
}

export function requireAuth(to, from, next) {
    if (!isLoggedIn()) {
        next({
            path: '/login'
          });
    } else {
      next();
    }
}
  