import TapoService from "../../services/tapo.service";

export default {
  namespaced: true,
  state: {
    loading: false,
    error: null,
    tapoDevice: null
  },
  mutations: {
    REQUEST(state) {
      state.loading = true;
      state.error = null;
    },
    TAPO_DEVICE_LOADED(state, device) {
      state.loading = false;
      state.tapoDevice = device;
    },
    FAILED(state, message) {
      state.loading = false;
      state.error = message;
    }
  },
  getters: {
    tapoDevice: state => state.tapoDevice
  },
  actions: {
    async fetchTapoDevice({commit}, deviceId) {
      commit("REQUEST");
      try {
        const result = await TapoService.fetchTapoDevice(deviceId);
        console.log(result.device)
        commit("TAPO_DEVICE_LOADED", result.device);
      } catch (err) {
        commit("FAILED", err.message);
        throw err;
      }
    }
  }
}
