<template>
    <b-navbar
            :sticky="true"
            toggleable="lg"
            type="light"
            style="background-color: rgba(255, 255, 255, 0.25);
        backdrop-filter: blur(5px);"
    >

        <b-container>
            <b-button v-b-toggle.sidebar-1>Toggle Sidebar</b-button>
            <b-navbar-toggle target="nav_collapse"/>

            <b-collapse is-nav id="nav_collapse">
                <b-navbar-nav>
                    <b-nav-item href="#" to="/" exact>Home</b-nav-item>
                    <b-nav-item href="#" to="/routines" exact>Routines</b-nav-item>
                    <b-nav-item href="#" to="/sensors" exact>Sensors</b-nav-item>
                    <b-nav-item href="#" to="/recipes" exact>Recipes</b-nav-item>
                    <b-nav-item href="#" to="/inventory" exact>Inventory</b-nav-item>
                    <b-nav-item href="#" to="/cameras" exact>Cameras</b-nav-item>
                </b-navbar-nav>
                <b-navbar-nav class="ml-auto">
                    <b-button
                            pill
                            img="./assets/add.png"
                            v-b-tooltip.hover
                            variant="primary"
                            title="Click here to add a new"
                            @click="onClickAdd"
                    ><img src="../assets/add.png" width="25" height="25"/>
                    </b-button>
                    <b-nav-item-dropdown right v-if="isLoggedIn">
                        <template slot="button-content">{{ username }}</template>
                        <b-dropdown-item v-on:click="logout">Signout</b-dropdown-item>
                    </b-nav-item-dropdown>
                </b-navbar-nav>
            </b-collapse>
        </b-container>
        <add-device-modal ref="deviceModal"/>
        <add-routine-modal ref="routineModal"/>
        <add-sensor-modal ref="sensorModal"/>
        <add-recipe-modal ref="recipeModal"/>
    </b-navbar>
</template>

<script>
  import AddDeviceModal from "./AddDeviceModal";
  import AddRecipeModal from "./AddRecipeModal";
  import AddRoutineModal from "./AddRoutineModal";
  import AddSensorModal from "./AddSensorModal";
  import {mapActions, mapState, mapGetters} from "vuex";

  export default {
    name: "app-toolbar",
    data() {
      return {};
    },
    components: {
      AddRecipeModal,
      AddDeviceModal,
      AddRoutineModal,
      AddSensorModal,
    },
    computed: {
      ...mapState("auth", ["username"]),
      ...mapGetters("auth", ["isLoggedIn"]),
    },

    methods: {
      ...mapActions("auth", ["logout"]),
      getAddButtonHint() {
        return "test";
      },
      onClickAdd() {
        switch (this.$route.name) {
          case "Home":
            this.$refs.deviceModal.$emit("add_device");
            break;
          case "Recipes":
            this.$refs.recipeModal.$emit("add_recipe");
            break;
          case "Routines":
            this.$refs.routineModal.$emit("add_routine");
            break;
          case "Sensors":
            this.$refs.sensorModal.$emit("add_sensor");
            break;
        }
      },
    },
  };
</script>
