<template>
    <div>
        <b-button pill @click="onButtonClick" center>
            <verte v-model="color" picker="wheel" model="rgb"/>
        </b-button>
    </div>
</template>

<script>
  import Verte from 'verte';
  import 'verte/dist/verte.css';

  export default {
    name: "LedStripDevice",
    components: {Verte},
    methods: {
      onButtonClick() {
        console.log(this.color)
      }
    },
    data() {
      return {
        color: "",
      }
    }
  }

</script>


<style>

</style>
