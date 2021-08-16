import ApiService from "./api.service";

export default {
  async wakeTapoDevice(deviceId) {
    const result = await ApiService().get("/api/v1/tapo/wake/" + deviceId);
    return result.data;
  },
  async fetchTapoDevice(deviceId) {
    const result = await ApiService().get("/api/v1/tapo/devices/" + deviceId);
    return result.data;
  },
  async deleteTapoDevice(deviceId) {
    const result = await ApiService().delete(
      "/api/v1/tapo/devices/" + deviceId
    );
    return result.data;
  },
  async fetchAllTapoDevices() {
    const result = await ApiService().get("/api/v1/tapo/devices");
    return result.data;
  },
  async commandDevice(deviceId, command, brightness) {
    const result = await ApiService().get(
      "/api/v1/tapo/lights/" +
        deviceId +
        "?command=" +
        command +
        "&brightness=" +
        brightness
    );
    return result.data;
  },
  async turnOnDevice(deviceId) {
    return await this.commandDevice(deviceId, "on", 100);
  },
  async turnOffDevice(deviceId) {
    return await this.commandDevice(deviceId, "off", 1);
  },
  async setDeviceBrightness(deviceId, brightness) {
    return await this.commandDevice(deviceId, "on", brightness);
  },
};
