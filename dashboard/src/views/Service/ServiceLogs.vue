<template>
<div>
    <br />
    <h3 class="text-success">{{ $route.query.uri }} - Logs</h3>

    <div class="card">
        <div class="card-body">
            <div v-if="logsText == ''" class="alert alert-danger" role="alert">
                No logs found.
            </div>
            
            <textarea v-if="logsText != ''" class="form-control" id="exampleFormControlTextarea1" :rows="rows" v-model="logsText" disabled></textarea>
        </div>
    </div>
</div>
</template>

<script>

export default {
    data() {
        return {
            logsText: "",
            rows: 6
        }
    },
    mounted(){
        this.logs();
    },
    methods: {
        logs: function() {
            this.$api.serviceDiscovery.manageService(this.$route.query.uri, this.$api.serviceDiscovery.ManagementActions.logs, (response) => {
                if (response.status != 200) {
                    return;
                }
                this.logsText = response.body.service_response;
                var length = response.body.service_response.split("\n").length;
                this.rows = length < 30 ? length : 30 ;
            });
        }
    }
}
</script>
