import Vue from "vue";
import VueResource from "vue-resource";
Vue.use(VueResource);

import App from "./App.vue";
import router from "./router";
import store from "./store";
require("bootstrap");

Vue.config.productionTip = false;

Vue.prototype.$chartColors = require("@/configs/chartColors");
Vue.prototype.$random = require("@/configs/random");
Vue.prototype.$utils = require("@/utils");

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount("#app");
