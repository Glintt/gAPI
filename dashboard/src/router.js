import Vue from "vue";
import Router from "vue-router";
import Home from "./views/Home.vue";
import Login from "./views/Login.vue";
import ListServices from "./views/ServiceDiscovery/ListServices.vue";
import NewService from "./views/ServiceDiscovery/NewService.vue";
import ViewService from "./views/Service/ViewService.vue";
import ServiceLogs from "./views/Service/ServiceLogs.vue";
import ByApi from "./views/Analytics/ByApi.vue";
import Realtime from "./views/Analytics/Realtime.vue";
import ListUsers from "./views/Users/ListUsers.vue";
import EditUser from "./views/Users/EditUser.vue";
import NewUser from "./views/Users/NewUser.vue";
import NewServiceGroup from "./views/ServiceDiscovery/NewServiceGroup.vue";
import ListServicesGroup from "./views/ServiceDiscovery/ListServicesGroup.vue";
import NewApplicationGroup from "./views/ServiceDiscovery/NewApplicationGroup.vue";
import ListApplicationsGroup from "./views/ServiceDiscovery/ListApplicationsGroup.vue";


var OAuthValidator = require("@/auth");
Vue.use(Router);

export default new Router({
  mode: 'history',
  routes: [
    {
      path: "/",
      name: "home",
      component: Home
    },
    {
      path: "/users",
      name: "users",
      component: ListUsers,
      beforeEnter: OAuthValidator.requireAdminAuth
    },
    {
      path: "/users/create",
      name: "users-create",
      component: NewUser,
      beforeEnter: OAuthValidator.requireAdminAuth
    },
    {
      path: "/profile",
      name: "profile",
      component: EditUser,
      beforeEnter: OAuthValidator.requireAuth
    },
    {
      path: "/login",
      name: "login",
      component: Login
    },
    {
      path: "/service-discovery/services",
      name: "service-discovery-services",
      component: ListServices
    },
    {
      path: "/service-discovery/services/new",
      name: "service-discovery-services-new",
      component: NewService,
      beforeEnter: OAuthValidator.requireAdminAuth
    },
    {
      path: "/service-discovery/service",
      name: "service-view",
      component: ViewService
    },
    {
      path: "/service-discovery/service/logs",
      name: "service-logs-view",
      component: ServiceLogs,
      beforeEnter: OAuthValidator.requireAuth
    },
    {
      path: "/service-discovery/groups/create",
      name: "service-groups-create",
      component: NewServiceGroup,
      beforeEnter: OAuthValidator.requireAdminAuth
    },
    {
      path: "/service-discovery/apps-groups/create",
      name: "service-apps-groups-create",
      component: NewApplicationGroup,
      beforeEnter: OAuthValidator.requireAdminAuth
    },
    {
      path: "/service-discovery/apps-groups",
      name: "service-apps-groups-list",
      component: ListApplicationsGroup,
      beforeEnter: OAuthValidator.requireAdminAuth
    },
    {
      path: "/service-discovery/groups",
      name: "service-groups-list",
      component: ListServicesGroup,
      beforeEnter: OAuthValidator.requireAdminAuth
    },
    {
      path: "/analytics/by-api",
      name: "analytics-by-api",
      component: ByApi,
      beforeEnter: OAuthValidator.requireAdminAuth
    },
    {
      path: "/analytics/realtime",
      name: "analytics-realtime",
      component: Realtime,
      beforeEnter: OAuthValidator.requireAdminAuth
    }
  ]
});
