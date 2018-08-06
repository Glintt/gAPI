<template>
    <div class="card mb-6">
        <div class="card-header text-white bg-success" @click="toggleCard('management_config')">Management configuration</div>
        <div class="card-body" v-if="showing">
            <div class="row">
                <div class="form-group col-sm-4 text-center a offset-sm-4">
                    <h6>
                        <i class="fas fa-desktop"></i> Service Management Service
                    </h6>
                </div>
            </div>
            <div class="row">
                <div class="form-group col-sm-3 offset-sm-3">
                    <label for="serviceDocumentation">
                        Host
                    </label>
                    <input :disabled="! isAdmin" type="text" v-model="service.ServiceManagementHost" class="form-control" id="ServiceManagementHost" aria-describedby="ServiceManagementHostHelp" placeholder="Service management webservices host">
                    <small id="ServiceManagementeHostHelp" class="form-text text-success">Host where service management webservices (restart, undeploy, ...) are located at.</small>
                </div>
                <div class="form-group col-sm-3">
                    <label for="serviceDocumentation">Port</label>
                    <input :disabled="! isAdmin" type="text" v-model="service.ServiceManagementPort" class="form-control" id="ServiceManagementPort" aria-describedby="ServiceManagementPortHelp" placeholder="Service management webservices port">
                    <small id="ServiceManagementPortHelp" class="form-text text-success">Port where service management webservices (restart, undeploy, ...) are located at.</small>
                </div>
            </div>
            <hr />
            <div class="row">
                <div class="form-group col-sm-3" v-for="(type, c) in managementTypes" v-bind:key="c">
                    <label for="serviceDocumentation">Service {{ type.action }} endpoint</label>
                    <input :disabled="! isAdmin" type="text" v-model="service.ServiceManagementEndpoints[type.action]" class="form-control" :id="type.action + 'ServiceEndpoint'" :aria-describedby="type.action + 'ServiceEndpointHelp'"  v-bind:placeholder="'Enter ' + type.action + ' service endpoint'">
                    <small :id="type.action + 'ServiceEndpointHelp'" class="form-text text-success">Endpoint to call to {{ type.action }} service.</small>
                </div>
            </div>
        </div>
    </div>
</template>


<script>
import { mapGetters } from 'vuex'
export default {
	name: 'service-management-config',
	props: ['service', 'showing'],
	mounted () {
		this.$api.serviceDiscovery.manageServiceTypes(response => {
			this.managementTypes = response.body
		})
	},
	computed: {
		...mapGetters({
			isAdmin: 'isAdmin',
			loggedInUser: 'loggedInUser'
		})
	},
	data () {
		return {
			managementTypes: {}
		}
	},
	methods: {
		toggleCard: function (cardName) {
			this.$emit('toggleCard', cardName)
		}
	}
}
</script>