import DeviceService from "../../services/devices.service";

export default {
  namespaced: true,
  state: {
    loading: false,
    error: null,
    devices: []
  },
  mutations: {
    REQUEST(state) {
      state.loading = true;
      state.error = null;
    },
    SUCCESS(state, devices) {
      state.loading = false;
      state.devices = devices;
    },
    FAILED(state, message) {
      state.loading = false;
      state.error = message;
    }
  },
  getters: {
    devices: state => state.devices.filter(device => device.status)
  },
  actions: {
    async fetchDevices({commit}) {
      commit("REQUEST");
      try {
        const result = await DeviceService.fetchDevices();
        commit("SUCCESS", result.devices);
      } catch (err) {
        commit("FAILED", err.message);
        throw err;
      }
    },
    async addNewDevice({commit}, data) {
      console.log('addNewDevice')
      console.log(data)
      commit("REQUEST");
      try {
        const result = await DeviceService.addNewDevice(data)
        console.log(result)
      } catch (err) {
        commit("FAILED", err.message);
        throw err;
      }
    }
  }
};
