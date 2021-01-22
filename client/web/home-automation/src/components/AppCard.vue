<template>
    <div v-if="this.state === 'loaded'">
        <b-card
                v-bind:sub-title="name"
                style="max-width: 540px; min-width: 200px; min-height: 425px; max-height: 500px"
                class="mb-2">
            <b-card-img v-bind:src="img" alt="Image" height="130" width="130"
                        class="mb-4"/>
            <p>{{category}}</p>
            <p>{{company}} : {{device_type}}</p>
            <div v-if="company==='tp-link'">
                <b-button variant="success" @click="turnOnDevice">On</b-button>
                <b-button @click="turnOffDevice()">Off</b-button>
                <input v-if="device_type=='L510E'" type="range" min="1" max="100"
                       v-model="brightness" @change="brightnessChanged()">
            </div>
            <div slot="footer">
                <b-button style="background-color: #4287f5;" v-bind:to="navigate()">Information
                </b-button>
                <b-button variant="danger" @click="deleteDevice()">X</b-button>
            </div>
        </b-card>
    </div>
    <div v-else>
        <h3>Loading results...</h3>
        <Loading :active.sync="this.state === 'loading'"
                 :is-full-page="true"/>
    </div>
</template>

<script>
  import TapoService from "../services/tapo.service"
  import {mapActions, mapState} from "vuex";
  import Loading from 'vue-loading-overlay'

  export default {
    name: "app-card",
    components: {Loading},
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
      }
    },
    computed: {
      ...mapState("tapo", {
        tapoDevice: state => state.tapoDevice
      })
    },
    methods: {
      ...mapActions("tapo", ["wakeTapoDevice", "fetchTapoDevice"]),
      navigate() {
        return "device/" + this.company + "/" + this.id;
      },
      turnOnDevice() {
        if (this.company === "tp-link") {
          console.log(this.id, this.brightness)
          if (this.device_type === 'L510E') {
            TapoService.setDeviceBrightness(this.id, this.brightness)
          } else
            TapoService.turnOnDevice(this.id);
        }
      },
      turnOffDevice() {
        if (this.company === "tp-link") {
          TapoService.turnOffDevice(this.id);
        }
      },
      async deleteDevice() {
        console.log(this.id)
        if (this.company === "tp-link") {
          const res = await TapoService.deleteTapoDevice(this.id);
          console.log(res)
        }
      },
      async brightnessChanged() {
        console.log(this.brightness)
        await TapoService.setDeviceBrightness(this.id, this.brightness)
      },
    },
    async mounted() {
      console.log(this.company)
      if (this.company === "tp-link") {
        await this.wakeTapoDevice(this.id);
        await this.fetchTapoDevice(this.id);
        this.state = 'loaded'
      } else {
        this.state = 'loaded'
      }
    }
  };
</script>

