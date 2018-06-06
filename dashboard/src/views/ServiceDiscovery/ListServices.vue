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
                <tr v-for="service in services">
                    <td>{{ service.Name }}</td>
                    <td>{{ service.MatchingURI }}</td>
                    <td>{{ service.APIDocumentation }}</td>
                    <td><i class="fas fa-heartbeat " :class="service.IsActive ? 'text-success' : 'text-danger'"></i></td>
                    <td v-show="isLoggedIn"><i class="fas " :class="service.Protected ? 'fa-lock text-success' : 'fa-unlock text-danger'"></i></td>
                    <td>
                        <router-link :to="'/service-discovery/service?uri='+service.MatchingURI" class="btn btn-sm btn-info" 
                            data-toggle="tooltip" data-placement="top" title="More info">
                            <i class="fas fa-info-circle"></i>
                        </router-link>
                        <button v-for="(type, index) in managementTypes" @click="manageService(service, type.action)" 
                            :class="'btn btn-sm btn-'+ type.background"
                            v-show="isLoggedIn && ! $api.serviceDiscovery.CustomManagementActions.includes(type.action)"
                            data-toggle="tooltip" data-placement="top" :title="type.description">
                            <i :class="type.icon"></i>
                        </button>
                        <router-link :to="'/service-discovery/service/logs?uri='+service.MatchingURI" class="btn btn-sm btn-info" v-show="isLoggedIn"
                            data-toggle="tooltip" data-placement="top" title="View service logs">
                            <i class="fas fa-file"></i>
                        </router-link>
                    </td>
                </tr>
            </tbody>
        </table>
        <ErrorMessage @modalClosed="statusModalClosed" :showing="statusMessage.showing && statusMessage.isError" :id="'requestError'" :error="statusMessage.msg" :title="'Error Occurred'"/>
        <SuccessModal @modalClosed="statusModalClosed" :showing="statusMessage.showing && !statusMessage.isError" :id="'requestSuccess'" :msg="statusMessage.msg" :title="'Success'"/>
        <ConfirmationModal @answerReceived="managementConfirmationReceived" @modalClosed="confirmationClosed" :showing="confirmation.showing" :id="'managementConfirm'" :msg="confirmation.msg" :title="confirmation.title"/>
    </div>
</template>

<script>
var serviceDiscoveryAPI = require("@/api/service-discovery");
import DataTable from "@/components/DataTable";
import ErrorMessage from "@/components/modals/ErrorMessage";
import ConfirmationModal from "@/components/modals/ConfirmationModal";
import SuccessModal from "@/components/modals/SuccessModal";

export default {
  name: "home",
  mounted() {
    this.updateData();
    this.$api.serviceDiscovery.manageServiceTypes(response => {
      this.managementTypes = response.body;
    });
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
      managementTypes: {},
      services: [],
      statusMessage: {
        msg: "",
        showing: false,
        isError: false
      },
      management: {
        action: "",
        service: null
      },
      confirmation: {
        showing: false,
        title: "",
        msg: ""
      },
      currentPage: 1,
      searchText: ""
    };
  },
  methods: {
    manageService: function(service, action) {
      this.confirmation.showing = true;
      this.confirmation.title = "Confirm - " + action;
      this.confirmation.msg =
        "Are you sure you want to " + action + " service " + service.Name + "?";
      this.management.service = service;
      this.management.action = action;
    },
    managementConfirmationReceived: function(answer) {
      if (answer == false) return;

      serviceDiscoveryAPI.manageService(
        this.management.service.MatchingURI,
        this.management.action,
        response => {
          this.statusMessage.msg = response.body.msg;
          this.statusMessage.isError = false;
          if (response.status != 200) {
            this.statusMessage.isError = true;
            if (response.body.service_response != undefined) {
              this.statusMessage.msg = response.body.service_response;
            }
          }
          this.statusMessage.showing = true;
        }
      );
    },
    confirmationClosed: function() {
      this.confirmation.showing = false;
      this.confirmation.msg = "";
    },
    statusModalClosed: function() {
      this.statusMessage.showing = false;
      this.statusMessage.msg = "";
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
    }
  },
  components: {
    ErrorMessage,
    SuccessModal,
    ConfirmationModal,
    DataTable
  }
};
</script>
