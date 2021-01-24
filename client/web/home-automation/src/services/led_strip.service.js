import ApiService from "./api.service";

export default {
  async fetchAllLedStripDevices() {
    const res = await ApiService().get("http://192.168.2.135/api/v1/leds/devices");
    console.log(res);
    return res.data;
  },
  async commandLedStripDevice(deviceId, data) {
    const res = await ApiService().post("http://192.168.2.135/api/v1/leds/devices/" + deviceId + "/command", data);
    console.log(res);
    return res.data;
  }
};
