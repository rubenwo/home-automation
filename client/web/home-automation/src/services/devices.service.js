import ApiService from "./api.service";

export default {
  async fetchDevices() {
    const res = await ApiService().get("http://192.168.2.135/api/v1/devices");
    console.log(res);
    return res.data;
  },
  async addNewDevice(device_type, data) {
    console.log(device_type, data)
    let url = "http://192.168.2.135/api/v1";
    switch (device_type) {
      case "tapo":
        url += "/tapo/devices/register";
        break;
        // case "hue":
        //   url += "/hue/devices/register";
        //   break;
    }
    const res = await ApiService.post(url, data);
    console.log(res)
    return res.data;
  }
};
