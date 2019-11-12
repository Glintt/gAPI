<template>
    <div class="home"> 
      
          <!-- <div class="row">
            <div class="alert alert-info col-sm-4 offset-sm-4" role="alert">
              <strong>gAPI Base Url:</strong>
              <span>&nbsp;&nbsp;{{ `${$config.API.getApiBaseUrl()}` }}</span><br />
              <small>Use this URL + gAPIPath to call microservices</small>
            </div>
          </div>       -->
      <div class="row">
        <div class="col-sm-2">
            <h5>Application Groups</h5><hr/>
            <ul class="list-group">
              <li :class="'list-group-item clickable ' + (appFilter == null ? 'active' : '')" @click="updateData()">All</li>
              <li :class="'list-group-item clickable ' + (appFilter == g.Id ? 'active' : '')" v-for="g in groups" @click="filterByApp(g)" v-bind:key="g.Id">
                {{g.Name}}
              </li>
            </ul>
        </div>
        <div class="col-sm-10">
          <div class="row">
            <div class="col-sm-12"> 
              <div class="row">
                <div class="col-sm-12 text-center form-inline ">
                  <form v-on:keyup.13="search" style="width: 100%">
                    <input class="form-control" style="width: 100%" v-model="searchText" placeholder="Search ..." />
                    <button hidden="true" class="btn btn-info" @click="search">
                      <i class="fas fa-search"></i>
                    </button>
                  </form>
                </div>
              </div>
              <br />
        
              <ListServices :services="services" :isLoggedIn="isLoggedIn" :loggedInUser="loggedInUser" v-on:showManageModal="showManageModal"/>
        
            </div>
          </div>
        </div>
      </div>
      <div class="row">
        <div class="col-sm-2 offset-sm-5 text-center">
          <nav aria-label="...">
            <ul class="pagination">
              <li :class="'page-item' + ((currentPage === 1) ? ' disabled' : '')">
                <a class="page-link" href="#" tabindex="-1" @click="currentPage - 1 < 1 ? currentPage = 1 : currentPage -= 1">Previous</a>
              </li>
              <li class="page-item" v-if="currentPage > 1">
                <a class="page-link" href="#" @click="currentPage -= 1" >{{ currentPage - 1}}</a>
              </li>
              <li class="page-item active">
                <a class="page-link" href="#">{{ currentPage }} <span class="sr-only">(current)</span></a>
              </li>
              <li class="page-item" v-if="services.length === 10">
                <a class="page-link" href="#" @click="currentPage += 1" >{{ currentPage + 1}}</a>
              </li>
              <li :class="'page-item' + ((services.length < 10) ? ' disabled' : '')">
                <a class="page-link" @click="currentPage += 1" href="#">Next</a>
              </li>
            </ul>
          </nav>
        </div>
      </div>
      
      <ServiceManagementModal @modalClosed="toggleManageModal" :showing="manageModal.showing" 
            :id="'manageModal'" 
            :title="'Manage Service - ' + manageModal.service.Name"
            :service="manageModal.service"/>
    </div>
</template>

<script>
import DataTable from "@/components/DataTable";
import ServiceManagementModal from "@/components/modals/ServiceManagementModal";
import ListServices from "@/components/service-discovery/ListServices";
import { mapGetters, mapActions } from "vuex";
var serviceDiscoveryAPI = require("@/api/service-discovery");

export default {
  name: "list-services",
  mounted() {
    this.updateData();
    this.fetchGroups();
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
      loggedInUser: "loggedInUser"
    }),
    ...mapGetters("appsGroups", [
      "groups",
      "ungroupedApplications",
      "possibleMatches"
    ])
  },
  data() {
    return {
      services: [],
      selectedGroup: null,
      currentPage: 1,
      searchText: "",
      appFilter:null,
      manageModal: {
        showing: false,
        service: {}
      }
    };
  },
  methods: {
    ...mapActions("appsGroups", [
      "fetchGroups"
    ]),
    filterByApp: function(appGroup) {
      this.appFilter = appGroup.Id;
      this.$api.serviceDiscovery.applicationGroupById(appGroup.Id, response => {
        this.selectedGroup = response.body;
        this.services = response.body.Services;
      }); 
    },
    search: function(event) {
      event.preventDefault();

      this.currentPage = 1;
      this.updateData();
    },
    updateData: function() {
      this.appFilter = null;
      serviceDiscoveryAPI.listServices(
        this.currentPage,
        this.searchText,
        response => {
          if (response.status !== 200) {
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
    DataTable,
    ListServices
  }
};
</script>
