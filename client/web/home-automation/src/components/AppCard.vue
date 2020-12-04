<template>
    <b-card
            v-bind:sub-title="name"
            style="max-width: 540px; min-width: 150px; height: 350px"
            class="mb-2">
        <b-card-img v-bind:src="img" alt="Image" height="130" width="130"
                    class="mb-4"></b-card-img>
        <p>{{category}}</p>
        <p>{{company}} : {{device_type}}</p>
        <div slot="footer">
            <b-button style="background-color: #4287f5;" v-bind:to="navigate()">Information
            </b-button>
        </div>
    </b-card>
</template>

<script>

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

    methods: {
      navigate() {
        return "device/" + this.company + "/" + this.id;
      },
    }
  };
</script>

