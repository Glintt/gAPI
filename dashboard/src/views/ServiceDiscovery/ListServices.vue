<template>
    <div class="home">
        <div class="alert alert-dark col-sm-3" role="alert">
            <strong>gAPI Base Url:</strong> 
            {{ 'http://' + $config.API.HOST + ':' + $config.API.PORT }}</br>
            <small>Use this URL + gAPIPath to call microservices</small>
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
                        <button @click="refreshService(service.MatchingURI)" class="btn btn-sm btn-info"  v-show="isLoggedIn">
                            <i class="fas fa-sync"></i>
                        </button>                       
                    </td>
                </tr>
            </tbody>
        </table>
        <ErrorMessage @modalClosed="errorClosed" :showing="error.showing" :id="'requestError'" :error="error.msg" :title="'Error Occurred'"/>
    </div>
</template>

<script>
    var serviceDiscoveryAPI = require("@/api/service-discovery");
    import ErrorMessage from "@/components/ErrorMessage";

    export default {
        name: "home",
        mounted(){
            serviceDiscoveryAPI.listServices((response) => {
                if (response.status!=200) {
                    this.services = [];
                    return;
                }
                this.services = response.body;
            })
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
                }
            }
        },
        methods:{
            refreshService: function(service){
                serviceDiscoveryAPI.refreshService(service, (response) => {
                    if (response.status != 200) {
                        this.error.msg = response.body.msg;
                        this.error.showing = true;
                    }
                })
            },
            errorClosed: function(){
                this.error.showing = false;
                this.error.msg = "";
            }
        },
        components:{
            ErrorMessage
        }
    }

</script>
