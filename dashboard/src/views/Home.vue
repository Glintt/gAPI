<template>
  <div class="home">
  <div class="jumbotron">
  <h1 class="display-4">Welcome to gAPI !</h1>
  <p class="lead">
    This API management dashboard allows you to quickly publish and manage an API.
  </p>
  <hr class="my-4">
  <p>Functionalities:</p>
  <ul>
    <li>Per API Analytics</li>
    <li>Centralized documentation</li>
    <li>Automatic logging</li>
    <li>OAuth Authentication</li>
  </ul>
  <p class="lead">
  <router-link to="/service-discovery/services/new" class="btn btn-primary btn-lg">Start Creating an API</router-link>
  </p>
</div>
  </div>
</template>

<script>
// @ is an alias to /src
import NavBar from '@/components/NavBar.vue'

export default {
	name: 'home',
	components: {
		NavBar
	}
}
</script>
