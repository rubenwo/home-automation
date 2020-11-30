import axios from "axios";
import store from "@/store";

export default () => {
  const http = axios.create({
    baseUrl: process.env.VUE_APP_BACKEND_URL,
    withCredentials: false
  });

  http.interceptors.response.use(undefined, async err => {
    console.log("Hi There", err.response.status);
    if (err.response.status === 401) {
      store.dispatch("");
    }
    return err.response;
  });
  return http;
};
