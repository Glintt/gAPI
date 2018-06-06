<template>
  <div id="app">
    <NavBar />
    <br />
    <div  class="container-fluid">
    <router-view/>
    </div>
  </div>
</template>
<script>
import NavBar from "@/components/NavBar.vue";

export default {
  components: {
    NavBar
  }
};
</script>