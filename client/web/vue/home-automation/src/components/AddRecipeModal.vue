<template>
    <b-modal size="xl" id="recipeModal" ref="recipeModal" @ok="handleOk" @cancel="handleCancel" @close="handleCancel">
        <b-row>
            <div>
                <p>Recipe:</p>
            </div>
        </b-row>
        <b-form-group>
            <b-input-group>
                <b-container fluid>
                    <b-row class="my-1">
                        <b-col sm="4">
                            <label>Name:</label>
                        </b-col>
                        <b-col sm="8">
                            <input size="sm" class="mx-1" placeholder="Recipe Name"
                                   v-model="recipe.name"/>
                        </b-col>
                    </b-row>
                </b-container>
            </b-input-group>
        </b-form-group>
        <b-form-group>
            <b-input-group>
                <b-container fluid>
                    <b-row class="my-1">
                        <b-col sm="4">
                            <label>Image Source:</label>
                        </b-col>
                        <b-col sm="8">
                            <input size="sm" class="mx-1" placeholder="http://..."
                                   v-model="recipe.img"/>
                        </b-col>
                    </b-row>
                </b-container>
            </b-input-group>
        </b-form-group>
        <b-row>
            <div>
                <p>Ingredients:</p>
            </div>
        </b-row>
        <b-row v-bind:key="'ingredients'+index" v-for="(ingredient, index) in recipe.ingredients">
            <b-input-group>
                <b-container fluid>
                    <b-row class="my-1">
                        <b-col sm="4">
                            <label>Ingredient name:</label>
                        </b-col>
                        <b-col sm="8">
                            <b-form-input size="sm" class="mx-1" placeholder="name"
                                          v-model="ingredient.name"/>
                        </b-col>
                    </b-row>
                </b-container>
            </b-input-group>
            <b-input-group>
                <b-container fluid>
                    <b-row class="my-1">
                        <b-col sm="4">
                            <label>Ingredient amount:</label>
                        </b-col>
                        <b-col sm="8">
                            <b-form-input size="sm" class="mx-1" placeholder="amount"
                                          v-model="ingredient.amount"/>
                        </b-col>
                    </b-row>
                </b-container>
            </b-input-group>
        </b-row>
        <b-button variant="success" @click="increaseIngredients">+</b-button>
        <b-button variant="danger" @click="decreaseIngredients">-</b-button>
        <b-row>
            <div>
                <p>Steps:</p>
            </div>
        </b-row>
        <b-row v-bind:key="'steps-'+index" v-for="(step, index) in recipe.steps">
            <b-input-group>
                <b-container fluid>
                    <b-row class="my-1">
                        <b-col sm="4">
                            <label>Instruction:</label>
                        </b-col>
                        <b-col sm="8">
                            <b-form-input size="sm" class="mx-1" placeholder="instruction"
                                          v-model="step.instruction"/>
                        </b-col>
                    </b-row>
                </b-container>
            </b-input-group>

        </b-row>
        <b-button variant="success" @click="increaseSteps">+</b-button>
        <b-button variant="danger" @click="decreaseSteps">-</b-button>

    </b-modal>

</template>

<script>
  import FoodService from '../services/food.service';

  export default {
    name: "AddRecipeModal",
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
    methods: {
      increaseSteps() {
        this.recipe.steps.push({
          instruction: "",
        })
      },
      decreaseSteps() {
        this.recipe.steps.pop();
      },
      increaseIngredients() {
        this.recipe.ingredients.push({
          name: "",
          amount: "",
        })
      },
      decreaseIngredients() {
        this.recipe.ingredients.pop();
      },
      handleOk(evt) {
        evt.preventDefault();
        this.handleSubmit();
      },
      handleCancel() {
      },
      async handleSubmit() {
        console.log(this.recipe)

        const res = await FoodService.addRecipe(this.recipe);
        console.log(res);

        this.recipe = {
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
        };


        this.$nextTick(() => {
          this.$refs.recipeModal.hide();
        });
      }
    },
    created() {
      this.$on('add_recipe', () => {
        this.$refs.recipeModal.show()
      })
    },
  }
</script>

<style scoped>

</style>
