import Vue from "vue";
import Vuex from "vuex";

import DeviceModule from "./modules/devices.module";
import TapoModule from "./modules/tapo.module";
import AuthModule from "./modules/auth.module";

Vue.use(Vuex);

export default new Vuex.Store({
  state: {},
  mutations: {},
  actions: {},
  modules: {
    devices: DeviceModule,
    tapo: TapoModule,
    auth: AuthModule,
  },
});
