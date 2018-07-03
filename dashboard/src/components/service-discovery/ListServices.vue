<template>
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
                <tr v-for="service in services" v-bind:key="service.Id" v-if="service.IsReachable || (service.UseGroupAttributes && service.GroupVisibility ) || loggedInUser">
                    <td>{{ service.Name }}</td>
                    <td>{{ service.MatchingURI }}</td>
                    <td>{{ service.APIDocumentation }}</td>
                    <td><i class="fas fa-heartbeat " :class="service.IsActive ? 'text-success' : 'text-danger'"></i></td>
                    <td v-show="isLoggedIn"><i class="fas " :class="service.Protected ? 'fa-lock text-success' : 'fa-unlock text-danger'"></i></td>
                    <td style="max-width: 20rem">
                      <router-link :to="'/service-discovery/service?uri='+service.MatchingURI" 
                          data-toggle="tooltip" title="More info" style="margin-right: 1em" >
                          <i class="fas fa-info-circle"></i>
                      </router-link>
                      <i class="fas fa-desktop text-success" 
                          data-toggle="tooltip" title="Manage Service" 
                          style="cursor:pointer" @click="showManageModal(service)" v-show="isLoggedIn && loggedInUser.IsAdmin"></i>
                      <!-- <button class="btn btn-success" @click="showManageModal(service)" v-show="isLoggedIn && loggedInUser.IsAdmin">
                        <i class="fas fa-desktop"></i> Manage
                      </button> -->
                    
                      <router-link :to="'/service-discovery/service/logs?uri='+service.MatchingURI" v-show="isLoggedIn"
                          data-toggle="tooltip" title="View service logs" style="margin-left: 1em">
                          <i class="fas fa-file text-warning"></i>
                      </router-link>
                    </td>
                </tr>
            </tbody>
        </table>
</template>

<script>
export default {
    name: 'list-services',
    props: ['services', 'isLoggedIn', 'loggedInUser'],
    methods: {
        showManageModal: function(service) {
            this.$emit('showManageModal', service)
        }
    }
}
</script>

<style>

</style>
