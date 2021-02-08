<template>
  <div v-if="this.state === 'loaded'">
    <b-card
      v-bind:sub-title="name"
      style="max-width: 540px; min-width: 200px; min-height: 425px; max-height: 500px; background-color: rgba(255, 255, 255, 0.7)"
      class="mb-2"
    >
      <b-card-img
        v-bind:src="img"
        alt="Image"
        height="130"
        width="130"
        class="mb-4"
      />
      <p>{{ category }}</p>
      <p>{{ company }} : {{ device_type }}</p>
      <div v-if="company === 'tp-link'">
        <b-button variant="success" @click="turnOnDevice">On</b-button>
        <b-button @click="turnOffDevice()">Off</b-button>
        <input
          v-if="device_type == 'L510E'"
          type="range"
          min="1"
          max="100"
          v-model="brightness"
          @change="brightnessChanged()"
        />
        <p>
          Status: {{ this.dev.device_info.device_on === "True" ? "On" : "Off" }}
        </p>
      </div>
      <div v-if="device_type === 'RGB_LED_STRIP'" style="text-align: center;">
        <b-button pill @click="onRGBClick()" center>
          <verte v-model="color" picker="wheel" model="rgb" />
        </b-button>
      </div>
      <div slot="footer">
        <b-button style="background-color: #4287f5;" v-bind:to="navigate()"
          >Information
        </b-button>
        <b-button variant="danger" @click="deleteDevice()">X</b-button>
      </div>
    </b-card>
  </div>
  <div v-else>
    <h3>Loading results...</h3>
    <Loading :active.sync="this.state === 'loading'" :is-full-page="true" />
  </div>
</template>

<script>
import Loading from "vue-loading-overlay";
import Verte from "verte";
import "verte/dist/verte.css";
import LedStripService from "../services/led_strip.service";
import TapoService from "../services/tapo.service";

export default {
  name: "app-card",
  components: { Loading, Verte },
  props: {
    name: {
      type: String,
      default: ""
    },
    img: {
      type: String,
      default: ""
    },
    category: {
      type: String,
      default: ""
    },
    company: {
      type: String,
      default: ""
    },
    device_type: {
      type: String,
      default: ""
    },
    id: {
      type: String,
      default: ""
    }
  },
  data() {
    return {
      brightness: 100,
      state: "loading",
      color: "",
      dev: {}
    };
  },
  methods: {
    navigate() {
      return "device/" + this.company + "/" + this.id;
    },
    async turnOnDevice() {
      if (this.company === "tp-link") {
        if (this.device_type === "L510E") {
          await TapoService.setDeviceBrightness(this.id, this.brightness);
        } else await TapoService.turnOnDevice(this.id);
        const deviceResult = await TapoService.fetchTapoDevice(this.id);
        this.dev = deviceResult.device;
      }
    },
    onRGBClick() {
      let rgb = this.color.replace(/[^\d,]/g, "").split(",");

      let command = {
        mode: "SINGLE_COLOR_RGB",
        red: parseInt(rgb[0]),
        green: parseInt(rgb[1]),
        blue: parseInt(rgb[2])
      };
      LedStripService.commandLedStripDevice(this.id, command);
    },
    async turnOffDevice() {
      if (this.company === "tp-link") {
        await TapoService.turnOffDevice(this.id);
        const deviceResult = await TapoService.fetchTapoDevice(this.id);
        this.dev = deviceResult.device;
      }
    },
    async deleteDevice() {
      if (this.company === "tp-link") {
        await TapoService.deleteTapoDevice(this.id);
      }
    },
    async brightnessChanged() {
      await TapoService.setDeviceBrightness(this.id, this.brightness);
    }
  },
  async mounted() {
    if (this.company === "tp-link") {
      await TapoService.wakeTapoDevice(this.id);
      const deviceResult = await TapoService.fetchTapoDevice(this.id);
      this.dev = deviceResult.device;
      this.state = "loaded";
    } else {
      this.state = "loaded";
    }
  }
};
</script>
