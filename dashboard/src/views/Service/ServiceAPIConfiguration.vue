<template>

    <div class="card mb-6">
        <div class="card-header text-white bg-info" @click="toggleCard('api_config')">API configuration</div>
        <div class="card-body" v-if="showing">
            <div class="row col-sm"> 
                <div class="form-group col-sm-6" v-if="isAdmin">
                    Associated Group
                    <select class="form-control" v-model="selectedGroup" @change="changed = selectedGroup == service.GroupId ? false : true">
                        <option :value="null"></option>                        
                        <option v-for="group in groups" :value="group.Id" :key="group.Id">{{ group.Name }}</option>
                    </select>
                    <button class="btn btn-sm btn-success" v-if="selectedGroup != null && changed" @click="associateToGroup">Associate</button>
                    <button class="btn btn-sm btn-danger" v-if="selectedGroup != null" @click="deassociateFromGroup">Deassociate</button>
                </div>
                <div class="form-group col-sm-6" v-if="!isAdmin">
                    Associated Group
                    <select class="form-control" v-model="selectedGroup" disabled="true">
                        <option :value="null"></option>                        
                        <option v-for="group in groups" :value="group.Id" :key="group.Id">{{ group.Name }}</option>
                    </select>
                </div>
            </div>
            <div class="row col-sm">
                <div class="form-group col-sm-6">
                    <label for="hostsName">Hosts</label>
                    <input type="text"  v-if="isAdmin"  v-model="hostToAdd" class="form-control" id="hostsName" aria-describedby="hostsHelp" placeholder="Enter hosts" :disabled="! isAdmin">
                    <small id="hostsHelp" class="form-text text-info">Hosts where the service is hosted.</small>
                    <button type="button" v-if="isAdmin"  @click="addHost" class="btn btn-sm btn-success">Add</button>
                    <ul class="list-group">
                        <li class="list-group-item" v-for="h in service.Hosts" :key="h">
                            {{ h }}
                            <button v-if="isAdmin" type="button" @click="removeHost(h)" class="btn btn-sm btn-danger">Delete</button>
                        </li>
                    </ul>
                </div>
                <div class="form-group col-sm-6">
                    <label for="serviceMatchingURI">Matching URI</label>
                    <input  :disabled="! isAdmin" type="text" v-model="service.MatchingURI" class="form-control" id="serviceMatchingUri" aria-describedby="serviceMatchingURIHelp" placeholder="Enter API Matching URI">
                    <small id="serviceMatchingURIHelp" class="form-text text-info">URI which will link to the API.</small>
                </div>
                <div class="form-group col-sm-6">
                    <label for="serviceToUri">URI</label>
                    <input  :disabled="! isAdmin" type="text" v-model="service.ToURI" class="form-control" id="serviceToUri" aria-describedby="serviceToUriHelp" placeholder="Enter API base uri">
                    <small id="serviceToUriHelp" class="form-text text-info">Service/API Base URI.</small>
                </div>
                <div class="form-group col-sm-6">
                    <label for="serviceDocumentation">Documentation Location</label>
                    <input  :disabled="! isAdmin" type="text" 
                        v-model="service.APIDocumentation" class="form-control" id="serviceDocumentation" aria-describedby="serviceDocumentationHelp" placeholder="Enter domain">
                    <small id="serviceDocumentationHelp" class="form-text text-info">API documentation URI.</small>
                </div>
                <div class="form-group col-sm-6">
                    <label for="serviceDocumentation">Healthcheck URL</label>
                    <input  :disabled="! isAdmin" type="text" 
                        v-model="service.HealthcheckUrl" class="form-control" id="serviceHealthcheckUrl" aria-describedby="serviceHealthcheckUrl" placeholder="Enter Healthcheck Url">
                    <small id="serviceHealthcheckUrl" class="form-text text-info">Healthcheck URL</small>
                </div>
            </div>
            <div class="row col-sm">
                <div class="form-group col-sm-6">
                    <i class="fas " :class="service.Protected ? 'fa-lock text-success' : 'fa-unlock text-danger'" @click="toggleProtection" />
                    <!-- <input type="checkbox"
                        :disabled="!isLoggedIn"
                            v-model="service.Protected" class="form-check-input" id="serviceProtected"> -->
                    <label class="form-check-label" for="serviceProtected">&nbsp;&nbsp;Protection</label>
                    <small id="serviceProtectedHelp" class="form-text text-info">Is Service Protected with OAuth?</small>
                </div>
                <div class="form-group col-sm-6">                                    
                    <i class="fas fa-archive " :class="service.IsCachingActive ? 'text-success' : 'text-danger'" @click="toggleCaching" />
                    <label class="form-check-label" for="serviceProtected">&nbsp;&nbsp;Cache</label>
                    <small id="serviceProtectedHelp" class="form-text text-info">Enable caching on this service? It will improve performance but be careful as it may affect results.</small> 
                </div>
                <div class="form-group col-sm-6">                      
                    <i class="fas " :class="service.IsReachable ? 'fa-eye text-success' : 'fa-eye-slash text-danger'" @click="toggleReachable" />
                    <label class="form-check-label" for="serviceReachable">&nbsp;&nbsp;Reachable</label>
                    <small id="serviceReachableHelp" class="form-text text-info">Is service reachable from external sources? If use group attributes is set to true, the it is used Group Reachability</small> 
                </div>
                <div class="form-group col-sm-6">                      
                    <i class="fas " :class="service.UseGroupAttributes ? 'fa-check text-success' : 'fa-times text-danger'" @click="toggleUseGroup" />
                    <label class="form-check-label" for="serviceUseGroup">&nbsp;&nbsp;Use Group Attributes</label>
                    <small id="serviceUseGroupHelp" class="form-text text-info">Use group attributes (reachability, ...)</small> 
                </div>
            </div>
            <hr />
            <h5>Group Info <small class="text-info">This section is not editable </small></h5>
            <br />
            <div class="row col-sm">                
                <div class="form-group col-sm-6">                      
                    <i class="fas " :class="service.GroupVisibility ? 'fa-eye text-success' : 'fa-eye-slash text-danger'" />
                    <label class="form-check-label" for="serviceReachable">&nbsp;&nbsp;Group Reachability</label>
                    <small id="serviceReachableHelp" class="form-text text-info">Group to which service is associated reachability from external. This value is used when Reachable status is set to false</small> 
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import { mapGetters } from 'vuex'

