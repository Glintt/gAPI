<template>
    <table class="table">
        <thead>
            <tr class="table-secondary" >
                <th scope="col">Name</th>
                <th scope="col">Reachable</th>
                <th scope="col" v-if="isLoggedIn && loggedInUser && loggedInUser.IsAdmin">Actions</th>
            </tr>
        </thead>
        <tbody>
            <tr v-for="sg in groups" v-bind:key="sg.Id">
                <td>
                    {{ sg.Name }}
                    <input class="form-control" v-model="sg.Name" v-show="editing == sg.Id"  />
                </td>
                <td>
                    <i class="fas " @click="sg.IsReachable = !sg.IsReachable" :class="sg.IsReachable ? 'fa-eye text-success' : 'fa-eye-slash text-danger'" v-if="isLoggedIn && loggedInUser && loggedInUser.IsAdmin" />
                    <i class="fas " v-if="!(isLoggedIn && loggedInUser && loggedInUser.IsAdmin)"  :class="sg.IsReachable ? 'fa-eye text-success' : 'fa-eye-slash text-danger'"/>
                </td>
                <td v-if="isLoggedIn && loggedInUser && loggedInUser.IsAdmin">
                    <button class="btn btn-sm btn-success" @click="editing = editing == sg.Id ? false : sg.Id">Edit</button>
                    <button class="btn btn-sm btn-primary" @click="updateGroup(sg)">Save</button>
                </td>
            </tr>
        </tbody>
    </table>
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
        ...mapGetters('serviceGroups', ['groups'])
    },
    data() {
        return {
            editing: false
        }
    },
    methods: {
        ...mapActions('serviceGroups', [
            'fetchGroups',
            'updateGroup'
        ])
    }
}
</script>

<style>

</style>
