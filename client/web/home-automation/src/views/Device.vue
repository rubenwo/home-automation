<template>
    <div class="row">
        <div v-if="this.company === 'tp-link'">
            <tapo-device/>
        </div>
        <div v-else-if="this.company ==='hue'">
            <hue-device/>
        </div>
    </div>
</template>

<script>
  import TapoDevice from "../components/TapoDevice";
  import HueDevice from "../components/HueDevice";

  export default {
    name: "Device",
    components: {
      HueDevice,
      TapoDevice,
    },
    created() {
      this.id = this.$route.params.id;
      this.company = this.$route.params.company;
    },
  }
</script>

<style>
    .row {
        margin-top: 10px;
    }
</style>
