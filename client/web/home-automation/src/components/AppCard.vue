<template>
    <b-card
            v-bind:sub-title="title"
            style="max-width: 540px; min-width: 150px; height: 350px"
            class="mb-2">
        <b-card-img v-bind:src="img" alt="Image" height="130" width="130" class="mb-4"></b-card-img>
        <div slot="footer">
            <b-button style="background-color: #c6002a;" v-bind:to="navigate()">Practice!
            </b-button>
        </div>
    </b-card>
</template>

<script>

  export default {
    name: "app-card",
    props: {
      title: {
        type: String,
        default: ""
      },
      img: {
        type: String,
        default: "https://labapp.rubenwoldhuis.nl/static/images/dangers/harmful_hazard.png"
      },
      id: {
        type: String,
        default: ""
      }
    },
    methods: {
      navigate() {
        return "device/" + this.id;
      },
    }
  };
</script>

