<template>
  <div class="row" style="margin-top: 10px">
    <b-col
      cols="4"
      sm="3"
      md="3"
      lg="2"
      xl="2"
      v-bind:key="recipe.id"
      v-for="recipe in this.recipes"
      style="margin-right: 25px; margin-left: 25px"
    >
      <recipe-card
        v-bind:name="recipe.name"
        v-bind:img="recipe.img"
        v-bind:id="recipe.id"
      />
    </b-col>
  </div>
</template>

<script>
import RecipeCard from "../components/RecipeCard";
import FoodService from "../services/food.service";

export default {
  name: "Recipes",
  components: { RecipeCard },
  data() {
    return {
      recipes: [],
    };
  },
  methods: {},
  async mounted() {
    const recipes = await FoodService.fetchRecipes(this.id);
    console.log(recipes);
    this.recipes = recipes.recipes;
  },
};
</script>

<style></style>
