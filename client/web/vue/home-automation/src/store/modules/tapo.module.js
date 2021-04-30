import TapoService from "../../services/tapo.service";

export default {
  namespaced: true,
  state: {
    loading: false,
    error: null,
    tapoDevice: null,
    tapoDevices: [],
    devs: {},
  },
  mutations: {
    REQUEST(state) {
      state.loading = true;
      state.error = null;
    },
    WAKE_DEVICE() {},
    TAPO_DEVICE_LOADED(state, device) {
      state.loading = false;
      state.tapoDevice = device;
    },
    TAPO_DEVICES_LOADED(state, devices) {
      state.loading = false;
      state.tapoDevices = devices;
    },
    FAILED(state, message) {
      state.loading = false;
      state.error = message;
    },
    DEV_LOADED(state, dev) {
      state.loading = false;
      state.devs[dev.device_id] = dev;
      console.log(state.devs);
    },
  },
  getters: {
    tapoDevice: (state) => state.tapoDevice,
  },
  actions: {
    async fetchTapoDevice({ commit }, deviceId) {
      commit("REQUEST");
      try {
        const result = await TapoService.fetchTapoDevice(deviceId);
        console.log(result.device);
        commit("TAPO_DEVICE_LOADED", result.device);
        commit("DEV_LOADED", result.device);
      } catch (err) {
        commit("FAILED", err.message);
        throw err;
      }
    },
    async fetchTapoDevices({ commit }) {
      commit("REQUEST");
      try {
        const result = await TapoService.fetchAllTapoDevices();
        console.log(result.devices);
        commit("TAPO_DEVICES_LOADED", result.devices);
      } catch (err) {
        commit("FAILED", err.message);
        throw err;
      }
    },
    async wakeTapoDevice({ commit }, deviceId) {
      commit("WAKE_DEVICE");
      const result = await TapoService.wakeTapoDevice(deviceId);
      console.log(result.devices);
    },
  },
};
