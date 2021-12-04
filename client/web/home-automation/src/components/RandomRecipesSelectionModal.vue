<template>
  <b-modal
    size="xl"
    id="randomRecipeModal"
    ref="randomRecipeModal"
    @ok="handleOk"
    @cancel="handleCancel"
    @close="handleCancel"
  >
    Hello World

    <div v-bind:key="recipe.id" v-for="recipe in recipes">
      <recipe-card
        v-bind:name="recipe.name"
        v-bind:img="recipe.img"
        v-bind:id="recipe.id"
      />
    </div>
  </b-modal>
</template>

<script>
import FoodService from "../services/food.service";
import RecipeCard from "./RecipeCard";

export default {
  name: "RandomRecipeModal",
  components:{RecipeCard},
  data() {
    return {
      recipes: [],
    };
  },
  methods: {
    handleOk(evt) {
      evt.preventDefault();
      this.handleSubmit();
    },
    handleCancel() {},
    async handleSubmit() {
      this.$nextTick(() => {
        this.$refs.randomRecipeModal.hide();
      });
    },
  },
  created() {
    this.$on("randomize_recipes", async () => {
      this.recipes = await FoodService.fetchRandomRecipes(2);
      this.$refs.randomRecipeModal.show();
    });
  },
};
</script>

<style scoped></style>
