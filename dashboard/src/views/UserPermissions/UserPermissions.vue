<template>
  <div>
    <div class="row col-md">
      <h3 class="text-primary">{{$route.params.username}}'s permissions</h3>
    </div>
    <br />
    <div class="row">
      <div class="col-md-6 offset-md-3 card card-body">
        <h5 class="text-success">Allowed</h5>
        <ul class="list-group">
          <li class="list-group-item" v-for="appGroup in permissions" v-bind:key="appGroup.Id">
            {{appGroup.Name}}
            <button class="btn btn-danger" @click="deleteAppGroup(appGroup.Id)">Remove</button>
          </li>
        </ul>
      </div>
    </div>
    <br/>
    <div class="row">
      <div class="col-md-6  offset-md-3 card card-body">
        <h5 class="text-info">Application groups</h5>
        <div class="form-group">
          <select class="form-control" id="appGroupsId" v-model="appToAdd">
            <option :value="appGroup.Id" v-for="appGroup in groups.filter(g => g.Services.length > 0)" v-bind:key="appGroup.Id">{{appGroup.Name}}</option>
          </select>
        </div>
        <button class="btn btn-primary" @click="addNewApp(appToAdd)">Add</button>
      </div>
    </div>
  </div>
</template>

<script>
import { mapActions, mapGetters } from "vuex";
var serviceDiscoveryAPI = require("@/api/service-discovery");

export default {
  props: ["showingUser"],
  data() {
    return { userId: null, appToAdd: null };
  },
  computed: {
    ...mapGetters("user_permissions", ["permissions"]),
    ...mapGetters("users", ["user", "usersList"]),
    ...mapGetters("appsGroups", ["groups"]),
  },
  mounted() {
    const username = this.$route.params.username;
    this.get(username);
    if (this.user === null) {
      this.updateList();
    }
    this.fetchGroups();
  },
  methods: {    
    ...mapActions("appsGroups", [
      "fetchGroups"
    ]),
    ...mapActions("user_permissions", ["add", "get", "remove"]),
    ...mapActions("users", ["updateList", "changeUser"]),
    
    addNewApp(id) {
      this.add({ user: this.$route.params.username,  applicationId: id})
    },

    deleteAppGroup(id) {
      this.remove(
        { user: this.$route.params.username,  applicationId: id}
      )
    }
  }
};
</script>

<style>
</style>
