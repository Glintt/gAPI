<template>
    <div class="home">
        
        <div class="row">
          <div class="col-sm-4 form-inline ">
            <form v-on:keyup.13="search" style="width: 100%">
              <input class="form-control" style="width: 80%" v-model="searchText" placeholder="Search ..." />
              <button class="btn btn-sm btn-info" @click="search">
                <i class="fas fa-search"></i>
              </button>
            </form>
          </div>
          <div class="alert alert-info col-sm-4" role="alert">
            <strong>gAPI Base Url:</strong>
            <span>&nbsp;&nbsp;{{ 'http://' + $config.API.HOST + ':' + $config.API.PORT }}</span><br />
            <small>Use this URL + gAPIPath to call microservices</small>
        </div>
          <div class="col-sm-3 offset-sm-1 form-inline">
              <button class="btn btn-sm btn-info" @click="currentPage - 1 < 1 ? currentPage = 1 : currentPage -= 1">
                <i class="fas fa-arrow-left"></i>
              </button>
              <input class="form-control sm-1 mr-sm-1" v-model="currentPage" />
              <button class="btn btn-sm btn-info" @click="currentPage += 1">
                <i class="fas fa-arrow-right"></i>
              </button>
          </div>
        </div>
        <br />
        <table class="table">
            <thead>
                <tr class="table-secondary" >
                    <th scope="col">Name</th>
                    <th scope="col">gAPI Path</th>
                    <th scope="col">API Documentation</th>
                    <th scope="col">Health</th>
                    <th scope="col" v-show="isLoggedIn">Secured?</th>
                    <th scope="col">Actions</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="service in services" v-bind:key="service.Id" v-if="service.IsReachable || (service.UseGroupAttributes && service.GroupVisibility ) || loggedInUser">
                    <td>{{ service.Name }}</td>
                    <td>{{ service.MatchingURI }}</td>
                    <td>{{ service.APIDocumentation }}</td>
                    <td><i class="fas fa-heartbeat " :class="service.IsActive ? 'text-success' : 'text-danger'"></i></td>
                    <td v-show="isLoggedIn"><i class="fas " :class="service.Protected ? 'fa-lock text-success' : 'fa-unlock text-danger'"></i></td>
                    <td style="max-width: 20rem">
                      <router-link :to="'/service-discovery/service?uri='+service.MatchingURI" class="btn btn-info" 
                          data-toggle="tooltip" data-placement="top" title="More info">
                          <i class="fas fa-info-circle"></i> Info
                      </router-link>

                      <button class="btn btn-success" @click="showManageModal(service)" v-show="isLoggedIn && loggedInUser.IsAdmin">
                        <i class="fas fa-desktop"></i> Manage
                      </button>
                    
                      <router-link :to="'/service-discovery/service/logs?uri='+service.MatchingURI" class="btn btn-warning" v-show="isLoggedIn"
                          data-toggle="tooltip" data-placement="top" title="View service logs">
                          <i class="fas fa-file"></i> Logs
                      </router-link>
                    </td>
                </tr>
            </tbody>
        </table>
        <ServiceManagementModal @modalClosed="toggleManageModal" :showing="manageModal.showing" 
              :id="'manageModal'" 
              :title="'Manage Service - ' + manageModal.service.Name"
              :service="manageModal.service"/>
    </div>
</template>

<script>
var serviceDiscoveryAPI = require("@/api/service-discovery");
import DataTable from "@/components/DataTable";
import ServiceManagementModal from "@/components/modals/ServiceManagementModal";
import { mapGetters } from 'vuex'

export default {
  name: "list-services",
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
    },
    ...mapGetters({
      loggedInUser: 'loggedInUser'
    })
  },
  data() {
    return {
      services: [],
      currentPage: 1,
      searchText: "",
      manageModal: {
        showing: false,
        service: {}
      }
    };
  },
  methods: {
    search: function() {
        event.preventDefault();

      this.currentPage = 1;
      this.updateData();
    },
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
        }
      );
    },
    showManageModal: function(service) {
      this.manageModal.service = service;
      this.toggleManageModal();
    },
    toggleManageModal: function() {
      this.manageModal.showing = !this.manageModal.showing;
    }
  },
  components: {
    ServiceManagementModal,
    DataTable
  }
};
</script>
