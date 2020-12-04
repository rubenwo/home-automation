import ApiService from "./api.service";

export default {
  async fetchTapoDevice(deviceId) {
    const result = await ApiService().get(
      "http://192.168.2.135/api/v1/tapo/devices/" + deviceId
    );
    console.log(result);
    return result.data;
  }
};
