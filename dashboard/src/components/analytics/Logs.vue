<template>
    <div class="row">
        <div class="col-md">            
            <data-table :headers="table.headers" :data="table.data" :actions="table.actions" v-on:copyPostman="copyPostman"></data-table>
        </div>
    </div>
</template>

<script>
    import DoughnutChart from '@/components/charts/DoughnutChart'
import BarChart from '@/components/charts/BarChart'
import DataTable from '@/components/DataTable'

var analyticsAPI = require('@/api/analytics')

export default {
    	name: 'analytics-logs',
    	props: ['selectedAPI'],
    	mounted () {
    		this.GetData()
    	},
    	watch: {
    		selectedAPI: function () {
    			this.ApplyFilter()
    		}
    	},
    	data () {
    		return {
    			table: {
    				headers: ['ID', 'Service', 'Host', 'Uri', 'UserAgent', 'Remote IP', 'Method', 'Date time', 'Query Arguments', 'Headers', 'Body', 'Response'],
    				data: [],
    				actions: {
    					postman: {
    						name: 'POSTMAN Copy',
    						event: 'copyPostman'
    					}
    				}
    			},
    
    			logs: []
    		}
    	},
    	methods: {
    
    		GetData: function () {
    			analyticsAPI.logs({
    				'endpoint': this.selectedAPI
    			}, (response) => {
    				this.logs = response.body.hits.hits
    				this.UpdateTableData()
    			})
		},
    		UpdateTableData: function () {
    			this.table.data = []
    			for (var i = 0; i < this.logs.length; i++) {
    				var log = this.logs[i]
    				var obj = {
    					'ID': log._id,
    					'Service': log._source.Service,
    					'Host': log._source.Host,
    					'Uri': log._source.Uri,
    					'UserAgent': log._source.UserAgent,
    					'RemoteIP': log._source.RemoteIp,
    					'Method': log._source.Method,
    					'DateTime': log._source.DateTime,
    					'QueryArgs': log._source.QueryArgs,
    					'Headers': log._source.Headers,
    					'Body': log._source.RequestBody,
    					'Response': log._source.Response
    				}

				this.table.data.push(obj)
    			}
    		},
    
    		ApplyFilter: function () {
    			this.GetData()
    		},
    		copyPostman: function (data) {
    			var tempInput = document.createElement('input')
			tempInput.style = 'position: absolute; left: -1000px; top: -1000px'
			tempInput.value = this.postmanRequestInfo(data)
    			document.body.appendChild(tempInput)
    			tempInput.select()
    			document.execCommand('copy')
			document.body.removeChild(tempInput)
    		},
    
    		postmanRequestInfo: function (data) {
    			var requestInfo = `curl -X ${data.Method} `

    			// add headers
    			var headers = JSON.parse(data.Headers)
    			for (var headerName in headers) {
    				requestInfo += `-H "${headerName}: ${headers[headerName]}" `
    			}
    
    			// add body
    			if (data.Body !== '') {
    				requestInfo += `-d '${data.Body}' `
    			}

    			// add host
    			requestInfo += data.Host + data.Uri + '?'

			// add query params
			var queryParams = JSON.parse(data.QueryArgs)
    			var anyParamAdded = false
    			for (var arg in queryParams) {
    				if (anyParamAdded) requestInfo += '&'
				requestInfo += arg + '=' + queryParams[arg]
    				anyParamAdded = true
    			}

    			return requestInfo
    		}
    	},
    	components: {
    		DoughnutChart,
    		BarChart,
    		DataTable
    	}
    }
</script>
