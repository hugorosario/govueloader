<template id="hello-world-form">
    <form>
      <input v-model="name" placeholder="Type your name" />
      <button @click.prevent="submitPress">Set my name</button>
    </form>
</template>
<script>
Vue.component("hello-world-form", {
  template: "#hello-world-form",
  data: function () {
    return {
      name: null,
    };
  },
  methods: {
    submitPress() {
      this.$root.$emit('setName', this.name);
    },
  },
});
</script>
<style>
.hello-world {
  color: goldenrod;
}
</style>
