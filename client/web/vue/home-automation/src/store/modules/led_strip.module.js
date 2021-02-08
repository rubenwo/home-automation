import LedStripService from "../../services/led_strip.service";

export default {
  namespaced: true,
  state: {
    loading: false,
    error: null,
    ledStripDevices: [],
    commandMessage: ""
  },
  mutations: {
    REQUEST(state) {
      state.loading = true;
      state.error = null;
    },
    LED_STRIP_DEVICE_COMMANDED(state, msg) {
      state.loading = false;
      state.commandMessage = msg;
    },
    LED_STRIP_DEVICES_LOADED(state, devices) {
      state.loading = false;
      state.tapoDevices = devices;
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
    async fetchAllLedStripDevices({ commit }) {
      commit("REQUEST");
      try {
        const result = await LedStripService.fetchAllLedStripDevices();
        console.log(result.devices);
        commit("LED_STRIP_DEVICES_LOADED", result.devices);
      } catch (err) {
        commit("FAILED", err.message);
        throw err;
      }
    },
    async commandLedStripDevice({ commit }, data) {
      commit("REQUEST");
      try {
        const result = await LedStripService.commandLedStripDevice(
          data.deviceId,
          data.command
        );
        console.log(result);
        commit("LED_STRIP_DEVICE_COMMANDED", result.message);
      } catch (err) {
        commit("FAILED", err.message);
        throw err;
      }
    }
  }
};
