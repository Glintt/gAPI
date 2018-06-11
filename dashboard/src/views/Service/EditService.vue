<template>
    <div class="row">
        <div class="col-sm">
            <InformationPanel v-if="informationStatus.isActive" :msg="informationStatus.msg" :className="informationStatus.className"></InformationPanel>
            <div class="row">
            </div>
            <div class="row">
                <div class="col-sm">
                    <div class="row">
                        <div class="col-sm-10" v-show="isLoggedIn">
                            <h2>{{ service.Name }} API Information</h2>
                        </div>
                        <div class="col-sm-2" v-show="isAdmin">
                            <button type="submit" class="btn btn-xs btn-primary" v-on:click="store">Save</button>
                            <button class="btn btn-danger"  @click="deleteService">Delete</button>
                            <button type="submit" class="btn btn-info" v-on:click="serviceUpdated">Preview</button>
                        </div>
                    </div>

                    <div class="card mb-12">
                        <div class="card-header text-white bg-primary" @click="toggleCard('basic')">
                            <div class="row">
                                <div :class="service.LastActiveTime != 0 ? 'col-sm-10' : 'col-sm-11'">
                                    Basic Information
                                </div>
                                <div :class="service.LastActiveTime != 0 ? 'col-sm-2' : 'col-sm-1'">
                                    <span>Health: </span>
                                    <i class="fas fa-heartbeat fa-lg" :class="isActiveClass"></i>
                                    <small v-if="service.LastActiveTime != 0"><br />Last Time Active: {{ this.$utils.convertMillisToTime(new Date().getTime() - service.LastActiveTime) }} </small>
                                </div>
                            </div>
                        </div>
                        <div class="card-body row" v-if="cards.basic.showing">
                            <div class="form-group col-sm">
                                <label for="serviceName">Name</label>
                                <input type="text"
                                    :disabled="! (isLoggedIn && isAdmin)"
                                    v-model="service.Name" class="form-control" id="serviceName" aria-describedby="nameHelp" placeholder="Enter name">
                                <small id="nameHelp" class="form-text text-primary">Give the service/API a name.</small>
                            </div>

                            <div class="form-group col-sm">
                                <label for="serviceMatchingUri">MatchingURI</label>
                                <input type="text"
                                    :disabled="! $permissions.HasPermission(pageType, loggedInUser)"
                                    v-model="service.MatchingURI" class="form-control" id="serviceMatchingUri" aria-describedby="serviceMatchingUriHelp" placeholder="Enter domain">
                                <small id="serviceMatchingUriHelp" class="form-text text-primary">Base URI which links to the service on API Management Platform.</small>
                            </div>

                            <div class="form-group col-sm">
                                <label for="serviceName">URL:</label>
                                <input type="text"
                                    :disabled="true"
                                    :value="'http://' + $config.API.HOST + ':' + $config.API.PORT + service.MatchingURI" 
                                    class="form-control" id="gapiBasePath" aria-describedby="nameHelp" placeholder="Enter name">
                                <small id="nameHelp" class="form-text text-primary">Base Path to call microservice.</small>
                            </div>
                        </div>
                    </div> 

                    <ServiceAPIConfiguration v-on:addHost="addHost" v-show="isLoggedIn" v-on:toggleCard="toggleCard" v-on:removeHost="removeHost" :showing="cards.api_config.showing" :service="service"/>
                    
                    <ServiceManagementConfig v-on:toggleCard="toggleCard" :showing="cards.management_config.showing" :service="service" v-show="isLoggedIn" />
                </div>
            </div>
        </div>    
    </div>
</template>

<script>
var serviceDiscoveryAPI = require("@/api/service-discovery");
import InformationPanel from "@/components/InformationPanel";
import ServiceAPIConfiguration from "@/views/Service/ServiceAPIConfiguration";
import ServiceManagementConfig from "@/views/Service/ServiceManagementConfig";
import { mapGetters } from 'vuex'

export default {
  name: "view-service",
  computed: {
    isLoggedIn() {
      return this.$oauthUtils.vmA.isLoggedIn();
    },
    ...mapGetters({
      isAdmin: 'isAdmin',
      loggedInUser: 'loggedInUser'
    })
  },
  mounted() {
    serviceDiscoveryAPI.getServices(this.$route.query.uri, response => {
      this.service = response.body;
      if (this.service.ServiceManagementEndpoints == null)
        this.service.ServiceManagementEndpoints = {};
      if (this.service.Hosts == null) this.service.Hosts = [];
      this.serviceUpdated();
      this.isActiveClass = this.service.IsActive
        ? "text-success"
        : "text-danger";
      this.serviceFetched = true;
    });

  },
  data() {

    return {
      pageType: 'ServiceDiscovery.EditService',
      serviceFetched: false,
      service: {
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
        ServiceManagementEndpoints: {}
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
          showing: false
        },
        management_config: {
          showing: false
        }
      }
    };
  },
  methods: {
    addHost: function(hostToAdd) {
      this.service.Hosts.push(hostToAdd);
      hostToAdd = "";
    },
    removeHost: function(hostToRemove) {
      var index = this.service.Hosts.indexOf(hostToRemove);
      this.service.Hosts.splice(index, 1);
    },
    toggleCard: function(cardName) {
        console.log(cardName)
      this.cards[cardName].showing = !this.cards[cardName].showing;
    },
    store: function() {
      serviceDiscoveryAPI.updateService(this.service, response => {
        this.informationStatus.isActive = true;

        if (response.status != 201) {
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
      serviceDiscoveryAPI.deleteService(this.service.MatchingURI, response => {
        this.informationStatus.isActive = true;
        if (response.status == 200) {
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
