<template>
  <table class="table">
    <thead>
      <tr class="table-secondary">
        <th scope="col">Name</th>
        <th scope="col">API Base Path</th>
        <th scope="col">Documentation</th>
        <th scope="col">Health</th>
        <th scope="col" v-show="isLoggedIn">Secured?</th>
        <th scope="col">Actions</th>
      </tr>
    </thead>
    <tbody>
      <tr
        v-for="service in services"
        v-bind:key="service.Id"
        v-if="service.IsReachable || (service.UseGroupAttributes && service.GroupVisibility ) || loggedInUser"
      >
        <td>{{ service.Name }}</td>
        <td>{{ apiPath(service)}}</td>
        <td>
          <a
            :href="$utils.urlConcat($config.API.getApiBaseUrl(),`api_docs/${service.Identifier}/documentation`)"
            target="_blank"
          >
            <i class="fas fa-book"></i> Documentation
          </a>
        </td>
        <td>
          <i class="fas fa-heartbeat" :class="service.IsActive ? 'text-success' : 'text-danger'"></i>
        </td>
        <td v-show="isLoggedIn">
          <i
            class="fas"
            :class="service.Protected ? 'fa-lock text-success' : 'fa-unlock text-danger'"
          ></i>
        </td>
        <td style="max-width: 20rem">
          <router-link
            :to="'/service-discovery/service?uri='+service.MatchingURI"
            data-toggle="tooltip"
            title="More info"
            style="margin-right: 1em"
          >
            <i class="fas fa-cog"></i>
          </router-link>
          <i
            class="fas fa-desktop text-success"
            data-toggle="tooltip"
            title="Manage Service"
            style="cursor:pointer"
            @click="showManageModal(service)"
            v-show="isLoggedIn && loggedInUser.IsAdmin"
          ></i>
          <!-- <button class="btn btn-success" @click="showManageModal(service)" v-show="isLoggedIn && loggedInUser.IsAdmin">
                        <i class="fas fa-desktop"></i> Manage
          </button>-->

          <router-link
            :to="'/service-discovery/service/logs?uri='+service.MatchingURI"
            v-show="isLoggedIn"
            data-toggle="tooltip"
            title="View service logs"
            style="margin-left: 1em"
          >
            <i class="fas fa-file text-warning"></i>
          </router-link>
        </td>
      </tr>
    </tbody>
  </table>
</template>

<script>
const pathModule = require("url");

export default {
  name: "list-services",
  props: ["services", "isLoggedIn", "loggedInUser"],
  methods: {
    showManageModal: function(service) {
      this.$emit("showManageModal", service);
    },
    apiPath(service) {
      return pathModule.resolve(
        this.$config.API.getApiBaseUrl(),
        service.MatchingURI
      );
    }
  }
};
</script>

<style>
</style>
