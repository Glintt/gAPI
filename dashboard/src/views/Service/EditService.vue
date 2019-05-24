<template>
  <div class="row">
    <div class="col-sm">
      <InformationPanel
        v-if="informationStatus.isActive"
        :msg="informationStatus.msg"
        :className="informationStatus.className"
      ></InformationPanel>
      <div class="row"></div>
      <div class="row">
        <div class="col-sm">
          <div class="row">
            <div class="col-sm-10" v-show="isLoggedIn">
              <h2>{{ service.Name }} API Information</h2>
            </div>
            <div class="col-sm-2" v-show="isAdmin">
              <button type="submit" class="btn btn-xs btn-primary" v-on:click="store">Save</button>
              <button class="btn btn-danger" @click="deleteService">Delete</button>
              <!-- <button type="submit" class="btn btn-info" v-on:click="serviceUpdated">Preview</button> -->
            </div>
          </div>
        </div>
      </div>

      
        <ul class="nav nav-tabs">
          <li class="nav-item">
            <a class="nav-link active" id="home-tab" data-toggle="tab" href="#home" role="tab" aria-controls="home" aria-selected="true">
              Basic Configuration
            </a>
            </li>
          <li v-show="isLoggedIn"  class="nav-item">
            <a class="nav-link" id="config-tab" data-toggle="tab" href="#config" role="tab" aria-controls="config" aria-selected="true">
              Api Configuration
            </a>
          </li>
          <li v-show="isLoggedIn" class="nav-item">
            <a class="nav-link" id="mng-config-tab" data-toggle="tab" href="#mng-config" role="tab" aria-controls="mng-config" aria-selected="true">
              Management Configuration
            </a>
          </li>
          
        </ul>

          <div class="tab-content" id="myTabContent">
            <div class="tab-pane fade show active" id="home" role="tabpanel" aria-labelledby="home-tab">
              <div class="card mb-12">
                      <div
                        class="card-header text-white bg-primary toggable-card"
                        @click="toggleCard('basic')"
                      >
                        <div class="row">
                          <div
                            :class="service.LastActiveTime !== 0 ? 'col-sm-10' : 'col-sm-11'"
                          >Basic Information</div>
                          <div :class="service.LastActiveTime !== 0 ? 'col-sm-2' : 'col-sm-1'">
                            <span>Health:</span>
                            <i class="fas fa-heartbeat fa-lg" :class="isActiveClass"></i>
                            <small v-if="service.LastActiveTime !== 0">
                              <br>
                              Last Time Active: {{ this.$utils.convertMillisToTime(new Date().getTime() - service.LastActiveTime) }}
                            </small>
                          </div>
                        </div>
                      </div>
                      <div class="card-body row" v-if="cards.basic.showing">
                        <div class="form-group col-sm">
                          <label for="serviceName">Name</label>
                          <input
                            type="text"
                            :disabled="! (isLoggedIn && isAdmin)"
                            v-model="service.Name"
                            class="form-control"
                            id="serviceName"
                            aria-describedby="nameHelp"
                            placeholder="Enter name"
                          >
                          <small id="nameHelp" class="form-text text-primary">Give the service/API a name.</small>
                        </div>

                        <div class="form-group col-sm">
                          <label for="serviceMatchingUri">MatchingURI</label>
                          <input
                            type="text"
                            :disabled="! $permissions.HasPermission(pageType, loggedInUser)"
                            v-model="service.MatchingURI"
                            class="form-control"
                            id="serviceMatchingUri"
                            aria-describedby="serviceMatchingUriHelp"
                            placeholder="Enter domain"
                          >
                          <small
                            id="serviceMatchingUriHelp"
                            class="form-text text-primary"
                          >Base URI which links to the service on API Management Platform.</small>
                        </div>

                        <div class="form-group col-sm">
                          <label for="serviceName">URL:</label>
                          <div class="input-group mb-2">
                            <input
                              type="text"
                              :disabled="true"
                              :value="serviceURL()"
                              class="form-control"
                              id="gapiBasePath"
                              aria-describedby="nameHelp"
                              placeholder="Enter name"
                            >
                            <button class="btn btn-success" @click="copyURL">
                              <i class="fas fa-clipboard"></i> Copy
                            </button>
                          </div>
                          <small id="nameHelp" class="form-text text-primary">Base Path to call microservice.</small>
                        </div>
                      </div>
                    </div>
            
            </div>
            <div class="tab-pane fade" id="config" role="tabpanel" aria-labelledby="config-tab">
              <ServiceAPIConfiguration
                      class="toggable-card"
                      v-on:addEndpointExclude="addEndpointExclude"
                      v-on:removeEndpointExclude="removeEndpointExclude"
                      v-on:addHost="addHost"
                      v-show="isLoggedIn"
                      v-on:toggleCard="toggleCard"
                      v-on:removeHost="removeHost"
                      :showing="cards.api_config.showing"
                      :service="service"
                    />
            </div>
            <div class="tab-pane fade" id="mng-config" role="tabpanel" aria-labelledby="mng-config-tab">
              <ServiceManagementConfig
                      class="toggable-card"
                      v-on:toggleCard="toggleCard"
                      :showing="cards.management_config.showing"
                      :service="service"
                      v-show="isLoggedIn"
                    />
                    </div>
          </div>
    </div>
  </div>
  
