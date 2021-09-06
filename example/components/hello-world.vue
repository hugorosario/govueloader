<template id="hello-world">
    <h1 class="hello-world">Hello, World 2!</h1>
</template>
<script>
    Vue.component("hello-world", {template: "#hello-world"});
</script>
<style scoped>
    .hello-world {
        color: goldenrod;
    }
</style>