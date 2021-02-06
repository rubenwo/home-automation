<template>
    <div class="recipe" style="background-color: rgba(255, 255, 255, 0.7);">
        <h3>{{recipe.name}}</h3>
        <b-img v-bind:src="recipe.img" fluid center/>
        <br>
        <li v-bind:key="'ingredient-'+index" v-for="(ingredient,index) in recipe.ingredients">{{ingredient.name}} :
            {{ingredient.amount}}
        </li>

        <br>
        <li v-bind:key="'step-'+index" v-for="(step, index) in recipe.steps">
            {{step.instruction}}
            <b-button>X</b-button>
        </li>
    </div>
</template>

<script>
  import FoodService from '../services/food.service';

  export default {
    name: "Recipe",
    data() {
      return {
        recipe: {
          name: "",
          img: "",
          ingredients: [
            {
              name: "",
              amount: "",
            }
          ],
          steps: [
            {
              instruction: "",
            }
          ]
        },

      }
    },
    async created() {
      this.id = this.$route.params.id;

      const res = await FoodService.fetchRecipes();
      for (let recipe of res.recipes) {
        if (this.id === recipe.id) {
          this.recipe = recipe;
          break;
        }
      }
    }
    ,
    mounted() {
    }
  }
</script>

<style>
    .recipe {
        margin-top: 10px;
        text-align: center;
    }
</style>
