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
      style="margin-right: 25px; margin-left: 25px"
    >
      <app-card
        v-bind:height="10"
        v-bind:name="device.name"
        v-bind:category="device.category"
        v-bind:company="device.product.company"
        v-bind:device_type="device.product.type"
        v-bind:img="getImgUrl(device)"
        v-bind:id="device.id"
      >
        >
      </app-card>
    </b-col>
  </div>
</template>

<script>
// @ is an alias to /src
import AppCard from "@/components/AppCard.vue";
import { mapActions, mapState, mapGetters } from "vuex";

export default {
  name: "Home",
  data() {
    return {};
  },
  computed: {
    ...mapState("devices", {
      devices: state => state.devices
    }),
    ...mapState("auth", ["error"])
  },
  methods: {
    ...mapActions("devices", ["fetchDevices"]),
    ...mapActions("auth", ["login", "logout"]),
    ...mapGetters("auth", ["isLoggedIn"]),
    getImgUrl(device) {
      switch (device.category) {
        case "plug":
          console.log("returning plug");
          return require("../assets/smart_plug_icon.png");
        case "light":
          console.log("returning light");
          return require("../assets/light_icon.jpg");
        case "led-strip":
          console.log("returning led-strip");
          return require("../assets/led_strip_icon.png");
      }
    }
  },
  async mounted() {
    console.log(this.isLoggedIn());
    if (!this.isLoggedIn()) {
      let data = {
        username: "admin",
        password: process.env.VUE_APP_ADMIN_PWD
      };
      await this.login(data);
    }
    await this.fetchDevices();
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
