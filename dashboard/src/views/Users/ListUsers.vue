<template>
  <div>
    <button class="offset-sm-11 btn btn-sm btn-warning" @click="closeEditing" v-if="user">Go back</button>

    <edit-user v-if="user" :showingUser="user"></edit-user>

    <data-table
      v-if="! user"
      :headers="headers"
      :searchable="true"
      v-on:changePage="changePage"
      v-on:search="search"
      :data="usersList.map(u => { let user = Object.assign({},u); delete user.OAuthClients; return user; })"
      :actions="actions"
      v-on:editUser="editUser"
      v-on:viewPermissions="viewPermissions"
    ></data-table>
  </div>
</template>

<script>
import { mapActions, mapGetters } from "vuex";
import DataTable from "@/components/DataTable";
import EditUser from "./EditUser";

export default {
  computed: {
    ...mapGetters("users", ["usersList", "user"])
  },
  mounted() {
    this.updateList();
    this.changeUser(null);
  },
  data() {
    return {
      headers: ["Id", "Username", "Email", "IsAdmin"],
      actions: {
        edit: {
          name: "Edit",
          event: "editUser"
        },
        permissions: {
          name: "Permissions",
          event: "viewPermissions"
        }
      },
      searchQuery: "",
      page: 1
    };
  },
  methods: {
    ...mapActions("users", ["updateList", "changeUser"]),
    editUser: function(user) {
      user = this.usersList.find(u => u.Id === user.Id)
      console.log(this.usersList.find(u => u.Id === user.Id))
      this.changeUser(user);
    },
    viewPermissions: function(user) {
      this.changeUser(user);
      this.$router.push("/user-permissions/" + user.Username);
    },
    viewPermissions: function(user) {
      this.changeUser(user);
      this.$router.push("/user-permissions/" + user.Username);
    },
    closeEditing: function() {
      this.changeUser(null);
    },
    changePage: function(page) {
      this.page = page;
      this.updateList({ query: this.searchQuery, page: this.page });
    },
    search: function(searchQuery) {
      this.searchQuery = searchQuery;
      this.updateList({ query: this.searchQuery, page: this.page });
    }
  },
  components: {
    DataTable,
    EditUser
  }
};
</script>

<style>
</style>
