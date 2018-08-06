<template>
<div>
  <div class="modal" tabindex="-1" role="dialog" :id="id">
    <div class="modal-dialog" role="document">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title text-info">{{ title }}</h5>
          <button type="button" class="close" data-dismiss="modal"  @click="closeModal" aria-label="Close">
            <span aria-hidden="true">&times;</span>
          </button>
        </div>
        <div class="modal-body">
          <p>{{ msg }}</p>
        </div>
        <div class="modal-footer">
          <button type="button" class="btn btn-success" data-dismiss="modal" @click="submitAnswer(true)">Yes</button>
          <button type="button" class="btn btn-danger" data-dismiss="modal" @click="submitAnswer(false)">No</button>
        </div>
      </div>
    </div>
  </div>

  <button :id="'openConfirmationModal'+id" type="button" data-toggle="modal" data-backdrop="static" :data-target="'#'+id"  style="display:none;">Launch modal</button>
</div>
</template>
<script>
export default {
  name: "home",
  props: ["id", "title", "msg", "showing"],
  watch: {
    showing: function() {
      if (this.showing === true) {
        this.openModal();
      }
    }
  },
  methods: {
    submitAnswer: function(answer) {
      this.$emit("answerReceived", answer);
      this.closeModal();
    },
    closeModal: function() {
      this.$emit("modalClosed");
    },
    openModal: function() {
      document.getElementById("openConfirmationModal" + this.id).click();
    }
  }
};
</script>
