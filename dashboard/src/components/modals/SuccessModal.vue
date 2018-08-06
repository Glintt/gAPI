<template>
<div>
  <div class="modal fade" tabindex="-1" role="dialog" :id="id">
    <div class="modal-dialog" role="document">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title text-success">{{ title }}</h5>
          <button type="button" class="close" data-dismiss="modal"  @click="closeModal" aria-label="Close">
            <span aria-hidden="true">&times;</span>
          </button>
        </div>
        <div class="modal-body">
          <div  v-html="msg"></div>
        </div>
        <div class="modal-footer">
          <button type="button" class="btn btn-success" data-dismiss="modal" @click="closeModal">Close</button>
        </div>
      </div>
    </div>
  </div>

  <button :id="'openSuccessModal'+id" type="button" data-toggle="modal" data-backdrop="static" :data-target="'#'+id" style="display:none;">Launch modal</button>
</div>  
</template>

<script>
export default {
  name: "SuccessModal",
  props: ["id", "title", "msg", "showing"],
  watch: {
    showing: function() {
      if (this.showing === true) {
        this.openModal();
      }
    }
  },
  methods: {
    closeModal: function() {
      this.$emit("modalClosed");
    },
    openModal: function() {
      document.getElementById("openSuccessModal" + this.id).click();
    }
  }
};
</script>
