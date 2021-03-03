import AuthService from "../../services/auth.service";

export default {
  namespaced: true,
  modules: {},
  state: {
    token: localStorage.getItem("token") || null,
    username: localStorage.getItem("username") || null,
    userid: localStorage.getItem("userid") || null,
    user: null,
    error: null
  },
  mutations: {
    SET_TOKEN(state, token) {
      state.token = token;
    },
    SET_USERNAME(state, username) {
      state.username = username;
    },
    SET_ID(state, userid) {
      state.userid = userid;
    },
    CLEAR_ERROR(state) {
      state.error = null;
    },
    CLEAR_ALL(state) {
      state.token = null;
      state.id = null;
      state.user = null;
      state.username = null;
    }
  },
  actions: {
    async login({ commit }, { username, password }) {
      commit("CLEAR_ERROR");
      console.log("logging in");
      const resp = await AuthService.login(username, password);
      console.log(resp);
      if (resp.status === 200) {
        const { username, user_id, token } = resp.data;
        localStorage.setItem("token", token);
        localStorage.setItem("username", username);
        localStorage.setItem("userid", user_id);

        commit("SET_USERNAME", username);
        commit("SET_TOKEN", token);
        commit("SET_ID", user_id);
      }
    },
    async logout({ commit }) {
      localStorage.removeItem("token");
      localStorage.removeItem("username");
      localStorage.removeItem("userid");
      commit("CLEAR_ALL");
      commit("CLEAR_ERROR");
      let data = {
        username: "admin",
        password: process.env.VUE_APP_ADMIN_PWD
      };
      await this.login(data);
    }
  },
  getters: {
    isLoggedIn: state => {
      if (state.token == null) return false;
      const parseJwt = token => {
        try {
          return JSON.parse(atob(token.split(".")[1]));
        } catch (e) {
          return null;
        }
      };
      let parsedToken = parseJwt(state.token);
      if (parsedToken == null) {
        // eslint-disable-next-line no-console
        console.log("error decoding!");
      }
      return parsedToken.exp > Math.floor(Date.now() / 1000);
    },
    getBearerToken: state => state.token
  }
};
