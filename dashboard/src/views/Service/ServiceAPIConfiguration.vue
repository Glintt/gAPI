<template>

    <div class="card mb-6">
        <div class="card-header text-white bg-info" @click="toggleCard('api_config')">API configuration</div>
        <div class="card-body" v-if="showing">
            <div class="row col-sm">
                <div class="form-group col-sm-6">
                    <label for="hostsName">Hosts</label>
                    <input type="text" v-model="hostToAdd" class="form-control" id="hostsName" aria-describedby="hostsHelp" placeholder="Enter hosts">
                    <small id="hostsHelp" class="form-text text-info">Hosts where the service is hosted.</small>
                    <button type="button" @click="addHost" class="btn btn-sm btn-success">Add</button>
                    <ul class="list-group">
                        <li class="list-group-item" v-for="h in service.Hosts" v-bind:key="h">
                            {{ h }}
                            <button type="button" @click="removeHost(h)" class="btn btn-sm btn-danger">Delete</button>
                        </li>
                    </ul>
                </div>
                <div class="form-group col-sm-6">
                    <label for="serviceToUri">URI</label>
                    <input type="text" v-model="service.ToURI" class="form-control" id="serviceToUri" aria-describedby="serviceToUriHelp" placeholder="Enter domain">
                    <small id="serviceToUriHelp" class="form-text text-info">Service/API Base URI.</small>
                </div>
                <div class="form-group col-sm-6">
                    <label for="serviceDocumentation">Documentation Location</label>
                    <input type="text" 
                        v-model="service.APIDocumentation" class="form-control" id="serviceDocumentation" aria-describedby="serviceDocumentationHelp" placeholder="Enter domain">
                    <small id="serviceDocumentationHelp" class="form-text text-info">API documentation URI.</small>
                </div>
                <div class="form-group col-sm-6">
                    <label for="serviceDocumentation">Healthcheck URL</label>
                    <input type="text" 
                        v-model="service.HealthcheckUrl" class="form-control" id="serviceHealthcheckUrl" aria-describedby="serviceHealthcheckUrl" placeholder="Enter Healthcheck Url">
                    <small id="serviceHealthcheckUrl" class="form-text text-info">Healthcheck URL</small>
                </div>
            </div>
            <div class="row col-sm">
                <div class="form-group col-sm-6">
                    <i class="fas " :class="service.Protected ? 'fa-lock text-success' : 'fa-unlock text-danger'" @click="service.Protected = !service.Protected" />
                    <!-- <input type="checkbox"
                        :disabled="!isLoggedIn"
                            v-model="service.Protected" class="form-check-input" id="serviceProtected"> -->
                    <label class="form-check-label" for="serviceProtected">&nbsp;&nbsp;Protection</label>
                    <small id="serviceProtectedHelp" class="form-text text-info">Is Service Protected with OAuth?</small>
                </div>
                <div class="form-group col-sm-6">                                    
                    <i class="fas fa-archive " :class="service.IsCachingActive ? 'text-success' : 'text-danger'" @click="service.IsCachingActive = !service.IsCachingActive" />
                    <label class="form-check-label" for="serviceProtected">&nbsp;&nbsp;Cache</label>
                    <small id="serviceProtectedHelp" class="form-text text-info">Enable caching on this service? It will improve performance but be careful as it may affect results.</small> 
                </div>
            </div>
        </div>
    </div>
</template>

<script>
export default {
  name: "service-api-configuration",
  props: ["iLoggedIn", "service", "showing"],
  data() {
    return {
      hostToAdd: ""
    };
  },
  methods: {
    removeHost: function(hostToRemove) {
      this.$emit("removeHost", hostToRemove);
    },
    addHost: function() {
      this.$emit("addHost", this.hostToAdd);
      this.hostToAdd = ""
    },
    toggleCard: function(cardName) {
        this.$emit("toggleCard", cardName)
    }
  }
};
</script>