export default {
  name: "service-api-configuration",
  props: ["iLoggedIn", "service", "showing"],
  data() {
    return {
      hostToAdd: "",
      groups: [],
      selectedGroup: null,
      changed : false
    };
  },
  computed: {
    ...mapGetters({
      isAdmin: 'isAdmin',
      loggedInUser: 'loggedInUser'
    })
  },
  mounted() {
    this.$api.serviceDiscovery.listServiceGroups(response => {
        this.groups = response.body;
    });
  },
  watch:{
      service: function() {
          this.selectedGroup = this.service.GroupId == "" ? null : this.service.GroupId
      }
  },
  methods: {
    associateToGroup: function() {
        this.$api.serviceDiscovery.addServiceToServiceGroup(this.selectedGroup, this.service.Id, response => {
            if (response.status == 201) {
                this.service.GroupId = this.selectedGroup
                this.service.GroupVisibility = this.groups.find(element => {
                    return element.Id == this.selectedGroup
                }).IsReachable
            }
        })
    },
    deassociateFromGroup: function() {
        this.$api.serviceDiscovery.deassociateServiceFromServiceGroup(this.selectedGroup, this.service.Id, response => {
            if (response.status == 201) {
                this.selectedGroup = null
                this.service.GroupId = null
                this.service.GroupVisibility = this.service.IsReachable
            }
        })
    },
    removeHost: function(hostToRemove) {
      this.$emit("removeHost", hostToRemove);
    },
    addHost: function() {
      this.$emit("addHost", this.hostToAdd);
      this.hostToAdd = ""
    },
    toggleCard: function(cardName) {
        this.$emit("toggleCard", cardName)
    },
    toggleCaching: function() {
        if (! this.isAdmin) return 
        this.service.IsCachingActive = !this.service.IsCachingActive
    },
    toggleProtection: function() {
        if (! this.isAdmin) return 
        this.service.Protected = !this.service.Protected
    },
    toggleReachable: function() {
        if (! this.isAdmin) return
        this.service.IsReachable = !this.service.IsReachable
    },
    toggleUseGroup: function() {
        if (! this.isAdmin) return
        this.service.UseGroupAttributes = !this.service.UseGroupAttributes
    }
  }
};
</script>