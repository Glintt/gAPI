<template>
<div>
  <div class="modal fade" tabindex="-1" role="dialog" :id="id">
    <div class="modal-dialog" role="document">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title text-danger">{{ title }}</h5>
          <button type="button" class="close" data-dismiss="modal"  @click="closeModal" aria-label="Close">
            <span aria-hidden="true">&times;</span>
          </button>
        </div>
        <div class="modal-body">
          <div class="row">
            <div class="col-sm text-center" v-for="(type, index) in managementTypes" v-bind:key="index">
              <button @click="manageService(service, type.action)" 
                  :class="'btn btn-'+ type.background"
                  v-show="! $api.serviceDiscovery.CustomManagementActions.includes(type.action)"
                  data-toggle="tooltip" data-placement="top" :title="type.description">
                  <i :class="type.icon"></i> {{ type.description }}
              </button>
              <br/><br/>
            </div>
          </div>
        </div>
        <div class="modal-footer">
          <button type="button" class="btn btn-danger" data-dismiss="modal" @click="closeModal">Close</button>
        </div>
      </div>
    </div>
  </div>

  <button :id="'openModal'+id" type="button" data-toggle="modal" data-backdrop="static" :data-target="'#'+id" style="display:none;">Launch modal</button>
  
  <ErrorMessage @modalClosed="statusModalClosed" :showing="statusMessage.showing && statusMessage.isError" :id="'requestError'" :error="statusMessage.msg" :title="'Error Occurred'"/>
  <SuccessModal @modalClosed="statusModalClosed" :showing="statusMessage.showing && !statusMessage.isError" :id="'requestSuccess'" :msg="statusMessage.msg" :title="'Success'"/>
  <ConfirmationModal @answerReceived="managementConfirmationReceived" @modalClosed="confirmationClosed" :showing="confirmation.showing" :id="'managementConfirm'" :msg="confirmation.msg" :title="confirmation.title"/>
        
</div>
</template>
<script>
import ConfirmationModal from "@/components/modals/ConfirmationModal";
import ErrorMessage from "@/components/modals/ErrorMessage";
import SuccessModal from "@/components/modals/SuccessModal";

export default {
  name: "service-management-modal",
  props: ["id", "title", "error", "showing", "service"],
  created() {
    this.$api.serviceDiscovery.manageServiceTypes(response => {
      this.managementTypes = response.body;
    });
  },
  data() {
    return {
      managementTypes: {},
      statusMessage: {
        msg: "",
        showing: false,
        isError: false
      },
      management: {
        action: "",
        service: this.service
      },
      confirmation: {
        showing: false,
        title: "",
        msg: ""
      }
    };
  },
  watch: {
    showing: function() {
      if (this.showing === true) this.openModal();
    }
  },
  methods: {
    closeModal: function() {
      this.$emit("modalClosed");
    },
    openModal: function() {
      document.getElementById("openModal" + this.id).click();
    },

    manageService: function(service, action) {
      this.confirmation.showing = true;
      this.confirmation.title = "Confirm - " + action;
      this.confirmation.msg =
        "Are you sure you want to " + action + " service " + service.Name + "?";
      this.management.service = service;
      this.management.action = action;
    },
    managementConfirmationReceived: function(answer) {
      if (answer === false) return;

      this.$api.serviceDiscovery.manageService(
        this.management.service.MatchingURI,
        this.management.action,
        response => {
          this.statusMessage.msg = response.body.msg;
          this.statusMessage.isError = false;
          if (response.status !== 200) {
            this.statusMessage.isError = true;
            if (response.body.service_response !== undefined) {
              this.statusMessage.msg = response.body.service_response;
            }
          }
          this.statusMessage.showing = true;
        }
      );
    },
    confirmationClosed: function() {
      this.confirmation.showing = false;
      this.confirmation.msg = "";
    },
    statusModalClosed: function() {
      this.statusMessage.showing = false;
      this.statusMessage.msg = "";
    }
  },
  components: {
    ConfirmationModal,
    ErrorMessage,
    SuccessModal
  }
};
</script>
