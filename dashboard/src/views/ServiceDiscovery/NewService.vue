<template>
    <div class="row">
        <div class="col-sm">
            <InformationPanel v-if="informationStatus.isActive" :msg="informationStatus.msg" :className="informationStatus.className"></InformationPanel>
            <h2>Add new API Service</h2>
            
            <form v-on:keyup.13="store">
                <div class="row">
                    <div class="col-sm-4">
                        <div class="form-group">
                            <label for="serviceName" class="text-info">Name</label>
                            <input type="text" v-model="service.Name" class="form-control" id="serviceName" aria-describedby="nameHelp" placeholder="Enter name">
                            <small id="nameHelp" class="form-text text-muted">Give the service/API a name.</small>
                        </div>
                    </div>
                </div>

                <ServiceAPIConfiguration v-on:addEndpointExclude="addEndpointExclude" v-on:removeEndpointExclude="removeEndpointExclude" v-on:addHost="addHost" v-on:toggleCard="toggleCard" v-on:removeHost="removeHost" :showing="cards.api_config.showing" :service="service"/>
                   
                <ServiceManagementConfig v-on:toggleCard="toggleCard" :showing="cards.management_config.showing" :service="service"/>         
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
    import InformationPanel from '@/components/InformationPanel'
import ServiceAPIConfiguration from '@/views/Service/ServiceAPIConfiguration'
import ServiceManagementConfig from '@/views/Service/ServiceManagementConfig'

export default {
    	name: 'new-service',
    	data () {
    		return {
    			hostToAdd: '',
    			service: {
    				Name: '',
    				Hosts: [],
    				MatchingURI: '',
    				ToURI: '',
    				Protected: false,
    				APIDocumentation: '',
    				IsCachingActive: false,
    				HealthcheckUrl: '',
    				ServiceManagementHost: '',
    				ServiceManagementPort: '',
    				ServiceManagementEndpoints: {},
    				IsReachable: false,
    				UseGroupAttributes: false,
    				ProtectedExclude: {}
    			},
    			informationStatus: {
    				isActive: false,
    				className: 'alert-success',
    				msg: ''
    			},
    			cards: {
    				basic: {
    					showing: true
    				},
    				api_config: {
    					showing: false
    				},
    				management_config: {
    					showing: false
    				}
    			}
    		}
    	},
    	methods: {
    		addEndpointExclude: function (endpointToExclude) {
    			this.service.ProtectedExclude[endpointToExclude.endpoint] = endpointToExclude.methods
    		},
    		removeEndpointExclude: function (endpointToExclude) {
    			var protect = Object.assign({}, this.service.ProtectedExclude)
    			delete protect[endpointToExclude]
    			this.service.ProtectedExclude = protect
    		},
    		addHost: function (hostToAdd) {
    			this.service.Hosts.push(hostToAdd)
    			hostToAdd = ''
		},
    		removeHost: function (hostToRemove) {
    			var index = this.service.Hosts.indexOf(hostToRemove)
    			this.service.Hosts.splice(index, 1)
    		},
    		toggleCard: function (cardName) {
    			console.log(cardName)
    			this.cards[cardName].showing = !this.cards[cardName].showing
    		},
    		store: function () {
    			this.$api.serviceDiscovery.storeService(this.service, (response) => {
    				if (response.status !== 201) {
    					this.informationStatus.msg = response.body.msg
    					this.informationStatus.isActive = true
    					this.informationStatus.className = 'alert-danger'
    				} else {
    					this.informationStatus.msg = 'Resource added successfully.'
					this.informationStatus.isActive = true
    					this.informationStatus.className = 'alert-success'
    				}
    			})
    		}
    	},
    	components: {
    		InformationPanel,
    		ServiceAPIConfiguration,
    		ServiceManagementConfig
    	}
    }
</script>
