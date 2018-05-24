<template>
    <div class="row">
        <div class="col-sm">
            <InformationPanel v-if="informationStatus.isActive" :msg="informationStatus.msg" :className="informationStatus.className"></InformationPanel>
            <div class="row">
            </div>
            <div class="row">
                <div class="col-sm">
                    <h2>gAPI Information</h2>
                    <div class="form-group" v-if="serviceFetched">
                        <span>Health: </span>
                        <i class="fas fa-heartbeat fa-lg" :class="isActiveClass"></i>
                        <small v-if="service.LastActiveTime != 0"> - Last Time Active: {{ this.$utils.convertMillisToTime(new Date().getTime() - service.LastActiveTime) }} </small>
                    </div>
                    <div class="form-group">
                        <label for="serviceName">Microservice Call Base Path</label>
                        <input type="text"
                            :disabled="true"
                             :value="'http://' + $config.API.HOST + ':' + $config.API.PORT + service.MatchingURI" 
                             class="form-control" id="gapiBasePath" aria-describedby="nameHelp" placeholder="Enter name">
                        <small id="nameHelp" class="form-text text-muted">Base Path to call microservice.</small>
                    </div>
                    <div class="form-group">
                        <label for="serviceName">Name</label>
                        <input type="text"
                            :disabled="!isLoggedIn"
                             v-model="service.Name" class="form-control" id="serviceName" aria-describedby="nameHelp" placeholder="Enter name">
                        <small id="nameHelp" class="form-text text-muted">Give the service/API a name.</small>
                    </div>
                    <div class="form-group">
                        <label for="serviceMatchingUri">MatchingURI</label>
                        <input type="text"
                            :disabled="!isLoggedIn"
                             v-model="service.MatchingURI" class="form-control" id="serviceMatchingUri" aria-describedby="serviceMatchingUriHelp" placeholder="Enter domain">
                        <small id="serviceMatchingUriHelp" class="form-text text-muted">Base URI which links to the service on API Management Platform.</small>
                    </div>
                </div>
            </div>
        
            <div class="row" v-show="isLoggedIn">
                <h2>MicroService Info</h2>
            </div>
            <div class="row" v-show="isLoggedIn">
                <div class="form-group col-sm-3">
                    <label for="hostsName">Hosts</label>
                    <input type="text" v-model="hostToAdd" class="form-control" id="hostsName" aria-describedby="hostsHelp" placeholder="Enter hosts">
                    <small id="hostsHelp" class="form-text text-muted">Hosts where the service is hosted.</small>
                    <button type="button" @click="addHost" class="btn btn-sm btn-success">Add</button>
                </div> 
                <ul class="list-group">
                    <li class="list-group-item" v-for="h in service.Hosts" v-bind:key="h">
                        {{ h }}
                        <button type="button" @click="removeHost(h)" class="btn btn-sm btn-danger">Delete</button>
                    </li>
                </ul>
                <div class="form-group col-sm-3">
                    <label for="serviceToUri">URI</label>
                    <input type="text" 
                        :disabled="!isLoggedIn"
                        v-model="service.ToURI" class="form-control" id="serviceToUri" aria-describedby="serviceToUriHelp" placeholder="Enter domain">
                    <small id="serviceToUriHelp" class="form-text text-muted">Service/API Base URI.</small>
                </div>
                <div class="form-group col-sm-3">
                    <label for="serviceDocumentation">Documentation Location</label>
                    <input type="text" 
                        :disabled="!isLoggedIn"
                        v-model="service.APIDocumentation" class="form-control" id="serviceDocumentation" aria-describedby="serviceDocumentationHelp" placeholder="Enter domain">
                    <small id="serviceDocumentationHelp" class="form-text text-muted">API documentation URI.</small>
                </div>
                <div class="form-group col-sm-3">
                    <label for="serviceDocumentation">Healthcheck URL</label>
                    <input type="text" 
                        :disabled="!isLoggedIn"
                        v-model="service.HealthcheckUrl" class="form-control" id="serviceHealthcheckUrl" aria-describedby="serviceHealthcheckUrl" placeholder="Enter Healthcheck Url">
                    <small id="serviceHealthcheckUrl" class="form-text text-muted">Healthcheck URL</small>
                </div>
                <div class="form-check col-sm-3">
                    <input type="checkbox"
                        :disabled="!isLoggedIn"
                            v-model="service.Protected" class="form-check-input" id="serviceProtected">
                    <label class="form-check-label" for="serviceProtected">Protected</label>
                    <small id="serviceProtectedHelp" class="form-text text-muted">Is Service Protected with OAuth?</small>
                </div>
                <div class="form-check col-sm-3">
                    <input type="checkbox" 
                        :disabled="!isLoggedIn"
                        v-model="service.IsCachingActive" class="form-check-input" id="serviceProtected">
                    <label class="form-check-label" for="serviceProtected">Enable Caching?</label>
                    <small id="serviceProtectedHelp" class="form-text text-muted">Enable caching on this service? It will improve performance but be careful as it may affect results.</small>
                </div>

                <div class="form-group col-sm-3">
                    <label for="serviceDocumentation">Service Management Service Host</label>
                    <input type="text" v-model="service.ServiceManagementHost" class="form-control" id="ServiceManagementHost" aria-describedby="ServiceManagementHostHelp" placeholder="Enter service management webservices host">
                    <small id="ServiceManagementeHostHelp" class="form-text text-muted">Host where service management webservices (restart, undeploy, ...) are located at.</small>
                </div>
                <div class="form-group col-sm-3">
                    <label for="serviceDocumentation">Service Management Port</label>
                    <input type="text" v-model="service.ServiceManagementPort" class="form-control" id="ServiceManagementPort" aria-describedby="ServiceManagementPortHelp" placeholder="Enter service management webservices port">
                    <small id="ServiceManagementPortHelp" class="form-text text-muted">Port where service management webservices (restart, undeploy, ...) are located at.</small>
                </div>

                <div class="row">
                    <div class="form-group col-sm-3" v-for="(type, c) in managementTypes" v-bind:key="type.action">
                        <label for="serviceDocumentation">Service {{ type.action }} endpoint</label>
                        <input type="text" v-model="service.ServiceManagementEndpoints[type.action]" class="form-control" :id="type.action + 'ServiceEndpoint'" :aria-describedby="type.action + 'ServiceEndpointHelp'"  v-bind:placeholder="'Enter ' + type.action + ' service endpoint'">
                        <small :id="type.action + 'ServiceEndpointHelp'" class="form-text text-muted">Endpoint to call to {{ type.action }} service.</small>
                    </div>
                </div>
            </div>
            <div class="row" v-show="isLoggedIn">
                <button type="submit" v-if="isLoggedIn" class="btn btn-primary" v-on:click="store">Save</button>
                <button type="submit" v-if="isLoggedIn" class="btn btn-info" v-on:click="serviceUpdated">Preview</button>
                <button class="btn btn-danger" v-if="isLoggedIn" @click="deleteService">Delete</button>
            </div>
        </div>    
    </div>
