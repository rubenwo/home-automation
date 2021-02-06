<template>
    <b-card
            v-bind:sub-title="name"
            style="max-width: 540px; min-width: 200px; min-height: 425px; max-height: 500px; background-color: rgba(255, 255, 255, 0.7)"
            class="mb-2">
        <b-card-img v-bind:src="img" alt="Image" height="130" width="130"
                    class="mb-4"/>

        <div slot="footer">
            <b-button style="background-color: #4287f5;" v-bind:to="navigate()">Information
            </b-button>
            <!--            <b-button variant="danger" @click="deleteDevice()">X</b-button>-->
        </div>
    </b-card>
</template>

<script>
  import FoodService from '../services/food.service';

  export default {
    name: "RecipeCard",
    props: {
      id: {
        type: String,
        default: ""
      },
      name: {
        type: String,
        default: ""
      },
      img: {
        type: String,
        default: ""
      },
    },
    methods: {
      navigate() {
        return "recipe/" + this.id;
      },
    },
    async mounted() {
      const recipes = await FoodService.fetchRecipes(this.id);
      console.log(recipes);
    }
  }
</script>

<style scoped>

</style>
