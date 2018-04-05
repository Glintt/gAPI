<template>
<table class="table">
  <thead>
    <tr>
      <th scope="col" v-for="colName in headers">{{ colName }} </th>
    </tr>
  </thead>
  <tbody>
    <tr v-for="row in data">
      <th scope="row" v-for="(item, index)  in row">{{item}}</th>
    </tr>
  </tbody>
</table>
</template>

<script>
    export default {
        name: "datatable",
        props : ["headers", "data"]
    }
</script>