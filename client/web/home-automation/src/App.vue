<template>
  <div id="app">
    <app-toolbar class="mb2" />

    <b-container class="content">
      <router-view />
    </b-container>
  </div>
</template>
<script>
  import AppToolbar from "./components/AppToolbar.vue";

  export default {
    name: "app",
    components: {
      AppToolbar,
    },
  };
</script>
<style>
  #app {
    font-family: "Avenir", Helvetica, Arial, sans-serif;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;

    padding: 0;
    margin: 0;
    width: 100%;
    min-height: 100vh;
    background-image: url("./assets/background.jpg");
  }
</style>
