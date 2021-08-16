import ApiService from "./api.service";

export default {
  async fetchDevices() {
    const res = await ApiService()
      .get("/api/v1/devices")
      .catch(() => {
        return null;
      });
    return res.data;
  },
  async addNewDevice(data) {
    let url = "/api/v1";
    switch (data.device_type) {
      case "tapo":
        url += "/tapo/devices/register";
        break;
      case "LED_STRIP":
        url += "/leds/devices/register";
        break;
      // case "hue":
      //   url += "/hue/devices/register";
      //   break;
    }
    const res = await ApiService().post(url, data.data);
    return res.data;
  },
};
