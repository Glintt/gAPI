import Vue from "vue";
import VueResource from "vue-resource";
Vue.use(VueResource);

import App from "./App.vue";
import router from "./router";
import store from "./store";
require("bootstrap");

Vue.config.productionTip = false;

Vue.prototype.$utils = require("@/utils");
Vue.prototype.$config = require("@/configs/urls").config;
Vue.prototype.$chartColors = require("@/configs/chartColors");
Vue.prototype.$random = require("@/configs/random");
Vue.prototype.$oauthUtils = require("@/auth");
Vue.prototype.$api = require("@/api/api").api;
Vue.prototype.$permissions = require("@/configs/permissions")

const UsersServices = require('@/api/users')

UsersServices.me(require("@/api/auth").getToken(), (response) => {
  store.dispatch('loggedInUserUpdate', response.body)
})

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount("#app");