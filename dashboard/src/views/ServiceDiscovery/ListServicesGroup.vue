<template>
    <table class="table">
        <thead>
            <tr class="table-secondary" >
                <th scope="col">Name</th>
                <th scope="col">Reachable</th>
                <th scope="col">Actions</th>
            </tr>
        </thead>
        <tbody>
            <tr v-for="sg in groups" v-bind:key="sg.Id">
                <td>
                    {{ sg.Name }}
                    <input class="form-control" v-model="sg.Name" v-show="editing" />
                </td>
                <td>
                    <i class="fas " :class="sg.IsReachable ? 'fa-eye text-success' : 'fa-eye-slash text-danger'" @click="sg.IsReachable = !sg.IsReachable"/>
                </td>
                <td>
                    <button class="btn btn-sm btn-success" @click="editing=!editing">Edit</button>
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
