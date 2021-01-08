<template>
    <b-card
            v-bind:sub-title="name"
            style="max-width: 540px; min-width: 175px; min-height: 425px; max-height: 500px"
            class="mb-2">
        <b-card-img v-bind:src="img" alt="Image" height="130" width="130"
                    class="mb-4"/>
        <p>{{category}}</p>
        <p>{{company}} : {{device_type}}</p>
        <div>
            <b-button variant="success" @click="turnOnDevice">On</b-button>
            <b-button @click="turnOffDevice()">Off</b-button>
            <input v-if="device_type=='L510E'" type="range" min="1" max="100"
                   v-model="brightness">
        </div>
        <div slot="footer">
            <b-button style="background-color: #4287f5;" v-bind:to="navigate()">Information
            </b-button>
        </div>
    </b-card>
</template>

<script>
  import TapoService from "../services/tapo.service"

  export default {
    name: "app-card",
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
      return {brightness: 100}
    },
    methods: {
      navigate() {
        return "device/" + this.company + "/" + this.id;
      },
      turnOnDevice() {
        console.log(this.id, this.brightness)
        if (this.device_type === 'L510E')
          TapoService.setDeviceBrightness(this.id, this.brightness)
        else
          TapoService.turnOnDevice(this.id);
      },
      turnOffDevice() {
        TapoService.turnOffDevice(this.id);
      }
    }
  };
</script>

