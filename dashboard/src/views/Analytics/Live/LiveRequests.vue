<template>
    <div class="card">
        <div class="card-body text-center">
            <h5 class="card-title">Requests per Second</h5>
            <hr/>
            <div class="card-text">            
                <line-chart :chartData="perSecondData" :options="options"></line-chart>
            </div>
        </div>
    </div>
</template>

<script>
import LineChart from "@/components/charts/LineChart";
import io from "socket.io-client";

export default {
  name: "live-requests",
  mounted() {
    var socket = io.connect(
      "http://" +
        this.$config.API.SOCKET_HOST +
        ":" +
        this.$config.API.SOCKET_PORT
    );
    socket.on("logs", msg => {
      this.updateMonitorInfo(msg);
    });

    socket.on("connect", () => {
      console.log("Connected to server!");
    });

    this.perSecondMonitor();
  },
  data() {
    return {
      updated: false,
      perSecondData: {
        labels: [],
        datasets: [
          {
            label: "Requests per Second",
            backgroundColor: this.$chartColors.colors[3],
            data: []
          }
        ]
      },
      options: { responsive: true, maintainAspectRatio: false },
      lastReceived: 0
    };
  },

  methods: {
    updateMonitorInfo: function(msgs) {
      this.lastReceived = parseInt(msgs);
    },

    perSecondMonitor: function() {
      var time = this.$utils.currentTimeString();

      this.perSecondData.labels.push(time);
      this.perSecondData.datasets[0].data.push(this.lastReceived);
      this.perSecondData = this.updateProps(this.perSecondData);

      this.lastReceived = 0;

      this.shiftData();
      setTimeout(this.perSecondMonitor, 1000);
    },

    shiftData: function() {
      if (this.perSecondData.datasets[0].data.length > 60) {
        this.perSecondData.labels.shift();
        this.perSecondData.datasets[0].data.shift();
      }
    },

    updateProps: function(data) {
      return {
        labels: data.labels,
        datasets: data.datasets
      };
    }
  },
  components: {
    LineChart
  }
};
</script>

<style scoped>
</style>
