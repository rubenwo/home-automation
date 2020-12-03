<template>
    <div class="row" v-if="!this.devices.length <= 0">
        <b-col
                cols="4"
                sm="3"
                md="3"
                lg="2"
                xl="2"
                v-bind:key="device.id"
                v-for="device in this.devices"
        >
            <app-card
                    v-bind:height="10"
                    v-bind:title="device.name"
                    v-bind:id="device.id"
                    style="max-width: 540px;"
            >
                >
            </app-card>
        </b-col>
    </div>
</template>

<script>
  // @ is an alias to /src
  import AppCard from "@/components/AppCard.vue";
  import {mapActions, mapState} from "vuex";

  export default {
    name: "Home",
    data() {
      return {};
    },
    computed: {
      ...mapState("devices", {
        devices: state => state.devices
      })
    },
    methods: {
      ...mapActions("devices", ["fetchDevices", "addNewDevice"])
    },
    async mounted() {
      await this.fetchDevices();
      let data = {
        "ip_address": "192.168.2.X",
        "email": "",
        "password": "",
        "device_type": "P100"
      }
      await this.addNewDevice("tapo", data);
    },
    components: {
      AppCard
    }
  };
</script>
<style>
    .row {
        margin-top: 10px;
    }
</style>
