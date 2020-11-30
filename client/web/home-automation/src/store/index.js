import Vue from "vue";
import Vuex from "vuex";

import DeviceModule from "./modules/devices.module";

Vue.use(Vuex);

export default new Vuex.Store({
  state: {},
  mutations: {},
  actions: {},
  modules: {
    devices: DeviceModule,
  }
});
