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
            <div class="col-md">
                <div class="card text-white bg-danger mb-3">
                    <div class="card-body text-center">
                        <h5 class="card-title">{{ globalInformation.totalErrors }}</h5>
                        <p class="card-text">Total Errors</p>
                    </div>
                </div>
            </div>
            <div class="col-md">
                <div class="card text-white bg-success mb-3">
                    <div class="card-body text-center">
                        <h5 class="card-title">{{ globalInformation.totalSuccess }}</h5>
                        <p class="card-text">Total Success</p>
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
    
        <div class="row">
            <div class="col-md">
                <div class="card ">
                    <div class="card-body text-center">
                        <h5 class="card-title">Most Called APIs</h5>
                        <hr/>
                        <div class="card-text">
                            <doughnut-chart :chartData="mostCalledApis" :options="{responsive: false}"></doughnut-chart>
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
                            <bar-chart :chartData="globalUserAgent" :options="{responsive: false}"></bar-chart>
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
                            <bar-chart :chartData="globalRemoteAddress" :options="{responsive: false}"></bar-chart>
                        </div>
                    </div>
                </div>
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
        props: ["selectedAPI"],
        watch:{
            selectedAPI : function(){
                this.GetData()
            }
        },
        mounted() {
            this.GetData()
        },
        data() {
            return {
                fullAnalytics:  {},
                elapsedTimeInfo: {},
                mostCalledApis: {},
                globalInformation: {},
                globalRemoteAddress: {},
                globalUserAgent:{},
            }
        },

        methods: {
            GetData:function(){
                analyticsAPI.byApi({
                    "endpoint": this.selectedAPI
                }, (response) => {
                    this.fullAnalytics = response.body.aggregations.api.buckets;
                    this.UpdateAllInfo()
                });
            },
            ResetData: function() {
                this.elapsedTimeInfo = {
                    max: 0,
                    min: -1,
                    avg: 0
                };
            },
            UpdateAllInfo: function() {
                this.ResetData()
                this.ElapsedTimeInfo()
                this.mostCalledApis = this.MostCalledApis()
                this.globalInformation = this.GlobalInformation()
                this.globalUserAgent = this.GlobalUserAgent()
                this.globalRemoteAddress = this.GlobalRemoteAddress()
            },
            ElapsedTimeInfo: function() {
                var totalTime = 0;
                var size = this.fullAnalytics.length;
    
                for (var i = 0; i < size; i++) {
                    this.elapsedTimeInfo.max = this.fullAnalytics[i].MaxElapsedTime.value > this.elapsedTimeInfo.max ? this.fullAnalytics[i].MaxElapsedTime.value : this.elapsedTimeInfo.max;
                    this.elapsedTimeInfo.min = (this.fullAnalytics[i].MinElapsedTime.value < this.elapsedTimeInfo.min || this.elapsedTimeInfo.min == -1) ? this.fullAnalytics[i].MinElapsedTime.value : this.elapsedTimeInfo.min;
                    totalTime += this.fullAnalytics[i].AvgElapsedTime.value;
                }
    
                this.elapsedTimeInfo.avg = this.$utils.timeRound(totalTime / size);
    
                return this.elapsedTimeInfo;
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