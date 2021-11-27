<template>
    <div>
        <br>
        <input class="form-control" type="text" placeholder="Search" aria-label="Search" v-model="searchInput"/>
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
    </div>
</template>

<script>
  import RecipeCard from "../components/RecipeCard";
  import FoodService from "../services/food.service";

  export default {
    name: "Recipes",
    components: {RecipeCard},
    data() {
      return {
        recipes: [],
        searchInput: ""
      };
    },
    methods: {},
    computed: {
      computedRecipes() {
        return this.recipes.filter(recipe => recipe.name.toLowerCase().includes(this.searchInput.toLowerCase()));
      }
    },
    async mounted() {
      const recipes = await FoodService.fetchRecipes(this.id);
      this.recipes = recipes.recipes;
    },
  };
</script>

<style></style>
