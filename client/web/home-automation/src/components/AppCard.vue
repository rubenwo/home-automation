<template>
    <div v-if="this.state === 'loaded'">
        <div v-if="this.dev !== null">
            <b-card
                    v-bind:title="name"
                    style="
        max-width: 200px;
        min-width: 200px;
        min-height: 550px;
        max-height: 550px;
        background-color: rgba(255, 255, 255, 0.25);
        backdrop-filter: blur(5px);
      "
                    class="mb-2"
            >
                <div v-if="device_type === 'RGB_LED_STRIP'" align="center">
                    <color-picker v-model="color" :width="130" :height="130" @color-change="onRGBClick"/>
                </div>
                <div v-else>
                    <b-card-img
                            v-bind:src="img"
                            alt="Image"
                            height="130"
                            width="130"
                            class="mb-4"
                    />
                </div>
                <p>{{ category }}</p>
                <p>{{ company }} : {{ device_type }}</p>
                <div v-if="company === 'tp-link'">
                    <toggle-button v-model="device_on" @change="onChangeEventHandler"/>
                    <input
                            v-if="device_type === 'L510E'"
                            type="range"
                            min="1"
                            max="100"
                            v-model="brightness"
                            @change="brightnessChanged()"
                    />
                    <p>
                        Status: {{ this.dev.device_info.device_on ? "On" : "Off" }}
                    </p>
                </div>

                <div v-if="device_type === 'RGB_LED_STRIP'">
                    <b-button pill @click="onRGBClickVerte" center>
                        <verte v-model="colorVerte" picker="wheel" model="rgb"/>
                    </b-button>

                    <b-button pill @click="onColorCycleClick" center>Color Cycle</b-button>
                    <b-button pill @click="onChristmasClick" center>Christmas</b-button>
                    <b-button pill @click="onBreathingClick" center>Pulse</b-button>

                </div>

                <div v-if="company === 'IKEA'">
                    <toggle-button v-if="category === '2'" v-model="device_on" @change="onChangeEventHandler"/>
                    <input
                            v-if="category === '2'"
                            type="range"
                            min="0"
                            max="255"
                            v-model="brightness"
                            @change="brightnessChanged()"
                    />

                </div>

                <div slot="footer">
                    <b-button style="background-color: #4287f5" v-bind:to="navigate()"
                    >Information
                    </b-button>
                    <b-button variant="danger" @click="deleteDevice()">X</b-button>
                </div>
            </b-card>
        </div>
        <div v-else>
            <b-card
                    v-bind:title="name"
                    style="
        max-width: 200px;
        min-width: 200px;
        min-height: 500px;
        max-height: 500px;
        background-color: rgba(255, 255, 255, 0.25);
        backdrop-filter: blur(5px);
      "
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
                <div style="color:#ff0033">Could not load this device. Check the backend</div>

            </b-card>
        </div>
    </div>
    <div v-else class="text-center">
        <h3>Loading results...</h3>
        <b-spinner style="width: 3rem; height: 3rem;" label="Large Spinner"/>
    </div>
</template>

