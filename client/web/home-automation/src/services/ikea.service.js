import ApiService from "./api.service";

export default {
  async fetchTradfriDevices() {
    const res = await ApiService()
        .get("/api/v1/tradfri/devices")
        .catch(() => {
          return null;
        });
    return res.data;
  },
  async fetchTradfriDevice(id) {
    const res = await ApiService()
        .get("/api/v1/tradfri/devices/" + id)
        .catch(() => {
          return null;
        });
    return res.data;
  },
  async commandTradfriDevices(command) {
    const res = await ApiService()
        .post("/api/v1/tradfri/devices/command", command)
        .catch(() => {
          return null;
        });
    return res.data;
  },
  async commandTradfriDevice(id, command) {
    const res = await ApiService()
        .post("/api/v1/tradfri/devices/" + id + "/command", command)
        .catch(() => {
          return null;
        });
    return res.data;
  },
};
