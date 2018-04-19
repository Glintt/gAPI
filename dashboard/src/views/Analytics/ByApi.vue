<template>
    <div class="container-fluid">
        <div class="row">
            <div class="col-md">
                <div class="form-group">
                    <label for="sel1">Filter by API:</label>
                    <select class="form-control" v-model="selectedAPI" id="sel1">
                        <option v-for="api in apisList">{{api}}</option>
                    </select>
                    <button class="btn btn-warning" v-if="selectedAPI!=null" @click="RemoveFilter()">Remove Filter</button>
                </div>
            </div>
        </div>
        <ByAPI :selectedAPI="selectedAPI"/>

        <br />
        <h3 class="title text-success">Errors List:</h3>
        <Logs :selectedAPI="selectedAPI" />
    </div>
</template>

<script>
    var analyticsAPI = require("@/api/analytics");
    
    import Logs from "@/components/analytics/Logs"
    import ByAPI from "@/components/analytics/ByAPI"
    
    export default {
        name: "analytics-by-api",
        mounted() {
            analyticsAPI.byApi({}, (response) => {
                this.fullAnalytics = response.body.aggregations.api.buckets;
                this.apisList = this.ApisList();
            });
        },
        data() {
            return {
                fullAnalytics:  {},
                selectedAPI: null,
                apisList: [],
            }
        },
        methods: {
            RemoveFilter: function() {
                analyticsAPI.byApi({}, (response) => {
                    this.fullAnalytics = response.body.aggregations.api.buckets;
                    this.selectedAPI = null;
                    this.UpdateTableData();
                })
            },
            ApisList: function() {
                var apis = []
                for (var i = 0; i < this.fullAnalytics.length; i++) {
                    apis.push(this.fullAnalytics[i].key);
                }
                return apis;
            }
        },
        components: {
            Logs,
            ByAPI
        }
    }
</script>
