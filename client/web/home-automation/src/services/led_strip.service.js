import ApiService from "./api.service";

export default {
  async fetchAllLedStripDevices() {
    const res = await ApiService().get("/api/v1/leds/devices");
    return res.data;
  },
  async commandLedStripDeviceSolid(deviceId, data) {
    const res = await ApiService().post(
        "/api/v1/leds/devices/" + deviceId + "/command/" + "solid",
        data
    );
    return res.data;
  },
  async commandLedStripDeviceColorCycle(deviceId) {
    const res = await ApiService().post(
        "/api/v1/leds/devices/" + deviceId + "/command/" + "colorcycle");
    return res.data;
  },
  async commandLedStripDeviceChristmas(deviceId) {
    const res = await ApiService().post(
        "/api/v1/leds/devices/" + deviceId + "/command/" + "christmas");
    return res.data;
  },
  async commandLedStripDeviceBreathing(deviceId, data) {
    const res = await ApiService().post(
        "/api/v1/leds/devices/" + deviceId + "/command/" + "breathing",
        data
    );
    return res.data;
  },
};
