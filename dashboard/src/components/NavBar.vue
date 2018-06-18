<template>
    <nav class="navbar navbar-expand-lg  navbar-light bg-light">
        <router-link to="/" class="navbar-brand" >
            <img src="/assets/gAPIlogo.png" width="45" height="30" alt="">
        </router-link>
        <div class="collapse navbar-collapse" id="navbarSupportedContent">
            <ul class="navbar-nav mr-auto ">
                <li class="nav-item dropdown">
                    <a href="" class="nav-link dropdown-toggle" role="button" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">Service Discovery</a>
                   
                    <div class="dropdown-menu" aria-labelledby="navbarDropdown">
                        <router-link to="/service-discovery/services/new" class="dropdown-item"
                            v-if="isLoggedIn && loggedInUser && loggedInUser.IsAdmin"><i class="fas fa-plus text-danger"></i> Register New API</router-link>
                        <router-link to="/service-discovery/services" class="dropdown-item" href="#"><i class="fas fa-server text-primary"></i> List APIs</router-link>
                        <router-link to="/service-discovery/groups/create" class="dropdown-item" href="#"><i class="fas fa-plus text-danger"></i> Add Service Group</router-link>
                        <router-link to="/service-discovery/groups" class="dropdown-item" href="#"><i class="fas fa-server text-info"></i> Service Groups</router-link>
                    </div>
                </li>
                <li class="nav-item dropdown" v-if="isLoggedIn && loggedInUser && loggedInUser.IsAdmin">
                    <a href="" class="nav-link dropdown-toggle" role="button" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">Analytics</a>

                    <div class="dropdown-menu" aria-labelledby="navbarDropdown">
                        <router-link to="/analytics/by-api" class="dropdown-item"><i class="fas fa-chart-pie text-info"></i> By API</router-link>
                        <router-link to="/analytics/realtime" class="dropdown-item"><i class="fas fa-chart-area text-warning"></i> Realtime</router-link>
                    </div>
                </li>
                <li class="nav-item dropdown" v-if="isLoggedIn && loggedInUser && loggedInUser.IsAdmin">
                    <a href="" class="nav-link dropdown-toggle" role="button" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">Users Administration</a>

                    <div class="dropdown-menu" aria-labelledby="navbarDropdown">
                        <router-link to="/users" class="dropdown-item"><i class="fas fa-users text-primary"></i> List All</router-link>
                        <router-link to="/users/create" class="dropdown-item"><i class="fas fa-user-plus text-success"></i> Add New</router-link>
                    </div>
                </li>
            </ul>

            <ul class="navbar-nav">
                <li class="nav-item navbar-nav" v-if="!isLoggedIn">
                     <router-link to="/login" class="nav-link">Login</router-link>
                </li>
                
                
                <li class="nav-item dropdown" v-if="isLoggedIn">
                    <a href="" class="nav-link dropdown-toggle" role="button" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
                        {{loggedInUser.Username}}
                    </a>

                    <div class="dropdown-menu">
                        <router-link to="/profile" class="nav-link"><i class="fas fa-user text-info"></i> Profile</router-link>
                        <a href="" @click="logout()" class="nav-link"><i class="fas fa-power-off text-danger"></i> Logout</a>
                    </div>
                </li>
            </ul>
        </div>
    </nav>
</template>

<script>
import { mapGetters } from 'vuex'


export default {
  computed: {
    isLoggedIn() {
      return this.$oauthUtils.vmA.isLoggedIn();
    },
    ...mapGetters({
      loggedInUser: 'loggedInUser'
    })
  },
  methods: {
    logout: function() {
      this.$oauthUtils.vmA.logout();
      this.$router.go("/login");
    }
  }
};
</script>
