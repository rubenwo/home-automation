<template>
  <div style="text-align: center">
    <b-button pill @click="onButtonClick" center>
      <verte v-model="color" picker="wheel" model="rgb" />
    </b-button>

    <p></p>
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
    default: ""
  },
  components: { Verte },
  methods: {
    onButtonClick() {
      console.log(this.color);
      let rgb = this.color.replace(/[^\d,]/g, "").split(",");

      console.log(rgb);
      let command = {
        mode: "SINGLE_COLOR_RGB",
        red: parseInt(rgb[0]),
        green: parseInt(rgb[1]),
        blue: parseInt(rgb[2])
      };
      LedStripService.commandLedStripDevice(this.id, command);
    }
  },
  data() {
    return {
      color: "",
      device: {}
    };
  },
  async mounted() {
    this.id = this.$route.params.id;
    console.log(this.id);
    const res = await LedStripService.fetchAllLedStripDevices();
    console.log(res);
    this.device = res;
  }
};
</script>

<style></style>
