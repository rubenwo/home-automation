<template>
    <div id="app" :style="{ backgroundImage: createBackgroundString }">
        <SideBar></SideBar>
        <app-toolbar class="mb2"/>

        <b-container class="content">
            <router-view/>
        </b-container>
    </div>
</template>
<script>
  import AppToolbar from "./components/AppToolbar.vue";
  import SideBar from "./components/SideBar.vue";

  export default {
    name: "app",
    components: {
      AppToolbar,
      SideBar
    },
    data() {
      return {
        angle: '50',
        color1: '#0D324D',
        color2: '#7F5A83'
      }
    },
    computed: {
      createBackgroundString() {
        return `linear-gradient(${this.angle}deg, ${this.color1}, ${this.color2})`;
      }
    }
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
        /*background-image: url("./assets/background.jpg");*/
    }
</style>
