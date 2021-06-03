import Vue from "vue";
import VueRouter from "vue-router";
import Home from "../views/Home.vue";
import Login from "../views/Login.vue";
import {MetaGuard} from "./guards";
import store from "@/store";

Vue.use(VueRouter);

const routes = [
  {
    path: "/",
    name: "Home",
    component: Home,
    meta: {
      requiresAuth: true,
      title: "Home",
    },
  },
  {
    path: "/device/:company/:id",
    name: "Device",
    component: () => import("../views/Device.vue"),
    meta: {
      requiresAuth: true,
      title: "Device",
    },
  },
  {
    path: "/recipes",
    name: "Recipes",
    component: () => import("../views/Recipes.vue"),
    meta: {
      requiresAuth: true,
      title: "Recipes",
    },
  },
  {
    path: "/sensors",
    name: "Sensors",
    component: () => import("../views/Sensors.vue"),
    meta: {
      requiresAuth: true,
      title: "Sensors",
    },
  },
  {
    path: "/routines",
    name: "Routines",
    component: () => import("../views/Routines.vue"),
    meta: {
      requiresAuth: true,
      title: "Routines",
    },
  },
  {
    path: "/recipe/:id",
    name: "Recipe",
    component: () => import("../views/Recipe"),
    meta: {
      requiresAuth: true,
      title: "Recipe",
    },
  },
  {
    path: "/routine/:id",
    name: "Routine",
    component: () => import("../views/Routine"),
    meta: {
      requiresAuth: true,
      title: "Routine",
    },
  },
  {
    path: "/inventory",
    name: "Inventory",
    component: () => import("../views/Inventory"),
    meta: {
      requiresAuth: true,
      title: "Inventory",
    },
  },
  {
    path: "/cameras",
    name: "Cameras",
    component: () => import("../views/Cameras"),
    meta: {
      requiresAuth: true,
      title: "Cameras",
    },
  },
  {
    path: "/login",
    name: "login",
    component: Login,
    meta: {
      title: "Login",
    },
  },
];

const router = new VueRouter({
  routes,
});

router.beforeEach(MetaGuard);
// eslint-disable-next-line no-unused-vars
router.afterEach((to, from) => {
  const pageTitle = to.meta.title;

  if (pageTitle) {
    document.title = `Home Automation - ${pageTitle}`;
  } else {
    document.title = "Home Automation";
  }
});

store.watch(
    (state, getters) => getters["auth/isLoggedIn"],
    (loggedIn) => {
      if (
          !loggedIn &&
          router.currentRoute.matched.some((record) => record.meta.requiresAuth)
      ) {
        router.push({
          name: "login",
          query: {redirect: router.currentRoute.path},
        });
      }
    }
);

export default router;
