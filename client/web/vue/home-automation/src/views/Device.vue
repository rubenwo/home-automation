<template>
  <div class="row" style="background-color: rgba(255, 255, 255, 0.7)">
    <div v-if="this.company === 'tp-link'">
      <tapo-device />
    </div>
    <div v-else-if="this.company === 'hue'">
      <hue-device />
    </div>
    <div v-else-if="this.company === 'esp32'">
      <led-strip-device />
    </div>
  </div>
</template>

<script>
import TapoDevice from "../components/TapoDevice";
import HueDevice from "../components/HueDevice";
import LedStripDevice from "../components/LedStripDevice";

export default {
  name: "Device",
  components: {
    LedStripDevice,
    HueDevice,
    TapoDevice
  },
  created() {
    this.id = this.$route.params.id;
    this.company = this.$route.params.company;
  }
};
</script>

<style>
.row {
  margin-top: 10px;
}
</style>
