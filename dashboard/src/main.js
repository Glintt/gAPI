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

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount("#app");