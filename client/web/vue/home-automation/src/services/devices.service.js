import ApiService from "./api.service";

export default {
  async fetchDevices() {
    console.log(ApiService().baseUrl)
    const res = await ApiService()
      .get("/api/v1/devices")
      .catch(() => {
      return null;
    });
    console.log(res);
    return res.data;
  },
  async addNewDevice(data) {
    console.log(data)
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
    console.log(res)
    return res.data;
  }
};
