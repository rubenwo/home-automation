import ApiService from './api.service';

export default {
  async fetchInventory() {
    const res = await ApiService()
        .get("/api/v1/inventory")
        .catch(() => {
          return null;
        });
    console.log(res);
    return res.data;
  },
  async addInventoryItem(data) {
    const res = await ApiService()
        .post("/api/v1/inventory", data)
        .catch(() => {
          return null;
        });
    console.log(res);
    return res.data;
  },
  async deleteInventoryItem(inventoryItemId) {
    const res = await ApiService()
        .delete("/api/v1/inventory/" + inventoryItemId)
        .catch(() => {
          return null;
        });
    console.log(res);
    return res.data;
  }
}
