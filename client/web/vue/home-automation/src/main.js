import Vue from "vue";
import App from "./App.vue";
import router from "./router";
import store from "./store";
import BootstrapVue from "bootstrap-vue";
import "bootstrap/dist/css/bootstrap.css";
import "bootstrap-vue/dist/bootstrap-vue.css";
import VueGraph from "vue-graph";

Vue.config.productionTip = false;
Vue.use(BootstrapVue);
Vue.use(VueGraph);

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount("#app");
