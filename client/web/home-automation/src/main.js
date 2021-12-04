import Vue from "vue";
import App from "./App.vue";
import router from "./router";
import store from "./store";
import BootstrapVue from "bootstrap-vue";
import "bootstrap/dist/css/bootstrap.css";
import "bootstrap-vue/dist/bootstrap-vue.css";

import VueEvents from "vue-events";
import VueAuthImage from 'vue-auth-image';
import axios from 'axios';


// library.add(faCoffee, faChevronDown, faChevronLeft, faChevronUp, faJs, faVuejs);

// Vue.component("font-awesome-icon", FontAwesomeIcon);

Vue.config.productionTip = false;
Vue.use(BootstrapVue);
Vue.use(VueEvents);
// register vue-auth-image directive
Vue.use(VueAuthImage);

axios.defaults.headers.common['Authorization'] = `Bearer ${store.getters["auth/getBearerToken"]}`;

new Vue({
  router,
  store,
  render: (h) => h(App),
}).$mount("#app");
