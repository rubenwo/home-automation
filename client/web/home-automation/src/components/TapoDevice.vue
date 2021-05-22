<template>
  <div>
    <b-img v-bind:src="getImgUrl()" fluid center />
    <p>{{ tapoDevice.device_name }}</p>

    <p>{{ tapoDevice }}</p>

    <div v-if="this.tapoDevice.device_type === 'L510E'">
      <p>LIGHT</p>
    </div>
    <div v-else-if="this.tapoDevice.device_type === 'P100'">
      <p>PLUG</p>
    </div>
  </div>
</template>

<script>
import { mapActions, mapState } from "vuex";

export default {
  name: "TapoDevice",
  id: {
    type: String,
    default: "",
  },
  computed: {
    ...mapState("tapo", {
      tapoDevice: (state) => state.tapoDevice,
    }),
  },
  methods: {
    ...mapActions("tapo", ["fetchTapoDevice"]),
    getImgUrl() {
      switch (this.tapoDevice.device_type) {
        case "P100":
          console.log("returning plug");
          return require("../assets/smart_plug_icon.png");
        case "L510E":
          console.log("returning light");
          return require("../assets/light_icon.jpg");
      }
    },
  },
  async mounted() {
    this.id = this.$route.params.id;
    await this.fetchTapoDevice(this.id);
    console.log(this.tapoDevice);
  },
};
</script>

<style scoped></style>
