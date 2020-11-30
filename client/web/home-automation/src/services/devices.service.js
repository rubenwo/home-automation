import ApiService from "./api.service";

export default {
  async fetchDevices() {
    const res = await ApiService().get("http://192.168.2.135/api/v1/devices");
    console.log(res);
    return res.data;
  }
};
