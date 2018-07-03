<template>
<div >
   <router-link to="/service-discovery/apps-groups/create" 
                            v-if="isLoggedIn && loggedInUser && loggedInUser.IsAdmin"
                             class="btn btn-default" href="#"><i class="fas fa-plus text-danger"></i> Add Application Group</router-link>
                        
<table class="table">
        <thead>
            <tr class="table-secondary" >
                <th scope="col">Name</th>
                <th scope="col" v-if="isLoggedIn && loggedInUser && loggedInUser.IsAdmin">Actions</th>
            </tr>
        </thead>
        <tbody>
            <tr v-for="sg in groups" v-bind:key="sg.Id">
                <td>
                    {{ sg.Name }}
                    <input class="form-control" v-model="sg.Name" v-show="editing == sg.Id"  />
                </td>
                <td v-if="isLoggedIn && loggedInUser && loggedInUser.IsAdmin">
                    <button class="btn btn-sm btn-success" @click="editing = editing == sg.Id ? false : sg.Id">Edit</button>
                    <button class="btn btn-sm btn-primary" @click="updateGroup(sg)">Save</button>
                    <button class="btn btn-sm btn-danger" @click="deleteGroup(sg)">Delete</button>
                    
                </td>
            </tr>
        </tbody>
    </table>

</div>
    
</template>

<script>
import {mapActions, mapGetters} from 'vuex'

export default {
    mounted() {
        this.fetchGroups()
    },
    computed: {
        isLoggedIn() {
            return this.$oauthUtils.vmA.isLoggedIn();
        },
        ...mapGetters({
            loggedInUser: 'loggedInUser'
        }),
        ...mapGetters('appsGroups', ['groups'])
    },
    data() {
        return {
            editing: false
        }
    },
    methods: {
        ...mapActions('appsGroups', [
            'fetchGroups',
            'updateGroup',
            'deleteGroup'
        ])
    }
}
</script>

<style>

</style>
