import AuthService from "../../services/auth.service";

export default {
  namespaced: true,
  modules: {},
  state: {
    authorization_token: localStorage.getItem("authorization_token") || null,
    refresh_token: localStorage.getItem("refresh_token") || null,
    username: localStorage.getItem("username") || null,
    userid: localStorage.getItem("userid") || null,
    user: null,
    error: null,
  },
  mutations: {
    SET_AUTHORIZATION_TOKEN(state, authorization_token) {
      state.authorization_token = authorization_token;
    },
    SET_REFRESH_TOKEN(state, refresh_token) {
      state.refresh_token = refresh_token;
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
      state.authorization_token = null;
      state.refresh_token = null;
      state.id = null;
      state.user = null;
      state.username = null;
    },
  },
  actions: {
    async login({ commit }, { username, password }) {
      commit("CLEAR_ERROR");
      console.log("logging in");
      const resp = await AuthService.login(username, password);
      console.log(resp);
      if (resp.status === 200) {
        const {
          username,
          user_id,
          authorization_token,
          refresh_token,
        } = resp.data;
        localStorage.setItem("authorization_token", authorization_token);
        localStorage.setItem("refresh_token", refresh_token);
        localStorage.setItem("username", username);
        localStorage.setItem("userid", user_id);

        commit("SET_USERNAME", username);
        commit("SET_AUTHORIZATION_TOKEN", authorization_token);
        commit("SET_REFRESH_TOKEN", refresh_token);
        commit("SET_ID", user_id);
        return true;
      }
      return false;
    },
    async logout({ commit }) {
      localStorage.removeItem("authorization_token");
      localStorage.removeItem("refresh_token");
      localStorage.removeItem("username");
      localStorage.removeItem("userid");
      commit("CLEAR_ALL");
      commit("CLEAR_ERROR");
    },
  },
  getters: {
    isLoggedIn: (state) => {
      if (state.authorization_token == null) return false;
      const parseJwt = (token) => {
        try {
          return JSON.parse(atob(token.split(".")[1]));
        } catch (e) {
          return null;
        }
      };
      let parsedToken = parseJwt(state.authorization_token);
      if (parsedToken == null) {
        // eslint-disable-next-line no-console
        console.log("error decoding!");
      }
      return parsedToken.exp > Math.floor(Date.now() / 1000);
    },
    getBearerToken: (state) => state.authorization_token,
  },
};
