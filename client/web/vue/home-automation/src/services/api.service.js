import axios from "axios";
import store from "@/store";

export default () => {
  //console.log(process.env.VUE_APP_BACKEND_URL);
  const http = axios.create({
    baseURL: "https://" + process.env.VUE_APP_BACKEND_URL,
    withCredentials: true,
    headers: {
      ...(store.getters["auth/getBearerToken"] && {
        Authorization: `Bearer ${store.getters["auth/getBearerToken"]}`,
      }),
    },
  });
  http.interceptors.response.use(undefined, async (err) => {
    console.log("Hi There", err.response.status);
    if (err.response.status === 401) {
      store.dispatch("auth/logout");
    }
    return err.response;
  });
  return http;
};
