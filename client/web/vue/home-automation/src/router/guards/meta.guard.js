import store from "@/store";

export default async (to, from, next) => {
  const loggedIn = store.getters["auth/isLoggedIn"];
  //  const user = store.getters['auth/getUser']
  //console.log(user)
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth);

  console.log("META GUARD");
  console.log(loggedIn);
  console.log(requiresAuth);

  // eslint-disable-next-line no-console
  console.log("[MetaGuard]: To:", to);
  // eslint-disable-next-line no-console
  if (!loggedIn && requiresAuth) {
    return next({
      name: "login",
      query: { redirect: to.path }
    });
  }

  if (loggedIn && to.name === "login") {
    return next();
  }

  return next();
};
