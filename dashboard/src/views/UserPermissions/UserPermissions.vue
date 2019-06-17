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
    return { services: [], permited: [] };
  },
  computed: {
    ...mapGetters("user_permissions", ["permissions"])
  },
  mounted() {
    const username = this.$route.params.username;
    this.get({ Username: username });
    this.fetchServices();
  },
  watch: {
    services: function() {
      this.permitedUpdate();
    }
  },
  methods: {
    ...mapActions("user_permissions", ["update", "get"]),
    permitedUpdate() {
      for (var s in this.services) {
        for (var p in this.permissions) {
          if (this.services[s].Id === this.permissions[p].ServiceId) {
            this.permited.push({
              UserId: this.permissions[p].UserId,
              ServiceId: this.permissions[p].ServiceId,
              Name: this.services[s].Name
            });
          }
        }
      }
      return false;
    },
    addPermited(id, name) {
      this.permited.push({
        UserId: this.permited[0].UserId,
        ServiceId: id,
        Name: name
      });
    },

    removePermited(id, name) {
      var userId = this.permited[0].UserId;
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
      console.log(id);
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
    }
  }
};
</script>

<style>
</style>
