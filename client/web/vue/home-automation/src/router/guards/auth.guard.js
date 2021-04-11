import store from "@/store";

export default async (to, from, next) => {
  const loggedIn = store.getters["auth/isLoggedIn"];

  console.log("AUTH GUARD");
  if (!loggedIn) {
    return next({ name: "login", query: { redirect: to.path } });
  }

  store.watch(
    () => store.getters["auth/isLoggedIn"],
    loggedIn => {
      if (!loggedIn) {
        return next({ name: "login", query: { redirect: to.path } });
      }
    }
  );

  return next();
};
