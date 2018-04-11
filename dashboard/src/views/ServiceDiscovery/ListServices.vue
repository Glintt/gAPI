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
                    <th scope="col" v-show="auth.isLoggedIn()">Secured?</th>
                    <th scope="col" v-show="auth.isLoggedIn()">Health</th>
                    <th scope="col">Actions</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="service in services" v-bind:key="service.Name">
                    <td>{{ service.Name }}</td>
                    <td>{{ service.MatchingURI }}</td>
                    <td>{{ service.APIDocumentation }}</td>
                    <td v-show="auth.isLoggedIn()"><i class="fas " :class="service.Protected ? 'fa-lock text-success' : 'fa-unlock text-danger'"></i></td>
                    <td v-show="auth.isLoggedIn()"><i class="fas fa-heartbeat " :class="service.IsActive ? 'text-success' : 'text-danger'"></i></td>
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
    var serviceDiscoveryAPI = require("@/api/service-discovery");

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
        data() {
            return {
                services : [],
                auth : require("@/auth")
            }
        }
    }

</script>
