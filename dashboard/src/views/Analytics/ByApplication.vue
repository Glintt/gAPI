<template>
    <div class="container-fluid">
        <div class="row">
            <div class="col-md">
                <div class="form-group">
                    <label for="sel1">Filter by Application:</label>
                    <select class="form-control" v-model="selectedAPI" id="sel1">
                        <option v-for="api in groups" v-bind:key="api.Id">{{api.Name}}</option>
                    </select>
                    <button class="btn btn-warning" v-if="selectedAPI!==null" @click="RemoveFilter()">Remove Filter</button>
                </div>
            </div>
        </div>
        <ByApplication :selectedAPI="selectedAPI" v-if="selectedAPI != null"/> 
        <span v-if="selectedAPI == null">Select an application to get analytics</span> 
    </div>
</template>

<script>
import ByApplication from "@/components/analytics/ByApplication";
import { mapActions, mapGetters } from "vuex";

var analyticsAPI = require("@/api/analytics");

export default {
  name: "analytics-by-application",
  mounted() {
    this.fetchGroups();
  },
  computed:{
    ...mapGetters("appsGroups", [
      "groups"
    ])
  },
  data() {
    return {
      fullAnalytics: {},
      selectedAPI: null,
    };
  },
  methods: {    
    ...mapActions("appsGroups", [
      "fetchGroups"
    ]),
    RemoveFilter: function() {
      this.selectedAPI = null
      this.UpdateTableData();
      this.fullAnalytics = {}
    }
  },
  components: {    
    ByApplication
  }
};
</script>
