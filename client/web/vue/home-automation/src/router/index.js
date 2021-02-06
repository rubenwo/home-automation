import Vue from "vue";
import VueRouter from "vue-router";
import Home from "../views/Home.vue";

Vue.use(VueRouter);

const routes = [
  {
    path: "/",
    name: "Home",
    component: Home
  },
  {
    path: "/device/:company/:id",
    name: "Device",
    component: () => import("../views/Device.vue")
  },
  {
    path: "/recipes",
    name: "Recipes",
    component: () => import("../views/Recipes.vue")
  },
  {
    path: "/recipe/:id",
    name: "Recipe",
    component: () => import("../views/Recipe")
  }
];

const router = new VueRouter({
  routes
});

export default router;
