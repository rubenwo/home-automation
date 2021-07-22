import ApiService from "./api.service";

export default {
  async fetchCameras() {
    const res = await ApiService()
        .get("/api/v1/cameras")
        .catch(() => {
          return null;
        });
    console.log(res);
    return res.data;
  },
  async addCamera(data) {
    const res = await ApiService()
        .post("/api/v1/cameras", data)
        .catch(() => {
          return null;
        });
    console.log(res);
    return res.data;
  },
  async updateCamera(cameraId, data) {
    const res = await ApiService()
        .put("/api/v1/cameras/" + cameraId, data)
        .catch(() => {
          return null;
        });
    console.log(res);
    return res.data;
  },
  async fetchCamera(cameraId) {
    const res = await ApiService()
        .get("/api/v1/cameras/" + cameraId)
        .catch(() => {
          return null;
        });
    console.log(res);
    return res.data;
  },
  async deleteCamera(cameraId) {
    const res = await ApiService()
        .delete("/api/v1/cameras/" + cameraId)
        .catch(() => {
          return null;
        });
    console.log(res);
    return res.data;
  }
};
