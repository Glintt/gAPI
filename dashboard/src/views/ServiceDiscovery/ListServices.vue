<template>
    <div class="home">
        <table class="table">
            <thead>
                <tr>
                    <th scope="col">Name</th>
                    <th scope="col">Domain</th>
                    <th scope="col">Port</th>
                    <th scope="col">gAPI Path</th>
                    <th scope="col">Secured?</th>
                    <th scope="col">Health</th>
                    <th scope="col">API Documentation</th>
                    <th scope="col">Actions</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="service in services">
                    <td>{{ service.Name }}</td>
                    <td>{{ service.Domain }}</td>
                    <td>{{ service.Port }}</td>
                    <td>{{ service.MatchingURI }}</td>
                    <td><i :class="service.Protected ? 'fas fa-lock text-success' : 'fas fa-unlock text-danger'"></i></td>
                    <td><i :class="service.IsActive ? 'fas fa-check-circle text-success' : 'fas fa-times-circle text-danger'"></i></td>
                    <td>{{ service.APIDocumentation }}</td>
                    <td>
                        <router-link :to="'/service-discovery/service?uri='+service.MatchingURI" class="navbar-brand" >View</router-link>
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
                this.services = response.body;
            })
        },
        data() {
            return {
                services : []
            }
        }
    }

</script>
