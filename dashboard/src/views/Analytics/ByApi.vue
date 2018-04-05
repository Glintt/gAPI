<template>
    <div class="container-fluid">
        <div class="row">
            <div class="col-md">
                <div class="form-group">
                    <label for="sel1">Filter by API:</label>
                    <select class="form-control" v-model="selectedAPI" id="sel1">
                            <option v-for="api in apisList">{{api}}</option>
                        </select>
                    <button class="btn btn-success" @click="ApplyFilter()">Apply</button>
                    <button class="btn btn-warning" v-if="selectedAPI!=null" @click="RemoveFilter()">Remove Filter</button>
                </div>
            </div>
        </div>
        <div class="row">
            <div class="col-md">
                <div class="card text-white bg-primary mb-3">
                    <div class="card-body text-center">
                        <h5 class="card-title">{{ fullAnalytics.length}}</h5>
                        <p class="card-text">Total called APIs</p>
                    </div>
                </div>
            </div>
            <div class="col-md">
                <div class="card text-white bg-danger mb-3">
                    <div class="card-body text-center">
                        <h5 class="card-title">{{ GlobalInformation().totalErrors }}</h5>
                        <p class="card-text">Total Errors</p>
                    </div>
                </div>
            </div>
            <div class="col-md">
                <div class="card text-white bg-success mb-3">
                    <div class="card-body text-center">
                        <h5 class="card-title">{{ GlobalInformation().totalSuccess }}</h5>
                        <p class="card-text">Total Success</p>
                    </div>
                </div>
            </div>
        </div>
    
        <div class="row">
            <div class="col-md">
                <div class="card ">
                    <div class="card-body text-center">
                        <h5 class="card-title">Most Called APIs</h5>
                        <hr/>
                        <div class="card-text">
                            <doughnut-chart :chartData="MostCalledApis()" :options="{responsive: false}"></doughnut-chart>
                        </div>
                    </div>
                </div>
            </div>
            <div class="col-md">
                <div class="card">
                    <div class="card-body text-center">
                        <h5 class="card-title">User Agent</h5>
                        <hr/>
                        <div class="card-text">
                            <bar-chart :chartData="GlobalUserAgent()" :options="{responsive: false}"></bar-chart>
                        </div>
                    </div>
                </div>
            </div>
            <div class="col-md">
                <div class="card">
                    <div class="card-body text-center">
                        <h5 class="card-title">Remote Address</h5>
                        <hr/>
                        <div class="card-text">
                            <bar-chart :chartData="GlobalRemoteAddress()" :options="{responsive: false}"></bar-chart>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    
        <div class="row">
            <div class="col-md">
                <data-table :headers="table.headers" :data="table.data"></data-table>
                
            </div>
        </div>
    </div>
</template>

