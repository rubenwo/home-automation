import ApiService from "./api.service";

export default {
  async fetchDevices() {
    const res = await ApiService().get("http://192.168.2.135/api/v1/devices");
    console.log(res);
    return res.data;
  },
  async addNewDevice(data) {
    console.log(data)
    let url = "http://192.168.2.135/api/v1";
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
