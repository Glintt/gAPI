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

Vue.prototype.$config.API.HOST = Vue.prototype.$utils.getCookieByName("API_HOST") == null ? Vue.prototype.$config.API.HOST : Vue.prototype.$utils.getCookieByName("API_HOST");
Vue.prototype.$config.API.PORT = Vue.prototype.$utils.getCookieByName("API_PORT") == null ? Vue.prototype.$config.API.PORT : Vue.prototype.$utils.getCookieByName("API_PORT");
Vue.prototype.$config.API.SOCKET_HOST = Vue.prototype.$utils.getCookieByName("SOCKET_HOST") == null ? Vue.prototype.$config.API.SOCKET_HOST : Vue.prototype.$utils.getCookieByName("SOCKET_HOST");
Vue.prototype.$config.API.SOCKET_PORT = Vue.prototype.$utils.getCookieByName("SOCKET_PORT") == null ? Vue.prototype.$config.API.SOCKET_PORT : Vue.prototype.$utils.getCookieByName("SOCKET_PORT");

console.log(Vue.prototype.$config.API)
console.log(Vue.prototype.$utils.getCookieByName("API_HOST"))

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount("#app");
