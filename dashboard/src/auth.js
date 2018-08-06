import Vue from "vue";
const OAUTH_API = require("@/api/auth");
var store = require("./store");

export const vmA = new Vue({
  data: {
    loggedIn: false,
    user: null
  },
  methods: {
    logout() {
      OAUTH_API.clearAccessToken();
      this.loggedIn = false;
    },
    isLoggedIn() {
      this.loggedIn =
        OAUTH_API.getToken() &&
        new Date().getTime() < OAUTH_API.getExpirationTime();
      return this.loggedIn;
    },
    authenticate(user, cb) {
      OAUTH_API.authenticate(user.username, user.password, response => {
        cb(response);
      });
    },
    currentUser() {
      return this.user;
    }
  }
});

export function requireAdminAuth(to, from, next) {
  if (!vmA.isLoggedIn()) {
    next({
      path: "/login"
    });
  } else {
    if (!store.default.state.loggedInUser.IsAdmin) return next({ path: "/" });
    else next();
  }
}

export function requireAuth(to, from, next) {
  if (!vmA.isLoggedIn()) {
    next({
      path: "/login"
    });
  } else {
    next();
  }
}
