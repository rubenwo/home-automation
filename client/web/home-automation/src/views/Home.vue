<template>
    <div v-if="devices.length > 0">
        <br>
        <input class="form-control" type="text" placeholder="Search" aria-label="Search" v-model="searchInput"/>
        <div v-bind:key="groupName" v-for="group, groupName in groups">
            <div>
                <h2 style="color:rgba(255, 255, 255, 0.45); text-align:center; border-bottom: 1px solid rgba(255, 255, 255, 0.45);">
                    {{groupName}}</h2>

                <div v-bind:key="device.id" v-for="device in group"
                     style="display:inline-block; margin-left:.95%; margin-right: .95%;">
                    <app-card
                            v-bind:height="10"
                            v-bind:name="device.name"
                            v-bind:category="device.category"
                            v-bind:company="device.product.company"
                            v-bind:device_type="device.product.type"
                            v-bind:img="getImgUrl(device)"
                            v-bind:id="device.id"
                    >
                    </app-card>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
  // @ is an alias to /src
  import AppCard from "@/components/AppCard.vue";
  import {mapActions, mapState} from "vuex";

  export default {
    name: "Home",
    data() {
      return {
        searchInput: ""
      };
    },
    computed: {
      ...mapState("devices", {
        devices: (state) => state.devices,
      }),
      groups() {
        const result = {};

        this.devices.forEach(device => {
          if (this.searchInput !== "") {
            if (!device.product.company.toLowerCase().includes(this.searchInput.toLowerCase())) {
              return;
            }
          }
          if (result[device.product.company] === undefined) {
            result[device.product.company] = []
          }

          result[device.product.company].push(device)
        });
        return result
      }
    },
    methods: {
      ...mapActions("devices", ["fetchDevices"]),
      getImgUrl(device) {
        if (device.product.company === 'IKEA') {
          switch (device.category) {
            case "0":
              return require("../assets/ikea_remote_control.png");
            case "2":
              return require("../assets/ikea_light_bulb.png");
          }
        }

        switch (device.category) {
          case "plug":
            return require("../assets/smart_plug_icon.png");
          case "light":
            return require("../assets/light_icon.png");
          case "led-strip":
            return require("../assets/led_strip_icon.png");
        }
      },

    },
    async mounted() {
      await this.fetchDevices();
    },
    components: {
      AppCard,
    },
  };
</script>
<style>
    .row {
        margin-top: 10px;
    }
</style>
