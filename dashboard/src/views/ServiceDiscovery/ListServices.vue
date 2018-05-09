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
                <tr v-for="service in services">
                    <td>{{ service.Name }}</td>
                    <td>{{ service.MatchingURI }}</td>
                    <td>{{ service.APIDocumentation }}</td>
                    <td><i class="fas fa-heartbeat " :class="service.IsActive ? 'text-success' : 'text-danger'"></i></td>
                    <td v-show="isLoggedIn"><i class="fas " :class="service.Protected ? 'fa-lock text-success' : 'fa-unlock text-danger'"></i></td>
                    <td>
                        <router-link :to="'/service-discovery/service?uri='+service.MatchingURI" class="navbar-brand" >
                            <i class="fas fa-info-circle"></i>
                        </router-link>
                        <button @click="refreshService(service)" class="btn btn-sm btn-info"  v-show="isLoggedIn">
                            <i class="fas fa-sync"></i>
                        </button>
                    </td>
                </tr>
            </tbody>
        </table>
        <ErrorMessage @modalClosed="errorClosed" :showing="error.showing" :id="'requestError'" :error="error.msg" :title="'Error Occurred'"/>
        <ConfirmationModal @answerReceived="restartConfirmationReceived" @modalClosed="confirmationClosed" :showing="confirmation.showing" :id="'refreshConfirm'" :msg="confirmation.msg" :title="'Refresh service'"/>
    </div>
</template>

<script>
    var serviceDiscoveryAPI = require("@/api/service-discovery");
    import DataTable from "@/components/DataTable";
    import ErrorMessage from "@/components/ErrorMessage";
    import ConfirmationModal from "@/components/ConfirmationModal";

    export default {
        name: "home",
        mounted(){
            this.updateData();
        },
        watch: {
            currentPage: function() {
                this.updateData();
            }
        },
        computed:{
            isLoggedIn(){
                return this.$oauthUtils.vmA.isLoggedIn()
            }
        },
        data() {
            return {
                services : [],
                error:{
                    msg: "",
                    showing: false
                },
                restart: {
                    service: null
                },
                confirmation: {
                    showing:false,
                    msg: ""
                },
                currentPage: 1,
                searchText: ""
            }
        },
        methods:{
            refreshService: function(service){
                this.confirmation.showing = true;
                this.confirmation.msg = "Are you sure you want to restart service " + service.Name + "?";
                this.restart.service = service;
            },
            restartConfirmationReceived: function(answer) {
                if (answer == false) return;

                serviceDiscoveryAPI.refreshService(this.restart.service.MatchingURI , (response) => {
                    if (response.status != 200) {
                        this.error.msg = response.body.msg;
                        this.error.showing = true;
                    }
                })
                this.restart.service = null;
            },
            confirmationClosed: function(){
                this.confirmation.showing = false;
                this.confirmation.msg = "";
            },
            errorClosed: function(){
                this.error.showing = false;
                this.error.msg = "";
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
                });
            }
        },
        components:{
            ErrorMessage,
            ConfirmationModal,
            DataTable
        }
    };
</script>
