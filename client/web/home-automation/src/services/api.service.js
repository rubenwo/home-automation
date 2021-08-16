import axios from "axios";
import store from "@/store";

export default () => {
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
    console.error(err);
    if (err.response.status === 401) {
      store.dispatch("auth/refreshToken");
    }
    return err.response;
  });
  return http;
};
