<template>
    
</template>

<script>
// LineChart.js
import { Pie } from 'vue-chartjs'

export default {
  extends: Pie,
  name:"pie-chart",
  props: ['data', 'options'],
  mounted () {
    this.renderChart(this.data, this.options)
  }
}

</script>