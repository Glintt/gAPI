<template>
    <div class="alert" v-bind:class="className" role="alert">
        {{ msg }}
    </div>
</template>

<script>
export default {
  name: "information-panel",
  props: ["msg", "className"]
};
</script>
