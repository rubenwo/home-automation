<template>
    <div class="login" v-if="!this.authenticating">
        <div>
            <b-alert v-model="showAuthError" variant="danger" dismissible
            >{{ this.context.msg }}
            </b-alert>

            <b-card align="center" title="Login" style="background-color: rgba(255, 255, 255, 0.25); margin-top:10px;">
                <b-alert :show="!!error" variant="danger">{{ error }}</b-alert>
                <b-form @submit.prevent="onSubmit">
                    <b-form-group
                            id="usernameInput"
                            label="Username:"
                            label-for="exampleInput1"
                    >
                        <b-form-input
                                id="username"
                                type="text"
                                v-model="form.username"
                                required
                                placeholder="Enter username"
                        />
                    </b-form-group>

                    <b-form-group
                            id="passwordInput"
                            label="Password:"
                            label-for="exampleInput2"
                    >
                        <b-form-input
                                id="password"
                                type="password"
                                v-model="form.password"
                                required
                                placeholder="Enter password"
                        />
                    </b-form-group>
                    <b-button type="submit" variant="primary">Login</b-button>
                </b-form>
            </b-card>
        </div>
    </div>
    <div v-else>
        <loading :active="authenticating" :is-full-page="true"/>
    </div>
</template>

<script>
  import {mapActions, mapGetters, mapState} from "vuex";

  import Loading from "vue-loading-overlay";
  import "vue-loading-overlay/dist/vue-loading.css";

  export default {
    name: "Login",
    components: {
      Loading,
    },
    data() {
      return {
        form: {
          username: "",
          password: "",
        },
        showRegisterSuccessAlert: false,
        showAuthError: false,
        context: {msg: ""},
        authenticating: false,
      };
    },
    computed: {
      ...mapState("auth", ["error"]),
    },
    methods: {
      ...mapActions("auth", ["login", "logout"]),
      ...mapGetters("auth", ["isLoggedIn"]),
      async onSubmit() {
        this.authenticating = true;
        const success = await this.login(this.form);
        this.authenticating = false;
        if (!success) {
          this.showAuthError = true;
          this.context.msg = "username/password is incorrect";
          return;
        }
        await this.$router.push(this.$route.query.redirect || "/");
      },
    },
  };
</script>

<style scoped></style>
