<template>
    <div style="text-align: center">
        <b-button pill @click="onButtonClick" center>
            <verte v-model="color" picker="wheel" model="rgb"/>
        </b-button>
    </div>
</template>

<script>
  import Verte from "verte";
  import "verte/dist/verte.css";
  import LedStripService from "../services/led_strip.service";

  export default {
    name: "LedStripDevice",
    id: {
      type: String,
      default: "",
    },
    components: {Verte},
    methods: {
      onButtonClick() {
        let rgb = this.color.replace(/[^\d,]/g, "").split(",");

        let command = {
          mode: "SINGLE_COLOR_RGB",
          red: parseInt(rgb[0]),
          green: parseInt(rgb[1]),
          blue: parseInt(rgb[2]),
        };
        LedStripService.commandLedStripDevice(this.id, command);
      },

    },
    data() {
      return {
        color: "",
        device: {},
      };
    },
    async mounted() {
      this.id = this.$route.params.id;
      const res = await LedStripService.fetchAllLedStripDevices();
      const devices = res.devices;

      for (let i = 0; i < devices.length; i++) {
        if (devices[i].id === this.id) {
          this.device = devices[i];
          break;
        }
      }
    },
  };
</script>

<style src="vue-color-gradient-picker/dist/index.css" lang="css"/>
