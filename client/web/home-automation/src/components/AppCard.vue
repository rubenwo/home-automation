<template>
    <b-card
            v-bind:sub-title="title"
            style="max-width: 540px; min-width: 150px; height: 350px"
            class="mb-2">
        <b-card-img v-bind:src="img" alt="Image" height="130" width="130" class="mb-4"></b-card-img>
        <b-progress v-bind:value="progress" :max="100" show-progress class="mb-3"></b-progress>
        <div slot="footer">
            <b-button style="background-color: #c6002a;" @click="updateProgress" v-bind:to="navigate()">Practice!
            </b-button>
        </div>
    </b-card>
</template>

<script>
  import {mapActions} from "vuex";

  export default {
    name: "app-card",
    props: {
      title: {
        type: String,
        default: ""
      },
      img: {
        type: String,
        default: ""
      },
      progress: {
        type: Number,
        default: 0
      },
      id: {
        type: Number,
        default: 0
      }
    },
    methods: {
      navigate() {
        return "module/" + this.id;
      },
      ...mapActions("auth", ["updateProgress"])
    }
  };
</script>

