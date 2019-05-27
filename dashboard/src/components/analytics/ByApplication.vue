<template>
    <div>
      <div class="row">
          <div class="col-md">
              <div class="card text-white bg-primary mb-3">
                  <div class="card-body text-center">
                      <h5 class="card-title">{{ fullAnalytics.length}}</h5>
                      <p class="card-text">Total called APIs</p>
                  </div>
              </div>
          </div>
      </div>
      <div class="row">
          <div class="col-md">
              <div class="card text-white bg-info mb-3">
                  <div class="card-body text-center">
                      <h5 class="card-title">{{ elapsedTimeInfo.avg }}</h5>
                      <p class="card-text">Average Elapsed Time (ms)</p>
                  </div>
              </div>
          </div>
          <div class="col-md">
              <div class="card text-white bg-danger mb-3">
                  <div class="card-body text-center">
                      <h5 class="card-title">{{ elapsedTimeInfo.max }}</h5>
                      <p class="card-text">Max Elapsed Time (ms)</p>
                  </div>
              </div>
          </div>
          <div class="col-md">
              <div class="card text-white bg-success mb-3">
                  <div class="card-body text-center">
                      <h5 class="card-title">{{ elapsedTimeInfo.min }}</h5>
                      <p class="card-text">Min Elapsed Time (ms)</p>
                  </div>
              </div>
          </div>
      </div>
    
  </div>
</template>

<script>
import DoughnutChart from "@/components/charts/DoughnutChart";
import BarChart from "@/components/charts/BarChart";
import DataTable from "@/components/DataTable";

var analyticsAPI = require("@/api/analytics");

export default {
  name: "analytics-by-application",
  props: ["selectedAPI"],
  watch: {
    selectedAPI: function() {
      this.GetData();
    }
  },
  mounted() {
    this.GetData();
  },
  data() {
    return {
      fullAnalytics: {},
      elapsedTimeInfo: {},
      mostCalledApis: {},
      globalInformation: {},
      globalRemoteAddress: {},
      globalUserAgent: {}
    };
  },

  methods: {
    GetData: function() {
      analyticsAPI.byApplication(
        {
          app_id: this.selectedAPI
        },
        response => {
          this.fullAnalytics = response.body.aggregations;
          this.fullAnalytics.length = response.body.hits.total;
          this.UpdateAllInfo();
        }
      );
    },
    ResetData: function() {
      this.elapsedTimeInfo = {
        max: 0,
        min: -1,
        avg: 0
      };
    },
    UpdateAllInfo: function() {
      this.ResetData();
      this.ElapsedTimeInfo();
    },
    ElapsedTimeInfo: function() {
      this.elapsedTimeInfo.max =this.fullAnalytics.MaxElapsedTime.value
      this.elapsedTimeInfo.min = this.fullAnalytics.MinElapsedTime.value;
      this.elapsedTimeInfo.avg = this.fullAnalytics.AvgElapsedTime.value;
    },
  },
  components: {
    DoughnutChart,
    BarChart,
    DataTable
  }
};
</script>
