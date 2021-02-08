import ApiService from './api.service';

export default {
  async login(username, password) {
    const res = await ApiService().post("/auth/login", {
      username: username,
      password: password
    }).catch(() => {
      return null;
    });
    return res;
  }
}
