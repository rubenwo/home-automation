import ApiService from "./api.service";

export default {
  async fetchTapoDevice(deviceId) {
    const result = await ApiService().get(
        "http://192.168.2.135/api/v1/tapo/devices/" + deviceId
    );
    console.log(result);
    return result.data;
  },
  async fetchAllTapoDevices() {
    const result = await ApiService().get(
        "http://192.168.2.135/api/v1/tapo/devices"
    );
    console.log(result);
    return result.data;
  },
  async commandDevice(deviceId, command, brightness) {
    const result = await ApiService().get(
        "http://192.168.2.135/api/v1/tapo/lights/" + deviceId + "?command=" + command + "&brightness=" + brightness
    );
    console.log(result);
    return result.data;
  },
  async turnOnDevice(deviceId) {
    return await this.commandDevice(deviceId, "on", 100);
  },
  async turnOffDevice(deviceId) {
    return await this.commandDevice(deviceId, "off", 0);
  },
  async setDeviceBrightness(deviceId, brightness) {
    return await this.commandDevice(deviceId, "on", brightness);
  }
};
