<template>
    <div id="app">
        <app-toolbar class="mb2"/>
        <b-container class="content">
            <router-view/>
        </b-container>
    </div>
</template>
<script>
  import AppToolbar from './components/AppToolbar.vue'

  export default {
    name: 'app',
    components: {
      AppToolbar
    }
  }
</script>
<style>
    #app {
        font-family: 'Avenir', Helvetica, Arial, sans-serif;
        -webkit-font-smoothing: antialiased;
        -moz-osx-font-smoothing: grayscale;

        color: #2c3e50;
    }
</style>
