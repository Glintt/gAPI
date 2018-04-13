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

Vue.http.get("/hosts", (response) =>{
  if (response.status != 200) {
    return;
  }

  Vue.prototype.$config.API.HOST = response.body.API_HOST;
  Vue.prototype.$config.API.PORT = response.body.API_PORT;
  Vue.prototype.$config.API.SOCKET_HOST = response.body.SOCKET_HOST;
  Vue.prototype.$config.API.SOCKET_PORT = response.body.SOCKET_PORT;
});


new Vue({
  router,
  store,
  render: h => h(App)
}).$mount("#app");