<template>
    <div class="recipe" style="        background-color: rgba(255, 255, 255, 0.25);
        backdrop-filter: blur(5px);">
        <h3>{{ recipe.name }}</h3>
        <b-img v-bind:src="recipe.img" fluid center/>
        <br/>
        <ul>
            <li
                    v-bind:key="'ingredient-' + index"
                    v-for="(ingredient, index) in recipe.ingredients"
            >
                {{ ingredient.name }} :
                {{ ingredient.amount }}
            </li>
        </ul>
        <br/>
        <ul>
            <li v-bind:key="'step-' + index" v-for="(step, index) in recipe.steps">
                {{ step.instruction }}
                <b-button>X</b-button>
            </li>
        </ul>
    </div>
</template>

<script>
  import FoodService from "../services/food.service";

  export default {
    name: "Recipe",
    data() {
      return {
        recipe: {
          id: 0,
          name: "",
          img: "",
          ingredients: [
            {
              id: 0,
              name: "",
              amount: "",
            },
          ],
          steps: [
            {
              id: 0,
              instruction: "",
            },
          ],
        },
      };
    },
    async created() {
      this.id = this.$route.params.id;

      const res = await FoodService.fetchRecipe(this.id);
      console.log(res);
      this.recipe = res.recipe;
    },
    mounted() {
    },
  };
</script>

<style>
    .recipe {
        margin-top: 10px;
        text-align: center;
    }
</style>
