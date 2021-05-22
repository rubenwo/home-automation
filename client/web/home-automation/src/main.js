import Vue from "vue";
import App from "./App.vue";
import router from "./router";
import store from "./store";
import BootstrapVue from "bootstrap-vue";
import "bootstrap/dist/css/bootstrap.css";
import "bootstrap-vue/dist/bootstrap-vue.css";
import { library } from "@fortawesome/fontawesome-svg-core";
import {
  faCoffee,
  faChevronDown,
  faChevronLeft,
  faChevronUp,
} from "@fortawesome/free-solid-svg-icons";
import { faJs, faVuejs } from "@fortawesome/free-brands-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";
import VueEvents from "vue-events";

library.add(faCoffee, faChevronDown, faChevronLeft, faChevronUp, faJs, faVuejs);

Vue.component("font-awesome-icon", FontAwesomeIcon);

Vue.config.productionTip = false;
Vue.use(BootstrapVue);
Vue.use(VueEvents);

new Vue({
  router,
  store,
  render: (h) => h(App),
}).$mount("#app");
