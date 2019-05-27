<template>
    <div class="container-fluid">
        <div class="row">
            <div class="col-md">
                <div class="form-group">
                    <label for="sel1">Filter by Application:</label>
                    <select class="form-control" v-model="selectedAPI" id="sel1">
                        <option v-for="api in applicationsList" v-bind:key="api">{{api}}</option>
                    </select>
                    <button class="btn btn-warning" v-if="selectedAPI!==null" @click="RemoveFilter()">Remove Filter</button>
                </div>
            </div>
        </div>
        <ByApplication :selectedAPI="selectedAPI"/> 
    </div>
</template>

<script>
import ByApplication from "@/components/analytics/ByApplication";

var analyticsAPI = require("@/api/analytics");

export default {
  name: "analytics-by-applicatino",
  mounted() {
    analyticsAPI.byApplication({}, response => {
      this.fullAnalytics = response.body.aggregations.api.buckets;
      this.applicationsList = this.GetApplicationsList();
    });
  },
  data() {
    return {
      fullAnalytics: {},
      selectedAPI: null,
      applicationsList: []
    };
  },
  methods: {
    RemoveFilter: function() {
      analyticsAPI.byApplication({}, response => {
        this.fullAnalytics = response.body.aggregations.api.buckets;
        this.selectedAPI = null;
        this.UpdateTableData();
      });
    },
    GetApplicationsList: function() {
      var apis = [];
      for (var i = 0; i < this.fullAnalytics.length; i++) {
        apis.push(this.fullAnalytics[i].key);
      }
      return apis;
    }
  },
  components: {
    
    ByApplication
  }
};
</script>
