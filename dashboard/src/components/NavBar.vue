<template>
    <nav class="navbar navbar-expand-lg  navbar-light bg-light">
        <router-link to="/" class="navbar-brand" >
            <img src="/assets/gAPIlogo.png" width="45" height="30" alt="">
        </router-link>
        <div class="collapse navbar-collapse" id="navbarSupportedContent">
            <ul class="navbar-nav mr-auto ">
                <li class="nav-item active">
                    
                </li>
                <li class="nav-item" v-if="!isLoggedIn">
                     <router-link to="/login" class="nav-link">Login</router-link>
                </li>
                <li class="nav-item" v-if="isLoggedIn">
                     <a @click="logout()" class="nav-link">Logout</a>
                </li>
                <li class="nav-item dropdown">
                    <a href="" class="nav-link dropdown-toggle" role="button" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">Service Discovery</a>
                   
                    <div class="dropdown-menu" aria-labelledby="navbarDropdown">
                        <router-link to="/service-discovery/services/new" class="dropdown-item"
                            v-if="isLoggedIn && loggedInUser && loggedInUser.IsAdmin">Register New API</router-link>
                        <router-link to="/service-discovery/services" class="dropdown-item" href="#">List APIs</router-link>
                    </div>
                </li>
                <li class="nav-item dropdown" v-if="isLoggedIn && loggedInUser && loggedInUser.IsAdmin">
                    <a href="" class="nav-link dropdown-toggle" role="button" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">Analytics</a>

                    <div class="dropdown-menu" aria-labelledby="navbarDropdown">
                        <router-link to="/analytics/by-api" class="dropdown-item">By API</router-link>
                        <router-link to="/analytics/realtime" class="dropdown-item">Realtime</router-link>
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
