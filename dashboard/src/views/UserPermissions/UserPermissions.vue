<template>
  <div>
    <div class="row">
      <h3>{{$route.params.username}}'s permissions</h3>
      <button
        class="btn btn-sm btn-info"
        @click="update({
          newPermissions:permited,
          user: {
              Username: $route.params.username
          }
      })"
      >Update</button>
    </div>
    <div class="row">
      <div class="col">
        <div class="form-check" v-for="service in permited" v-bind:key="service.ServiceId">
          <input
            type="checkbox"
            class="form-check-input"
            id="removePermitedcheck"
            checked
            @click="removePermited(service.ServiceId,service.Name)"
          >
          <label class="form-check-label" for="removePermitedcheck">{{ service.Name }}</label>
        </div>
      </div>
      <div class="col">
        <div class="form-check" v-for="service in services" v-bind:key="service.Identifier">
          <input
            type="checkbox"
            class="form-check-input"
            id="addPermitedcheck"
            @click="addPermited(service.Id,service.Name)"
          >
          <label class="form-check-label" for="addPermitedcheck">{{ service.Name }}</label>
        </div>
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
    return { services: [], permited: [], userId: null };
  },
  computed: {
    ...mapGetters("user_permissions", ["permissions"]),
    ...mapGetters("users", ["user", "usersList"])
  },
  mounted() {
    const username = this.$route.params.username;
    this.get({ Username: username });
    this.fetchServices();
    if (this.user === null) {
      this.updateList();
    }
  },
  watch: {
    services: function() {
      this.permitedUpdate();
    }
  },
  methods: {
    ...mapActions("user_permissions", ["update", "get"]),
    ...mapActions("users", ["updateList", "changeUser"]),
    permitedUpdate() {
      if (this.user === null) this.updateUser();
      for (var s in this.services) {
        for (var p in this.permissions) {
          if (this.services[s].Id === this.permissions[p].ServiceId) {
            this.permited.push({
              UserId: this.user.Id,
              ServiceId: this.permissions[p].ServiceId,
              Name: this.services[s].Name
            });
          }
        }
      }
      return false;
    },
    addPermited(id, name) {
      if (this.user === null) this.updateUser();
      this.permited.push({
        UserId: this.user.Id,
        ServiceId: id,
        Name: name
      });
    },

    removePermited(id, name) {
      if (this.user === null) this.updateUser();
      var userId = this.user.Id;
      var temp = [];
      for (var i = 0; i < this.permited.length; i++) {
        if (this.permited[i].ServiceId !== id) {
          temp.push({
            UserId: userId,
            ServiceId: this.permited[i].ServiceId,
            Name: this.permited[i].Name
          });
        }
      }
      this.permited = temp;
    },

    fetchServices() {
      serviceDiscoveryAPI.listServices(-1, "", response => {
        if (response.status !== 200) {
          this.services = [];
          return;
        }
        this.services = response.body;
      });
    },

    updateUser() {
      var user = null;
      for (var i = 0; i < this.usersList.length; i++) {
        if (this.usersList[i].Username === this.$route.params.username)
          user = this.usersList[i];
      }
      this.changeUser(user);
    }
  }
};
</script>

<style>
</style>