</template>

<script>
import InformationPanel from "@/components/InformationPanel";
import ServiceAPIConfiguration from "@/views/Service/ServiceAPIConfiguration";
import ServiceManagementConfig from "@/views/Service/ServiceManagementConfig";
import { mapGetters } from "vuex";
var serviceDiscoveryAPI = require("@/api/service-discovery");

export default {
  name: "view-service",
  computed: {
    isLoggedIn() {
      return this.$oauthUtils.vmA.isLoggedIn();
    },
    ...mapGetters({
      isAdmin: "isAdmin",
      loggedInUser: "loggedInUser"
    })
  },
  mounted() {
    serviceDiscoveryAPI.getServices(this.$route.query.uri, response => {
      this.service = response.body;
      if (this.service.ServiceManagementEndpoints === null) {
        this.service.ServiceManagementEndpoints = {};
      }
      if (this.service.Hosts === null) this.service.Hosts = [];
      this.serviceUpdated();
      this.isActiveClass = this.service.IsActive
        ? "text-success"
        : "text-danger";
      this.serviceFetched = true;
    });
  },
  data() {
    return {
      pageType: "ServiceDiscovery.EditService",
      serviceFetched: false,
      service: {
        Id: "",
        Name: "",
        Hosts: [],
        MatchingURI: "",
        ToURI: "",
        Protected: false,
        APIDocumentation: "",
        IsActive: true,
        HealthcheckUrl: "",
        ServiceManagementHost: "",
        ServiceManagementPort: "",
        ServiceManagementEndpoints: {},
        ProtectedExclude: {}
      },
      isActiveClass: "text-success",
      informationStatus: {
        isActive: false,
        className: "alert-success",
        msg: ""
      },
      cards: {
        basic: {
          showing: true
        },
        api_config: {
          showing: true
        },
        management_config: {
          showing: true
        }
      }
    };
  },
  methods: {
    serviceURL: function() {
      return this.$utils.urlConcat(
        this.$config.API.getApiBaseUrl(),
        this.service.MatchingURI
      );
    },
    copyURL: function() {
      var tempInput = document.createElement("input");
      tempInput.style = "position: absolute; left: -1000px; top: -1000px";
      tempInput.value = this.serviceURL();
      document.body.appendChild(tempInput);
      tempInput.select();
      document.execCommand("copy");
      document.body.removeChild(tempInput);
    },
    addEndpointExclude: function(endpointToExclude) {
      this.service.ProtectedExclude[endpointToExclude.endpoint] =
        endpointToExclude.methods;
    },
    removeEndpointExclude: function(endpointToExclude) {
      var protect = Object.assign({}, this.service.ProtectedExclude);
      delete protect[endpointToExclude];
      this.service.ProtectedExclude = protect;
    },
    addHost: function(hostToAdd) {
      this.service.Hosts.push(hostToAdd);
      hostToAdd = "";
    },
    removeHost: function(hostToRemove) {
      var index = this.service.Hosts.indexOf(hostToRemove);
      this.service.Hosts.splice(index, 1);
    },
    toggleCard: function(cardName) {
      this.cards[cardName].showing = !this.cards[cardName].showing;
    },
    store: function() {
      serviceDiscoveryAPI.updateService(this.service, response => {
        this.informationStatus.isActive = true;

        if (response.status !== 201) {
          this.informationStatus.msg = response.body.msg;
          this.informationStatus.className = "alert-danger";
        } else {
          this.informationStatus.msg = "Resource added successfully.";
          this.informationStatus.className = "alert-success";
        }

        this.serviceUpdated();
      });
    },
    deleteService: function() {
      serviceDiscoveryAPI.deleteService(this.service.Id, response => {
        this.informationStatus.isActive = true;
        if (response.status === 200) {
          this.informationStatus.msg = "Resource removed successfully.";
          this.informationStatus.className = "alert-success";

          setTimeout(() => {
            this.$router.go("/service-discovery/services");
          }, 400);
        } else {
          this.informationStatus.msg = "Error removing resource.";
          this.informationStatus.className = "alert-danger";
        }
      });
    },
    serviceUpdated: function() {
      this.service.IsActive = this.service.IsActive;
      this.$emit("serviceUpdated", this.service);
    }
  },
  components: {
    InformationPanel,
    ServiceAPIConfiguration,
    ServiceManagementConfig
  }
};
</script>
