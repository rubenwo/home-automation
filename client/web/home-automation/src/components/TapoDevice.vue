<template>
    <div>
        <p>{{tapoDevice.device_name}}</p>

        <div v-if="this.tapoDevice.device_type === 'L510E'">
            <p>LIGHT</p>
        </div>
        <div v-else-if="this.tapoDevice.device_type === 'P100'">
            <p>PLUG</p>
        </div>
    </div>
</template>

<script>
  import {mapActions, mapState} from "vuex";

  export default {
    name: "TapoDevice",
    id: {
      type: String,
      default: ""
    },
    computed: {
      ...mapState("tapo", {
        tapoDevice: state => state.tapoDevice
      })
    },
    methods: {
      ...mapActions("tapo", ["fetchTapoDevice"]),
    },
    async mounted() {
      this.id = this.$route.params.id;
      await this.fetchTapoDevice(this.id);
      console.log(this.tapoDevice)
    },
  }
</script>

<style scoped>

</style>
