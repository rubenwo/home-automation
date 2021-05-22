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
      store.dispatch("auth/refreshToken");
      // store.dispatch("auth/logout");
    }
    return err.response;
  });
  return http;
};

// let isAlreadyFetchingAccessToken = false
// let subscribers = []
//
// function onAccessTokenFetched(access_token) {
//   subscribers = subscribers.filter(callback => callback(access_token))
// }
//
// function addSubscriber(callback) {
//   subscribers.push(callback)
// }
//
// axios.interceptors.response.use(function (response) {
//   return response
// }, function (error) {
//   const { config, response: { status } } = error
//   const originalRequest = config
//
//   if (status === 401) {
//     if (!isAlreadyFetchingAccessToken) {
//       isAlreadyFetchingAccessToken = true
//
//       // instead of this store call you would put your code to get new token
//       store.dispatch(fetchAccessToken()).then((access_token) => {
//         isAlreadyFetchingAccessToken = false
//         onAccessTokenFetched(access_token)
//       })
//     }
//
//     const retryOriginalRequest = new Promise((resolve) => {
//       addSubscriber(access_token => {
//         originalRequest.headers.Authorization = 'Bearer ' + access_token
//         resolve(axios(originalRequest))
//       })
//     })
//     return retryOriginalRequest
//   }
//   return Promise.reject(error)
// })
