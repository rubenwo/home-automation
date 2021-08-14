<template>
    <div style="text-align: center">
        <b-button pill @click="onButtonClick" center>
            <verte v-model="color" picker="wheel" model="rgb"/>
        </b-button>

        <p></p>


        <ColorPicker
                :color="colour"
                :onStartChange="colour => onColorChange(colour, 'start')"
                :onChange="colour => onColorChange(colour, 'change')"
                :onEndChange="colour => onColorChange(colour, 'end')"
        />
    </div>
</template>

<script>
  import Verte from "verte";
  import "verte/dist/verte.css";
  import LedStripService from "../services/led_strip.service";
  import {ColorPicker} from 'vue-color-gradient-picker';

  export default {
    name: "LedStripDevice",
    id: {
      type: String,
      default: "",
    },
    components: {Verte, ColorPicker},
    methods: {
      onButtonClick() {
        console.log(this.color);
        let rgb = this.color.replace(/[^\d,]/g, "").split(",");

        console.log(rgb);
        let command = {
          mode: "SINGLE_COLOR_RGB",
          red: parseInt(rgb[0]),
          green: parseInt(rgb[1]),
          blue: parseInt(rgb[2]),
        };
        LedStripService.commandLedStripDevice(this.id, command);
      },

      // eslint-disable-next-line no-unused-vars
      onColorChange(attrs, name) {
        console.log("Change");
        this.colour = { ...attrs };
      }
    },
    data() {
      return {
        color: "",
        device: {},
        colour: {
          red: 255,
          green: 0,
          blue: 0,
          alpha: 1
        }
      };
    },
    async mounted() {
      this.id = this.$route.params.id;
      console.log(this.id);
      const res = await LedStripService.fetchAllLedStripDevices();
      const devices = res.devices;

      for (let i = 0; i < devices.length; i++) {
        if (devices[i].id === this.id) {
          this.device = devices[i];
          break;
        }
      }
      console.log(this.device);

    },
  };
</script>

<style src="vue-color-gradient-picker/dist/index.css" lang="css" />