<script>
  import "vue-loading-overlay/dist/vue-loading.css";
  import ColorPicker from 'vue-color-picker-wheel';
  import Verte from "verte";
  import "verte/dist/verte.css";
  import LedStripService from "../services/led_strip.service";
  import TapoService from "../services/tapo.service";
  import {ToggleButton} from 'vue-js-toggle-button'
  import TradfriService from "../services/ikea.service";

  export default {
    name: "app-card",
    components: {ColorPicker, ToggleButton, Verte},
    props: {
      name: {
        type: String,
        default: "",
      },
      img: {
        type: String,
        default: "",
      },
      category: {
        type: String,
        default: "",
      },
      company: {
        type: String,
        default: "",
      },
      device_type: {
        type: String,
        default: "",
      },
      id: {
        type: String,
        default: "",
      },
    },
    data() {
      return {
        device_on: false,
        brightness: 100,
        state: "loading",
        colorPickerInitiated: false,
        color: "#E17D0F",
        colorVerte: "",
        dev: {},
      };
    },
    methods: {
      onChangeEventHandler(e) {
        if (e.value === true) this.turnOnDevice();
        else if (e.value === false) this.turnOffDevice();
      },
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
          this.device_on = this.dev.device_info.device_on;
        } else if (this.company === "IKEA") {
          await TradfriService.commandTradfriDevice(this.id, {
            device_type: "light",
            dimmable_light_command: {
              power: 1
            }
          });
        }
      },
      onRGBClick() {
        if (!this.colorPickerInitiated) {
          this.colorPickerInitiated = true;
          return;
        }
        let rgb = /^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(this.color);
        rgb = {
          r: parseInt(rgb[1], 16),
          g: parseInt(rgb[2], 16),
          b: parseInt(rgb[3], 16)
        };

        let command = {
          mode: "SINGLE_COLOR_RGB",
          red: rgb.r,
          green: rgb.g,
          blue: rgb.b,
        };

        LedStripService.commandLedStripDeviceSolid(this.id, command);
      },
      onColorCycleClick() {
        LedStripService.commandLedStripDeviceColorCycle(this.id);
      },
      onChristmasClick() {
        LedStripService.commandLedStripDeviceChristmas(this.id);
      },
      onBreathingClick() {
        let rgb = /^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(this.color);
        rgb = {
          r: parseInt(rgb[1], 16),
          g: parseInt(rgb[2], 16),
          b: parseInt(rgb[3], 16)
        };

        let command = {
          red: rgb.r,
          green: rgb.g,
          blue: rgb.b,
        };
        LedStripService.commandLedStripDeviceBreathing(this.id, command);
      },
      onRGBClickVerte() {
        let rgbToHex = (r, g, b) => {
          return "#" + ((1 << 24) + (r << 16) + (g << 8) + b).toString(16).slice(1);
        };
        let rgb = this.colorVerte.replace(/[^\d,]/g, "").split(",");
        this.color = rgbToHex(parseInt(rgb[0]), parseInt(rgb[1]), parseInt(rgb[2]));
      },
      async turnOffDevice() {
        if (this.company === "tp-link") {
          await TapoService.turnOffDevice(this.id);
          const deviceResult = await TapoService.fetchTapoDevice(this.id);
          this.dev = deviceResult.device;
          this.device_on = this.dev.device_info.device_on;
        } else if (this.company === "IKEA") {
          await TradfriService.commandTradfriDevice(this.id, {
            device_type: "light",
            dimmable_light_command: {
              power: 0
            }
          });
        }
      },
      async deleteDevice() {
        if (this.company === "tp-link") {
          await TapoService.deleteTapoDevice(this.id);
        }
      },
      async brightnessChanged() {
        if (this.company === "tp-link") {
          await TapoService.setDeviceBrightness(this.id, this.brightness);
        } else if (this.company === "IKEA") {
          await TradfriService.commandTradfriDevice(this.id, {
            device_type: "light",
            dimmable_light_command: {
              brightness: parseInt(this.brightness)
            }
          });
        }
      },
    },
    async mounted() {
      if (this.company === "tp-link") {
        // await TapoService.wakeTapoDevice(this.id);
        const deviceResult = await TapoService.fetchTapoDevice(this.id);
        if (deviceResult === "couldn't retrieve tapo device info\n") {
          this.dev = null;
        } else {
          this.dev = deviceResult.device;
          this.device_on = this.dev.device_info.device_on;
        }
        this.state = "loaded";
      } else if (this.device_type === 'RGB_LED_STRIP') {
        // TODO: get the current RGB color from the
        this.state = "loaded";
      } else if (this.company === "IKEA") {
        let device = await TradfriService.fetchTradfriDevice(this.id);
        if (device.device_type === 'light') {
          this.brightness = device.dimmable_device_data.brightness;
          this.device_on = device.dimmable_device_data.power === 1
        }
        this.state = "loaded";
      } else {
        this.state = "loaded";
      }
    },
  };
</script>
