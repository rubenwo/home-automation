import ApiService from "./api.service";

export default {
  async fetchRecipes() {
    const res = await ApiService()
        .get("/api/v1/recipes")
        .catch(() => {
          return null;
        });
    console.log(res);
    return res.data;
  },
  async addRecipe(data) {
    const res = await ApiService()
        .post(
            "/api/v1/recipes",
            data,
        ).catch(() => {
          return null;
        })
    console.log(res);
    return res.data;
  },
  async fetchRecipe(recipeId) {
    const res = await ApiService()
        .get("/api/v1/recipes/" + recipeId)
        .catch(() => {
          return null;
        });
    console.log(res);
    return res.data;
  },
  async deleteRecipe(recipeId) {
    const res = await ApiService()
        .delete("/api/v1/recipes/" + recipeId)
        .catch(() => {
          return null;
        });
    console.log(res);
    return res.data;
  }
}
