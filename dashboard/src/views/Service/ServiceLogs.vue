<template>
<div>
    <br />
    <h3 class="text-success">{{ $route.query.uri }} - Logs</h3>

    <div class="card">
        <div class="card-body">
            <div v-if="logsText === '' && !loading" class="alert alert-danger" role="alert">
                No logs found.
            </div>
            <div v-if="loading" class="alert alert-warning" role="alert">
                Loading logs... Wait a moment
            </div>
            
            <pre v-if="logsText !== ''" class="form-control prelog" id="exampleFormControlTextarea1" :rows="rows"  disabled>{{logsText}}</pre>
        </div>
    </div>
</div>
</template>

<style scoped>
.prelog {
  height: auto;
  max-height: 600px;
  overflow: auto;
  background-color: #eeeeee;
  word-break: normal !important;
  word-wrap: normal !important;
  white-space: pre !important;
}
</style>

<script>
export default {
  data() {
    return {
      logsText: "",
      rows: 6,
      loading: true
    };
  },
  mounted() {
    this.logs();
  },
  methods: {
    logs: function() {
      this.loading = true;
      this.$api.serviceDiscovery.manageService(
        this.$route.query.uri,
        "logs",
        response => {
          this.loading = false;
          if (response.status !== 200) {
            return;
          }
          this.logsText = response.body.service_response;
          var length = response.body.service_response.split("\n").length;
          this.rows = length < 30 ? length : 30;
        }
      );
    }
  }
};
</script>
