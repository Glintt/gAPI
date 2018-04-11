<template>
<div>
    <div class="row">
        <div class="col-sm">
            <EditService v-on:serviceUpdated="updateDocumentationEndpoint"></EditService>
        </div>   
    </div>

    <hr/> 
    
    <iframe v-if="ServiceDocumentationEndpoint != null" frameborder="0" style="min-height:600px;width:100%;top:0px;left:0px;right:0px;bottom:0px" height="100%" width="100%" :src="ServiceDocumentationEndpoint"></iframe>
    
</div>
   
</template>

<script>
    var serviceDiscoveryAPI = require("@/api/service-discovery");
    
    import EditService from "@/views/Service/EditService";
    
    export default {
        name: "view-service",
        data() {
            return {
                ServiceDocumentationEndpoint : null,
                auth : require("@/auth")
            }
        },
        methods:{
            updateDocumentationEndpoint : function(service){
                this.ServiceDocumentationEndpoint = "http://" + service.Domain + ":" + service.Port + service.APIDocumentation;
            }
        },
        components:{
            EditService
        }
    }
</script>
