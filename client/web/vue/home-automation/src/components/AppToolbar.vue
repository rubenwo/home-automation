<template>
    <b-navbar :sticky="true" toggleable="lg" type="dark" style="background-color: rgba(70, 70, 70, 0.7)">
        <b-container>
            <b-navbar-brand href="#" to="/">Home Automation</b-navbar-brand>

            <b-navbar-toggle target="nav_collapse"/>

            <b-collapse is-nav id="nav_collapse">
                <b-navbar-nav>
                    <b-nav-item href="#" to="/" exact>Home</b-nav-item>
                    <b-nav-item href="#" to="/recipes" exact>Recipes</b-nav-item>
                </b-navbar-nav>
                <b-navbar-nav class="ml-auto">

                    <b-button
                            pill
                            img="./assets/add.png"
                            v-b-tooltip.hover
                            variant="primary"
                            title="Click here to add a new device!"
                            @click="onClickAdd"
                    ><img src="../assets/add.png" width="25" height="25"/>
                    </b-button>
                </b-navbar-nav>
            </b-collapse>
        </b-container>
        <add-device-modal ref="deviceModal"/>
        <add-recipe-modal ref="recipeModal"/>
    </b-navbar>
</template>

<script>
  import {mapActions, mapState} from "vuex";
  import AddDeviceModal from "./AddDeviceModal";
  import AddRecipeModal from "./AddRecipeModal";

  export default {
    name: "app-toolbar",
    data() {
      return {};
    },
    components: {
      AddRecipeModal,
      AddDeviceModal
    },
    computed: {
      ...mapState("devices", {
        devices: state => state.devices
      })
    },

    methods: {
      ...mapActions("devices", ["addNewDevice"]),
      onClickAdd() {
        console.log(this.$route)
        switch (this.$route.name) {
          case "Home":
            this.$refs.deviceModal.$emit('add_device');
            break;
          case "Recipes":
            this.$refs.recipeModal.$emit('add_recipe');
            break;
        }
      },
    }
  };
</script>
