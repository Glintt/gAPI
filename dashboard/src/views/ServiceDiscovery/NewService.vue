<template>
    <div class="row">
        <div class="col-sm">
            <InformationPanel v-if="informationStatus.isActive" :msg="informationStatus.msg" :className="informationStatus.className"></InformationPanel>
            <form v-on:keyup.13="store">
                <div class="row">

                    <div class="col-sm">
                        <div class="form-group">
                            <label for="serviceName">Name</label>
                            <input type="text" v-model="service.Name" class="form-control" id="serviceName" aria-describedby="nameHelp" placeholder="Enter name">
                            <small id="nameHelp" class="form-text text-muted">Give the service/API a name.</small>
                        </div>
                        <div class="form-group">
                            <label for="domainName">Domain</label>
                            <input type="text" v-model="service.Domain" class="form-control" id="domainName" aria-describedby="domainHelp" placeholder="Enter domain">
                            <small id="domainHelp" class="form-text text-muted">Domain/IP where the service is hosted.</small>
                        </div>
                        <div class="form-group">
                            <label for="servicePort">Port</label>
                            <input type="text" v-model="service.Port" class="form-control" id="servicePort" aria-describedby="servicePortHelp" placeholder="Enter port">
                            <small id="servicePortHelp" class="form-text text-muted">Port where the service is exposed.</small>
                        </div>
                        <div class="form-group">
                            <label for="serviceMatchingUri">MatchingURI</label>
                            <input type="text" v-model="service.MatchingURI" class="form-control" id="serviceMatchingUri" aria-describedby="serviceMatchingUriHelp" placeholder="Enter domain">
                            <small id="serviceMatchingUriHelp" class="form-text text-muted">Base URI which links to the service on API Management Platform.</small>
                        </div>
                    </div>
                    <div class="col-sm">
                        <div class="form-group">
                            <label for="serviceToUri">To URI</label>
                            <input type="text" v-model="service.ToURI" class="form-control" id="serviceToUri" aria-describedby="serviceToUriHelp" placeholder="Enter domain">
                            <small id="serviceToUriHelp" class="form-text text-muted">Service/API Base URI.</small>
                        </div>
                        <div class="form-check">
                            <input type="checkbox" v-model="service.Protected" class="form-check-input" id="serviceProtected">
                            <label class="form-check-label" for="serviceProtected">Protected</label>
                            <small id="serviceProtectedHelp" class="form-text text-muted">Is Service Protected with OAuth?</small>
                        </div>
        
                        <div class="form-check">
                            <input type="checkbox" v-model="service.IsCachingActive" class="form-check-input" id="serviceCaching">
                            <label class="form-check-label" for="serviceCaching">Activate Cache?</label>
                            <small id="serviceCachingHelp" class="form-text text-muted">Is caching enabled for the service?</small>
                        </div>
        
                        <div class="form-group">
                            <label for="serviceDocumentation">Documentation Location</label>
                            <input type="text" v-model="service.APIDocumentation" class="form-control" id="serviceDocumentation" aria-describedby="serviceDocumentationHelp" placeholder="Enter domain">
                            <small id="serviceDocumentationHelp" class="form-text text-muted">API documentation URI.</small>
                        </div>
                    </div>
                </div>
            </form>
    
            <div class="row">
                <div class="col-sm">
                    <button type="submit" class="btn btn-primary" v-on:click="store" >Save</button>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
    var serviceDiscoveryAPI = require("@/api/service-discovery");
    import InformationPanel from "@/components/InformationPanel";

    export default {
        name: "home",
        data() {
            return {
                service: {
                    Name: "",
                    Domain: "",
                    Port: "",
                    MatchingURI: "",
                    ToURI: "",
                    Protected: false,
                    APIDocumentation: "",
                    IsCachingActive : false
                },
                informationStatus:{
                    isActive : false,
                    className: 'alert-success',
                    msg : ""
                }
            }
        },
        methods: {
            store : function(){
                serviceDiscoveryAPI.storeService(this.service, (response) => {
                    if(response.status != 201)
                    {
                        this.informationStatus.msg = response.body.msg;
                        this.informationStatus.isActive = true;
                        this.informationStatus.className = 'alert-danger';
                    }else{
                        this.informationStatus.msg = "Resource added successfully.";
                        this.informationStatus.isActive = true;
                        this.informationStatus.className = 'alert-success';
                    }
                    

                })
            }
        },
        components:{
            InformationPanel
        }
    }
</script>
