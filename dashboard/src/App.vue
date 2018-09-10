<template>
  <div id="app">
    <NavBar v-if="! isFullscreen"/>
    <br />
    <div class="container-fluid" style="margin-top:0px">
    <router-view/>
    </div>
  </div>
</template>
<style>
small,
strong,
p,
body,
span {
  font-family: "Montserrat", sans-serif;
}
.btn {
  margin-left: 0.2em;
  margin-right: 0.2em;
}
.toggable-card {
  cursor: pointer;
}
</style>

<script>
import NavBar from "@/components/NavBar.vue";
import { mapGetters } from "vuex";

export default {
  computed: {
    ...mapGetters("fullscreen", {
      isFullscreen: "isFullscreen"
    })
  },
  components: {
    NavBar
  }
};
</script>
