<template id="hello-world">
  <div>
    <h1 class="hello-world">Hello, {{ name }}!</h1>
    <hello-world-form></hello-world-form>
    <h2><a href="/vuetify">Check this nice Vuetify app</a></h2>
  </div>
</template>
<script>
Vue.component("hello-world", {
  template: "#hello-world",
  data: function () {
    return {
      name: "world",
    };
  },
  mounted: function () { 
      this.$root.$on('setName', (newname) => { 
        this.name = newname;
      })
    }  
});
</script>
<style>
.hello-world {
  color: goldenrod;
}
</style>
