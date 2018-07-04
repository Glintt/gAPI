<template>
<div >
    <router-link to="/service-discovery/apps-groups/create" 
                            v-if="isLoggedIn && loggedInUser && loggedInUser.IsAdmin"
                            class="btn btn-default" href="#"><i class="fas fa-plus text-danger"></i> Add Application Group</router-link>

    <table class="table">
        <thead>
            <tr class="table-secondary" >
                <th scope="col">Name</th>
                <th scope="col"># APIs</th>
                <th scope="col" v-if="isLoggedIn && loggedInUser && loggedInUser.IsAdmin" style="width: 25%">Actions</th>
            </tr>
        </thead>
        <tbody>
            <tr v-for="sg in groups" v-bind:key="sg.Id">
                <td>
                    {{ sg.Name }}
                    <input class="form-control" v-model="sg.Name" v-show="editing == sg.Id"  />
                </td>
                 <td>
                    {{ sg.Services.length }}
                </td>
                <td v-if="isLoggedIn && loggedInUser && loggedInUser.IsAdmin">
                    <button class="btn btn-sm btn-success" @click="editing = editing == sg.Id ? false : sg.Id">Edit</button>
                    <button class="btn btn-sm btn-primary" @click="updateGroup(sg)">Save</button>
                    <button class="btn btn-sm btn-danger" @click="deleteGroup(sg)">Delete</button>
                    <button class="btn btn-sm btn-info" @click="showAPIs(sg)">Show APIs</button>
                    <button class="btn btn-sm btn-warning" @click="findMatches(sg)">Find Matches</button>
                </td>
            </tr>
            <tr>
                <td>Ungrouped APIs</td>
                <td>
                    {{ ungroupedApplications.length }}
                </td>
                <td v-if="isLoggedIn && loggedInUser && loggedInUser.IsAdmin">
                    <button class="btn btn-sm btn-info" @click="showUngroupedAPIs">Show APIs</button>
                </td>
            </tr>
        </tbody>
    </table>

    <div class="row" v-if="services != null && showing=='apis'">
        <div class="col-sm-12">
            <h4>{{ selectedGroup.Name }} - APIs</h4>
            <hr/>

            <table class="table">
                <thead>
                    <tr class="table-secondary" >
                        <th scope="col">Name</th>
                        <th scope="col">gAPI Path</th>
                        <th scope="col">API Documentation</th>
                        <th scope="col">Actions</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="service in services" v-bind:key="service.Id" v-if="service.IsReachable || (service.UseGroupAttributes && service.GroupVisibility ) || loggedInUser">
                        <td>{{ service.Name }}</td>
                        <td>{{ service.MatchingURI }}</td>
                        <td>{{ service.APIDocumentation }}</td>
                        <td style="max-width: 20rem">
                        <router-link :to="'/service-discovery/service?uri='+service.MatchingURI" 
                            data-toggle="tooltip" title="More info" style="margin-right: 1em" >
                            <i class="fas fa-info-circle"></i>
                        </router-link>
                        <button class="btn btn-danger" data-toggle="tooltip" title="Associate" @click="deassociate(service, selectedGroup)">
                            Deassociate From App Group
                        </button>
                        </td>
                    </tr>
                </tbody>
            </table>
           <!--  <ListServices :services="" :isLoggedIn="isLoggedIn" :loggedInUser="loggedInUser"/> -->

        </div>
    </div>

    <div class="row" v-if="possibleMatches.length > 0 && showing=='possibleMatches'">
        <div class="col-sm-12">
            <h4>{{ selectedGroup.Name }} - Possible Matches APIs</h4>
            <hr/>

            <table class="table">
                <thead>
                    <tr class="table-secondary" >
                        <th scope="col">Name</th>
                        <th scope="col">gAPI Path</th>
                        <th scope="col">API Documentation</th>
                        <th scope="col">Actions</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="service in possibleMatches" v-bind:key="service.Id" v-if="service.IsReachable || (service.UseGroupAttributes && service.GroupVisibility ) || loggedInUser">
                        <td>{{ service.Name }}</td>
                        <td>{{ service.MatchingURI }}</td>
                        <td>{{ service.APIDocumentation }}</td>
                        <td style="max-width: 20rem">
                        <router-link :to="'/service-discovery/service?uri='+service.MatchingURI" 
                            data-toggle="tooltip" title="More info" style="margin-right: 1em" >
                            <i class="fas fa-info-circle"></i>
                        </router-link>
                        <button class="btn btn-success" data-toggle="tooltip" title="Associate" @click="associate(service, selectedGroup)">
                            Associate to App Group
                        </button>
                        </td>
                    </tr>
                </tbody>
            </table>
                
        </div>
    </div>

</div>
    
</template>

<script>
import {mapActions, mapGetters} from 'vuex'
import ListServices from "@/components/service-discovery/ListServices"

export default {
    mounted() {
        this.getData()
    },
    computed: {
        isLoggedIn() {
            return this.$oauthUtils.vmA.isLoggedIn();
        },
        ...mapGetters({
            loggedInUser: 'loggedInUser'
        }),
        ...mapGetters('appsGroups', ['groups', 'ungroupedApplications', 'possibleMatches'])
    },
    data() {
        return {
            editing: false,
            services: null,
            selectedGroup: null,
            showing: null
        }
    },
    methods: {
        getData() {
            this.fetchGroups()
            this.listUngroupedApplications()
        },
        ...mapActions('appsGroups', [
            'fetchGroups',
            'updateGroup',
            'deleteGroup',
            'listUngroupedApplications',
            'findPossibleMatches',
        ]),
        showAPIs: function(appGroup) {
            this.$api.serviceDiscovery.applicationGroupById(appGroup.Id, response => {
                this.selectedGroup = response.body
                this.services = response.body.Services
                this.showing='apis'
            })
        },
        showUngroupedAPIs: function() {
            this.selectedGroup = {
                Name: "Ungrouped applications"
            }

            this.services = this.ungroupedApplications
            this.showing='apis'
        },
        findMatches: function(appGroup) {
            this.selectedGroup = appGroup
            this.findPossibleMatches(appGroup)
            this.showing='possibleMatches'
        },
        associate: function(s, g) {
            this.$api.serviceDiscovery.addServiceToAppsGroup( g.Id, s.Id, response => {
                if (response.status == 201) {
                    this.findPossibleMatches(g)
                    this.getData()
                }
            })
        },
        deassociate(s,g) {
            this.$api.serviceDiscovery.deassociateServiceFromAppsGroup( g.Id, s.Id, response => {
                if (response.status == 201) {
                    this.getData()
                    this.showAPIs(g)
                }
            })
        }
    },
    components: {
        ListServices
    }
}
</script>

<style>

</style>
