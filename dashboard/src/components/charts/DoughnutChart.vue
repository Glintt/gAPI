

<script>
import {Doughnut, mixins} from 'vue-chartjs'
const { reactiveProp } = mixins

export default {
	name: 'doughnut-chart',
	extends: Doughnut,
	mixins: [reactiveProp],
	props: ['options'],

	mounted () {
		this.renderChart(this.chartData, this.options)
	}
}
</script>