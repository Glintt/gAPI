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
          <div  v-html="error"></div>
        </div>
        <div class="modal-footer">
          <button type="button" class="btn btn-danger" data-dismiss="modal" @click="closeModal">Close</button>
        </div>
      </div>
    </div>
  </div>

  <button :id="'openModal'+id" type="button" data-toggle="modal" data-backdrop="static" :data-target="'#'+id" style="display:none;">Launch modal</button>
</div>
</template>
<script>

    export default {
        name: "error-message-modal",
        props: ["id","title", "error", "showing"],
        watch:{
          showing: function(){
            if (this.showing == true)
              this.openModal();
          }
        },
        methods:{
          closeModal:function(){
            this.$emit("modalClosed")
          },
          openModal: function(){
            document.getElementById("openModal"+this.id).click();
          }
        }
    }
</script>