<template>
<div>    
    <div class="row" v-if="! documentationExpanded">
        <div class="col-sm">
            <EditService  v-on:serviceUpdated="updateDocumentationEndpoint"></EditService>
        </div>   
    </div>
    <br  v-if="! documentationExpanded" />
    
    <div class="row">
        <div class="col-sm-6">
            <h2 v-if="! documentationExpanded">Documentation</h2>            
        </div>
        <div class="offset-sm-5 col-sm-1 text-right" style="cursor: pointer;" @click="expandDocumentation">
            <i class="fas fa-expand-arrows-alt text-primary"></i>
        </div>
    </div>
    <hr  v-if="! documentationExpanded"/> 

    <iframe ref="documentation" id="documentation" v-if="ServiceDocumentationEndpoint !== null" frameborder="0" :style="documentationStyle" :height="documentationHeight" width="100%" :src="ServiceDocumentationEndpoint"></iframe>
</div>
   
</template>
<style>

#documentation {
    -moz-transition: height .5s;
    -ms-transition: height .5s;
    -o-transition: height .5s;
    -webkit-transition: height .5s;
    transition: height .5s;
    height: 0;
    overflow: hidden;
  }
</style>

<script>
import EditService from '@/views/Service/EditService'
import { mapActions } from 'vuex'

// var serviceDiscoveryAPI = require('@/api/service-discovery')

export default {
	name: 'view-service',
	data () {
		return {
			ServiceDocumentationEndpoint: null,
			documentationExpanded: false,
			documentationHeight: '200%',
			documentationStyle: 'min-height:600px;width:100%;top:0px;left:0px;right:0px;bottom:0px'
		}
	},
	methods: {
		...mapActions('fullscreen', [
			'openFullScreen',
			'closeFullScreen'
		]),
		updateDocumentationEndpoint: function (service) {
			this.ServiceDocumentationEndpoint = service.APIDocumentation
		},
		expandDocumentation: function () {
			this.documentationExpanded = !this.documentationExpanded
			if (this.documentationExpanded) {
				console.log()
				this.documentationStyle = 'min-height:' + (screen.height - 165) + 'px;width:100%;top:0px;left:0px;right:0px;bottom:0px'
				this.openFullScreen()
			} else {
				this.closeFullScreen()
			}
		}
	},
	components: {
		EditService
	}
}
</script>
