import ApiService from "./api.service";

export default {
  async fetchAllLedStripDevices() {
    const res = await ApiService().get("/api/v1/leds/devices");
    console.log(res);
    return res.data;
  },
  async commandLedStripDevice(deviceId, data) {
    const res = await ApiService().post(
      "/api/v1/leds/devices/" + deviceId + "/command",
      data
    );
    console.log(res);
    return res.data;
  }
};
