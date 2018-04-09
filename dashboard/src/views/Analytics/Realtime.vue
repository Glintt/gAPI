<template>
        <div class="row">
            <div class="col-md-4">
                <LiveRequests></LiveRequests>
            </div>
        </div>
</template>

<script>

    import LiveRequests from "@/views/Analytics/Live/LiveRequests"
    export default {
        name: "analytics-realtime",
        
        components: {
            LiveRequests
        }
    }
</script>