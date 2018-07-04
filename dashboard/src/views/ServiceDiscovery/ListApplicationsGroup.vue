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

    <div class="row" v-if="services != null">
        <div class="col-sm-12">
            <h4>{{ selectedGroup.Name }} - APIs</h4>
            <hr/>
            <ListServices :services="services" :isLoggedIn="isLoggedIn" :loggedInUser="loggedInUser"/>

        </div>
    </div>

    <div class="row" v-if="possibleMatches.length > 0">
        <div class="col-sm-12">
            <h4>{{ selectedGroup.Name }} - Possible Matches APIs</h4>
            <hr/>
            <ListServices :services="possibleMatches" :isLoggedIn="isLoggedIn" :loggedInUser="loggedInUser"/>
        </div>
    </div>

</div>
    
</template>

<script>
import {mapActions, mapGetters} from 'vuex'
import ListServices from "@/components/service-discovery/ListServices"

export default {
    mounted() {
        this.fetchGroups()
        this.listUngroupedApplications()
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
            selectedGroup: null
        }
    },
    methods: {
        ...mapActions('appsGroups', [
            'fetchGroups',
            'updateGroup',
            'deleteGroup',
            'listUngroupedApplications',
            'findPossibleMatches'
        ]),
        showAPIs: function(appGroup) {
            this.$api.serviceDiscovery.applicationGroupById(appGroup.Id, response => {
                this.selectedGroup = response.body
                this.services = response.body.Services
            })
        },
        showUngroupedAPIs: function() {
            this.selectedGroup = {
                Name: "Ungrouped applications"
            }

            this.services = this.ungroupedApplications
        },
        findMatches: function(appGroup) {
            this.selectedGroup = appGroup
            this.findPossibleMatches(appGroup)
        }
    },
    components: {
        ListServices
    }
}
</script>

<style>

</style>
