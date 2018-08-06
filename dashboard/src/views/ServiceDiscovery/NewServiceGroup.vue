<template>
    <div class="row">
        <div class="col-sm">
            <InformationPanel v-if="informationStatus.isActive" :msg="informationStatus.msg" :className="informationStatus.className"></InformationPanel>
            <h2>Add new Service Layer</h2>
            
            <form v-on:keyup.13="store">
                <div class="row">
                    <div class="col-sm-6 offset-sm-1">
                        <div class="form-group">
                            <label for="groupName" class="text-info">Name</label>
                            <input type="text" v-model="group.Name" class="form-control" id="groupName" aria-describedby="nameHelp" placeholder="Enter name">
                            <small id="nameHelp" class="form-text text-muted">Give the layer a name.</small>
                        </div>
                        
                        <div class="form-group">                      
                            <i class="fas " :class="group.IsReachable ? 'fa-eye text-success' : 'fa-eye-slash text-danger'" @click="toggleReachable" />
                            <label class="form-check-label" for="groupReachable">&nbsp;&nbsp;Reachable</label>
                            <small id="groupReachableHelp" class="form-text text-info">Is group reachable from external sources?</small> 
                        </div>
                    </div>
                </div>  
            </form>
    
            <div class="row">
                <div class="col-sm offset-sm-1">
                    <button type="submit" class="btn btn-primary" v-on:click="store" >Save</button>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import InformationPanel from '@/components/InformationPanel'
import { mapGetters } from 'vuex'

export default {
	name: 'new-group',
	data () {
		return {
			group: {
				Name: '',
				IsReachable: false,
				Services: []
			},
			informationStatus: {
				isActive: false,
				className: 'alert-success',
				msg: ''
			}
		}
	},
	computed: {
		...mapGetters({
			isAdmin: 'isAdmin'
		})
	},
	methods: {
		store: function () {
			this.$api.serviceDiscovery.storeServiceGroup(this.group, (response) => {
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
		},
		toggleReachable: function () {
			if (!this.isAdmin) return
			this.group.IsReachable = !this.group.IsReachable
		}
	},
	components: {
		InformationPanel
	}
}
</script>

<style>

</style>
