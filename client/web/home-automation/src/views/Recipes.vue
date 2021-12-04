<template>
  <div>
    <br />
    <b-button style="width: 100%" @click="randomizeSelection"
      >Randomize</b-button
    >
    <input
      class="form-control"
      type="text"
      placeholder="Search"
      aria-label="Search"
      v-model="searchInput"
    />
    <b-row style="margin-top: 10px">
      <b-col
        cols="4"
        sm="3"
        md="3"
        lg="2"
        xl="2"
        v-bind:key="recipe.id"
        v-for="recipe in this.computedRecipes"
        style="margin-right: 25px; margin-left: 25px"
      >
        <recipe-card
          v-bind:name="recipe.name"
          v-bind:img="recipe.img"
          v-bind:id="recipe.id"
        />
      </b-col>
    </b-row>
    <random-recipe-modal ref="randomRecipeModal" />
  </div>
</template>

<script>
import RecipeCard from "../components/RecipeCard";
import FoodService from "../services/food.service";
import RandomRecipeModal from "../components/RandomRecipesSelectionModal.vue";
export default {
  name: "Recipes",
  components: { RecipeCard, RandomRecipeModal },
  data() {
    return {
      recipes: [],
      searchInput: "",
    };
  },
  methods: {
    randomizeSelection() {
      console.log("clicked randomize button");
      this.$refs.randomRecipeModal.$emit("randomize_recipes");
    },
  },
  computed: {
    computedRecipes() {
      return this.recipes.filter((recipe) =>
        recipe.name.toLowerCase().includes(this.searchInput.toLowerCase())
      );
    },
  },
  async mounted() {
    const recipes = await FoodService.fetchRecipes(this.id);
    this.recipes = recipes.recipes;
  },
};
</script>

<style></style>