<script>
    var analyticsAPI = require("@/api/analytics");
    
    import DoughnutChart from "@/components/charts/DoughnutChart"
    import BarChart from "@/components/charts/BarChart"
    import DataTable from "@/components/DataTable"
    
    export default {
        name: "analytics-by-api",
        mounted() {
            analyticsAPI.byApi({}, (response) => {
                this.fullAnalytics = response.body.aggregations.api.buckets;
                this.apisList = this.ApisList();
                this.logs = response.body.hits.hits;
                this.UpdateTableData();
            });
        },
        data() {
            return {
                table: {
                    headers:["ID","Service", "Host", "UserAgent", "Remote IP", "Method","Date time","Query Arguments", "Headers"],
                    data: []
                },
                
                fullAnalytics:  {},
                selectedAPI: null,
                apisList: [],
                logs : []
            }
        },
        methods: {
            UpdateTableData:function(){
              
                for(var i = 0; i < this.logs.length; i++){
                    var log = this.logs[i];
                    var obj = {
                        "ID": log._id,
                        "Service":log._source.Service,
                        "Host":log._source.Host,
                        "UserAgent":log._source.UserAgent,
                        "RemoteIP":log._source.RemoteIp,
                        "Method":log._source.Method,
                        "DateTime":log._source.DateTime,
                        "QueryArgs":log._source.QueryArgs,
                        "Headers":log._source.Headers
                    };
                    this.table.data.push(obj)
                }  
            },
            RemoveFilter: function() {
                analyticsAPI.byApi({}, (response) => {
                    this.fullAnalytics = response.body.aggregations.api.buckets;
                    this.selectedAPI = null;
                    this.logs = response.body.hits.hits;
                    this.UpdateTableData();
                })
            },
            ApplyFilter: function() {
                analyticsAPI.byApi({
                    "endpoint": this.selectedAPI
                }, (response) => {
                    this.fullAnalytics = response.body.aggregations.api.buckets;

                    this.logs = response.body.hits.hits;
                    this.UpdateTableData();
                })
            },
            ApisList: function() {
                var apis = []
                for (var i = 0; i < this.fullAnalytics.length; i++) {
                    apis.push(this.fullAnalytics[i].key);
                }
                return apis;
            },
            MostCalledApis: function() {
                var calledApis = {
                    labels: [],
                    datasets: [{
                        label: "Most called APIs",
                        backgroundColor: [],
                        data: []
                    }]
                }
    
                for (var i = 0; i < this.fullAnalytics.length && i < 3; i++) {
                    var key = this.fullAnalytics[i].key;
                    var value = this.fullAnalytics[i].doc_count;
    
                    calledApis.labels.push(key);
                    calledApis.datasets[0].data.push(value);
                    calledApis.datasets[0].backgroundColor.push(this.$chartColors.colors[this.$random.randomBetween(0, 4)]);
                }
                return calledApis;
            },
    
            GlobalInformation: function() {
                var globalInfo = {
                    totalErrors: 0,
                    totalSuccess: 0
                }
    
                for (var i = 0; i < this.fullAnalytics.length; i++) {
                    var statusCodeForEndpoint = this.fullAnalytics[i].StatusCode.buckets;
    
                    for (var j = 0; j < statusCodeForEndpoint.length; j++) {
                        var key = statusCodeForEndpoint[j].key;
                        var value = statusCodeForEndpoint[j].doc_count;
    
                        if (key >= 400) {
                            globalInfo.totalErrors = globalInfo.totalErrors + value;
                        }
                        if (key < 400) {
                            globalInfo.totalSuccess = globalInfo.totalSuccess + value;
                        }
                    }
                }
    
                return globalInfo;
            },
            GlobalUserAgent: function() {
                var userAgentApis = {
                    labels: [],
                    datasets: [{
                        label: "User Agent APIs",
                        backgroundColor: [],
                        data: []
                    }]
                };
    
                var userAgent = {};
    
                for (var i = 0; i < this.fullAnalytics.length; i++) {
                    var userAgentForEndpoint = this.fullAnalytics[i].UserAgent.buckets;
    
                    for (var j = 0; j < userAgentForEndpoint.length; j++) {
                        var key = userAgentForEndpoint[j].key;
                        var value = userAgentForEndpoint[j].doc_count;
    
                        if (userAgent.hasOwnProperty(key)) {
                            userAgent[key] = userAgent[key] + value;
                        } else {
                            userAgent[key] = value;
                        }
                    }
                }
    
                for (var prop in userAgent) {
                    if (userAgentApis.labels.length >= 6)
                        return userAgentApis;
                    userAgentApis.labels.push(prop);
                    userAgentApis.datasets[0].data.push(userAgent[prop]);
                    userAgentApis.datasets[0].backgroundColor.push(this.$chartColors.colors[this.$random.randomBetween(0, this.$chartColors.colors.length)]);
                }
                return userAgentApis;
            },
            GlobalRemoteAddress: function() {
                var remoteAddrApis = {
                    labels: [],
                    datasets: [{
                        label: "Remote Addr APIs",
                        backgroundColor: [],
                        data: []
                    }]
                };
    
                var remoteAddr = {};
    
                for (var i = 0; i < this.fullAnalytics.length; i++) {
                    var remoteAddrForEndpoint = this.fullAnalytics[i].RemoteAddr.buckets;
    
                    for (var j = 0; j < remoteAddrForEndpoint.length; j++) {
                        var key = remoteAddrForEndpoint[j].key;
                        var value = remoteAddrForEndpoint[j].doc_count;
    
                        if (remoteAddr.hasOwnProperty(key)) {
                            remoteAddr[key] = remoteAddr[key] + value;
                        } else {
                            remoteAddr[key] = value;
                        }
                    }
                }
    
                for (var prop in remoteAddr) {
                    if (remoteAddrApis.labels.length >= 6)
                        return remoteAddrApis;
                    remoteAddrApis.labels.push(prop);
                    remoteAddrApis.datasets[0].data.push(remoteAddr[prop]);
                    remoteAddrApis.datasets[0].backgroundColor.push(this.$chartColors.colors[this.$random.randomBetween(0, this.$chartColors.colors.length)]);
                }
                return remoteAddrApis;
            }
        },
        components: {
            DoughnutChart,
            BarChart,
            DataTable
        }
    }
</script>
