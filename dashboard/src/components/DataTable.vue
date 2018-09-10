<template>
<div>
  <div class="row" v-if="searchable">
    <div class="col-sm-5 form-inline ">
        <input class="form-control" v-model="searchText"/>
        <button class="btn btn-sm btn-info" @click="search">
        <i class="fas fa-search"></i>
        </button>
    </div>
    
    <div class="col-sm-3 offset-sm-4 form-inline">
        <button class="btn btn-sm btn-info" @click="changePage(currentPage - 1 < 1 ? currentPage = 1 : currentPage -= 1)">
          <i class="fas fa-arrow-left"></i>
        </button>
        <input class="form-control sm-1 mr-sm-1" v-model="currentPage" />
        <button class="btn btn-sm btn-info" @click="changePage(currentPage += 1)">
          <i class="fas fa-arrow-right"></i>
        </button>
    </div>
  </div>
        
<table class="table">
  <thead>
    <tr>
      <th scope="col" v-for="colName in headers" v-bind:key="colName">{{ colName }} </th>
      <th scope="col" v-if="Object.keys(actions).length > 0">Actions </th>
    </tr>
  </thead>
  <tbody>
    <tr v-for="(row, index) in data" v-bind:key="index">
      <th scope="row" v-for="(item, index)  in row" v-bind:key="index">
        <span>{{ simplifyString(item) }}</span>
        <br />
        <span class="text-info toggable-card" v-if="item !== undefined && (item.length > maxLength || isJSON(item))" @click="showMore(item, index)" data-toggle="modal" data-target="#showMoreModal">
          Show more
        </span>
      </th>
      <th scope="row" v-if="Object.keys(actions).length > 0">
        <button v-for="action in actions" :class="'btn btn-' + (action.className||'primary')" v-bind:key="action.name"
          @click="actionEmit(action.event, row)">{{ action.name }}</button>
      </th>
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
  props: ["headers", "data", "actions", "searchable"],
  data() {
    return {
      maxLength: 25,
      showingMore: {
        index: null,
        item: null
      },
      searchText: "",
      currentPage: 1
    };
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

      if (upToNCharacters === jsonStr) {
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
    },

    search: function() {
      this.$emit("search", this.searchText);
    },

    changePage: function(value) {
      this.currentPage = value;
      this.$emit("changePage", this.currentPage);
    }
  }
};
</script>