</template>

<script>
    var serviceDiscoveryAPI = require("@/api/service-discovery");
    import InformationPanel from "@/components/InformationPanel";

    export default {
        name: "view-service",
        computed:{
            isLoggedIn(){
                return this.$oauthUtils.vmA.isLoggedIn();
            }
        },
        mounted() {
            this.$api.serviceDiscovery.manageServiceTypes(response => {
                this.managementTypes = response.body;
            })
            serviceDiscoveryAPI.getServices(this.$route.query.uri, (response) => {
                this.service = response.body;
                if(this.service.ServiceManagementEndpoints == null) this.service.ServiceManagementEndpoints = {};
                if(this.service.Hosts == null) this.service.Hosts = [];
                this.serviceUpdated();
                this.isActiveClass = this.service.IsActive ? 'text-success' : 'text-danger';
                this.serviceFetched = true;
            })
        },
        data() {
            return {
                hostToAdd: "",
                managementTypes:{},
                serviceFetched:false,
                service: {
                    "Name": "",
                    Hosts:[],
                    "MatchingURI": "",
                    "ToURI": "",
                    "Protected": false,
                    "APIDocumentation": "",
                    "IsActive": true,
                    "HealthcheckUrl" : "",
                    ServiceManagementHost : "",
                    ServiceManagementPort : "",
                    ServiceManagementEndpoints:{}
                },
                isActiveClass : 'text-success',
                informationStatus:{
                    isActive : false,
                    className: 'alert-success',
                    msg : ""
                }
            }
        },
        methods: {
            addHost : function() {
                this.service.Hosts.push(this.hostToAdd);
                this.hostToAdd = "";
            },
            removeHost: function(hostToRemove) {
                var index = this.service.Hosts.indexOf(hostToRemove);
                this.service.Hosts.splice(index, 1);
            },
            store: function() {
                serviceDiscoveryAPI.updateService(this.service, (response) => {
                    this.informationStatus.isActive = true;

                    if (response.status != 201) {
                        this.informationStatus.msg = response.body.msg;
                        this.informationStatus.className = 'alert-danger';
                    } else {
                        this.informationStatus.msg = "Resource added successfully.";
                        this.informationStatus.className = 'alert-success';
                    }
                    
                    this.serviceUpdated()
                });
            },
            deleteService : function(){
                serviceDiscoveryAPI.deleteService(this.service.MatchingURI, (response) => {
                    this.informationStatus.isActive = true;
                    if(response.status == 200){
                        this.informationStatus.msg = "Resource removed successfully.";
                        this.informationStatus.className = 'alert-success';

                        setTimeout(() => {
                            this.$router.go("/service-discovery/services");
                        }, 400);
                    } else {
                        this.informationStatus.msg = "Error removing resource.";
                        this.informationStatus.className = 'alert-danger';
                    }
                })
            },
            serviceUpdated:function(){
                this.service.IsActive = this.service.IsActive;
                this.$emit("serviceUpdated", this.service);
            }
        },
        components:{
            InformationPanel
        }
    }
</script>
