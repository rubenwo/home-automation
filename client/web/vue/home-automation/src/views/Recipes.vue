<template>
    <div style="margin-top: 10px">
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

  export default {
    name: "Recipes",
    components: {RecipeCard},
    data() {
      return {
        recipes: [
          {
            id: "0001",
            name: "Pasta",
            img: this.getImgUrl()
          }
        ],
      }
    },
    methods: {
      getImgUrl() {
        return require('../assets/pasta_test_img.jpeg')
      }
    }
  }
</script>

<style>

</style>
