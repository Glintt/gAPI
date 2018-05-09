<template>
    <div class="home">
        <div class="alert alert-dark col-sm-3" role="alert">
            <strong>gAPI Base Url:</strong> 
            {{ 'http://' + $config.API.HOST + ':' + $config.API.PORT }}<br />
            <small>Use this URL + gAPIPath to call microservices</small>
        </div>
        <div class="row">
          <div class="col-sm-5 form-inline ">
              <input class="form-control" v-model="searchText"/>
              <button class="btn btn-sm btn-info" @click="updateData">
                <i class="fas fa-search"></i>
              </button>
          </div>
          <div class="col-sm-3 offset-sm-4 form-inline">
              <button class="btn btn-sm btn-info" @click="currentPage - 1 < 1 ? currentPage = 1 : currentPage -= 1">
                <i class="fas fa-arrow-left"></i>
              </button>
              <input class="form-control sm-1 mr-sm-1" v-model="currentPage" />
              <button class="btn btn-sm btn-info" @click="currentPage += 1">
                <i class="fas fa-arrow-right"></i>
              </button>
          </div>
        </div>
        <table class="table">
            <thead>
                <tr class="text-success">
                    <th scope="col">Name</th>
                    <th scope="col">gAPI Path</th>
                    <th scope="col">API Documentation</th>
                    <th scope="col">Health</th>
                    <th scope="col" v-show="isLoggedIn">Secured?</th>
                    <th scope="col">Actions</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="service in services" v-bind:key="service.Name">
                    <td>{{ service.Name }}</td>
                    <td>{{ service.MatchingURI }}</td>
                    <td>{{ service.APIDocumentation }}</td>
                    <td><i class="fas fa-heartbeat " :class="service.IsActive ? 'text-success' : 'text-danger'"></i></td>
                    <td v-show="isLoggedIn"><i class="fas " :class="service.Protected ? 'fa-lock text-success' : 'fa-unlock text-danger'"></i></td>
                    <td>
                        <router-link :to="'/service-discovery/service?uri='+service.MatchingURI" class="navbar-brand" >
                            <i class="fas fa-info-circle"></i>
                        </router-link>
                    </td>
                </tr>
            </tbody>
        </table>
    </div>
</template>

<script>
import DataTable from "@/components/DataTable";
var serviceDiscoveryAPI = require("@/api/service-discovery");

export default {
  name: "home",
  mounted() {
    this.updateData();
  },
  watch: {
    currentPage: function() {
      this.updateData();
    }
  },
  computed: {
    isLoggedIn() {
      return this.$oauthUtils.vmA.isLoggedIn();
    }
  },
  data() {
    return {
      services: [],
      currentPage: 1,
      searchText: ""
    };
  },
  methods: {
    updateData: function() {
      serviceDiscoveryAPI.listServices(
        this.currentPage, 
        this.searchText, 
        response => {
          if (response.status != 200) {
            this.services = [];
            return;
          }
          this.services = response.body;
        });
    }
  },
  components: {
    DataTable
  }
};
</script>
