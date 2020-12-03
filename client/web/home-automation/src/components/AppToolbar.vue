<template>
    <b-navbar :sticky="true" toggleable="lg" type="dark" style="background-color: #4287f5;">
        <b-container>
            <b-navbar-brand href="#" to="/">LabApp</b-navbar-brand>

            <b-navbar-toggle target="nav_collapse"/>

            <b-collapse is-nav id="nav_collapse">
                <b-navbar-nav>
                    <b-nav-item href="#" to="/" exact>Home</b-nav-item>
                    <b-nav-item href="#" to="/device" exact>Device</b-nav-item>
                    <!-- <b-nav-item href="#" to="/settings">Settings</b-nav-item> -->
                </b-navbar-nav>
                <b-navbar-nav class="ml-auto">

                    <b-button
                            pill
                            img="./assets/add.png"
                            v-b-tooltip.hover
                            variant="primary"
                            title="Click here to add a new device!"
                            @click="onClickAdd"
                            v-b-modal.help-modal
                    ><img src="../assets/add.png" width="25" height="25"/>
                    </b-button>
                </b-navbar-nav>
            </b-collapse>
        </b-container>
    </b-navbar>
</template>

<script>
  import {mapActions, mapState} from "vuex";

  export default {
    name: "app-toolbar",
    data() {
      return {};
    },

    computed: {
      ...mapState("devices", {
        devices: state => state.devices
      })
    },

    methods: {
      ...mapActions("devices", ["addNewDevice"]),
      onClickAdd: async () => {
        await this.addNewDevice("tapo", {
          "ip_address": "192.168.2.X",
          "email": "",
          "password": "",
          "device_type": "P100"
        });
        console.log("clicked add")
      }
    }
  };
</script>
