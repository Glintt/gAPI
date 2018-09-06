<template>
    <div v-if="active != null && plugins != null">
        <table class="table" v-for="(row, index) in plugins" v-bind:key="index">
            <thead>
                <tr>
                <th scope="col">Plugin Name</th>
                <th scope="col">Plugin Type</th>
                <th scope="col">IsActive</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="(item, i) in row" v-bind:key="i">
                    <th scope="row">{{ item }}</th>
                    <th scope="row">{{ index }}</th>
                    <th scope="row" v-if="active.BeforeRequest != null">
                        <i class="fas fa-check text-success" v-if="active.BeforeRequest.indexOf(item) !== -1"></i>
                        <i class="fas fa-times text-danger" v-if="active.BeforeRequest.indexOf(item) === -1"></i>
                    </th>
                </tr>
            </tbody>
        </table>
    </div>
    
</template>

<script>
import { mapGetters, mapActions } from "vuex";

export default {
  computed: {
    ...mapGetters("plugins", ["plugins", "active"])
  },
  mounted() {
    this.updatePlugins();
  },
  methods: {
    ...mapActions("plugins", ["updatePlugins"])
  }
};
</script>
