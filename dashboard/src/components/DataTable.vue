<template>
<div>
<table class="table">
  <thead>
    <tr>
      <th scope="col" v-for="colName in headers">{{ colName }} </th>
      <th scope="col" v-if="Object.keys(actions).length > 0">Actions </th>
    </tr>
  </thead>
  <tbody>
    <tr v-for="row in data">
      <th scope="row" v-for="(item, index)  in row">
        <span>{{ simplifyString(item) }}</span>
        <br />
        <span class="text-info" v-if="item !== undefined && (item.length > maxLength || isJSON(item))" @click="showMore(item, index)" data-toggle="modal" data-target="#showMoreModal">
          Show more
        </span>
      </th>
      <th scope="row" v-if="Object.keys(actions).length > 0">
        <button v-for="action in actions" :class="'btn btn-' + (action.className||'primary')"
          @click="actionEmit(action.event, row)">{{ action.name }}</button>
      </th>
    </tr>
      
    </tr>
  </tbody>
</table>


<div class="modal fade" id="showMoreModal" tabindex="-1" role="dialog" aria-labelledby="exampleModalLabel" aria-hidden="true">
  <div class="modal-dialog" role="document">
    <div class="modal-content">
      <div class="modal-header">
        <h5 class="modal-title" id="exampleModalLabel">
        {{ showingMore.index }}
        </h5>
        <button type="button" class="close" data-dismiss="modal" aria-label="Close">
          <span aria-hidden="true">&times;</span>
        </button>
      </div>
      <div class="modal-body">
        <pre>{{ beautifyString(showingMore.item) }}</pre>
      </div>
      <div class="modal-footer">
        <button type="button" class="btn btn-secondary" data-dismiss="modal">Close</button>
        <button type="button" class="btn btn-primary">Save changes</button>
      </div>
    </div>
  </div>
</div>
</div>
</template>

<script>
export default {
  name: "datatable",
  props: ["headers", "data", "actions"],
  data() {
    return {
      maxLength: 25,
      showingMore: {
        index: null,
        item: null
      }
    };
  },
  mounted() {
    console.log(this.data);
  },
  methods: {
    isJSON: function(jsonStr) {
      try {
        JSON.parse(jsonStr);
      } catch (e) {
        return false;
      }
      return true;
    },

    showMore: function(item, index) {
      this.showingMore.index = index;
      this.showingMore.item = item;
    },
    simplifyString: function(jsonStr) {
      if (jsonStr === undefined || this.isJSON(jsonStr)) {
        return "";
      }
      var upToNCharacters = jsonStr.substring(
        0,
        Math.min(jsonStr.length, this.maxLength)
      );

      if (upToNCharacters == jsonStr) {
        return jsonStr;
      }
      return upToNCharacters + "...";
    },

    beautifyString: function(str) {
      if (this.isJSON(str)) {
        return JSON.parse(str);
      }
      return str;
    },

    actionEmit: function(event, row) {
      this.$emit(event, row);
    }
  }
};
</script>