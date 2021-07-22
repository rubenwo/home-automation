<template>
    <div>

        <div class="row" style="margin-top: 10px">
            <b-col
                    cols="4"
                    sm="3"
                    md="3"
                    lg="2"
                    xl="2"
                    v-bind:key="camera.id"
                    v-for="camera in this.cameras"
                    style="margin-right: 25px; margin-left: 25px"
            >
                <p>{{camera}}</p>
                <img :src="getStreamUrl(camera.id)" alt="video streaming from the camera">
            </b-col>
        </div>
    </div>
</template>

<script>
  import CameraService from '../services/camera.service';

  export default {
    name: "Cameras",
    data() {
      return {
        cameras: []
      }
    },
    methods: {
      getStreamUrl(cameraId) {
        if (process.env.VUE_APP_BACKEND_URL !== "") {
          return "https://" + process.env.VUE_APP_BACKEND_URL + "/api/v1/stream/" + cameraId;
        } else {
          return "/api/v1/stream/" + cameraId;
        }
      }
    },
    async mounted() {
      const cameras = await CameraService.fetchCameras();
      console.log(cameras);
      this.cameras = cameras;
    },
  };
</script>

<style scoped></style>
