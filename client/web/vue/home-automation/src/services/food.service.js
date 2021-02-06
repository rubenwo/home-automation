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
  }
}
