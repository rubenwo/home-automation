<template>
    <b-modal size="xl" id="deviceModal" ref="modal" @ok="handleOk" @cancel="handleCancel" @close="handleCancel">
        <b-form-group>
            <b-input-group>
                <b-container fluid>
                    <b-row class="my-1">
                        <b-col sm="4">
                            <label>IP Address:</label>
                        </b-col>
                        <b-col sm="8">
                            <input size="sm" class="mx-1" placeholder="ip address"
                                   v-model="newItem.ip"/>
                        </b-col>
                    </b-row>
                </b-container>
            </b-input-group>
            <b-input-group>
                <b-container fluid>
                    <b-row class="my-1">
                        <b-col sm="4">
                            <label>Email:</label>
                        </b-col>
                        <b-col sm="8">
                            <input size="sm" class="mx-1" placeholder="email"
                                   v-model="newItem.email"/>
                        </b-col>
                    </b-row>
                </b-container>
            </b-input-group>
            <b-input-group>
                <b-container fluid>
                    <b-row class="my-1">
                        <b-col sm="4">
                            <label>Password:</label>
                        </b-col>
                        <b-col sm="8">
                            <input type="password" size="sm" class="mx-1" placeholder="password"
                                   v-model="newItem.password"/>
                        </b-col>
                    </b-row>
                </b-container>
            </b-input-group>
            <b-input-group>
                <b-container fluid>
                    <b-row class="my-1">
                        <b-col sm="4">
                            <label>Device Type:</label>
                        </b-col>
                        <b-col sm="8">
                            <b-form-select v-model="newItem.device_type" :options="device_type_options"></b-form-select>
                        </b-col>
                    </b-row>
                </b-container>
            </b-input-group>
        </b-form-group>
    </b-modal>

</template>

<script>
  import {mapActions, mapState} from "vuex";

  export default {
    name: "AddDeviceModal",
    data: () => ({
      newItem: {
        ip: "",
        email: "",
        password: "",
        device_type: "",
      },
      device_type_options: [
        {value: "L510E", text: "L510E"},
        {value: "P100", text: "P100"}
      ]
    }),
    computed: {
      ...mapState("devices", {
        devices: state => state.devices
      })
    },
    methods: {
      ...mapActions("devices", ["addNewDevice"]),

      handleOk(evt) {
        evt.preventDefault();
        console.log(this.newItem)
        this.handleSubmit();
      },
      handleCancel() {
      },
      async handleSubmit() {
        let input = {
          device_type: "tapo",
          data: {
            "ip_address": this.newItem.ip,
            "email": this.newItem.email,
            "password": this.newItem.password,
            "device_type": this.newItem.device_type
          }
        }
        await this.addNewDevice(input);

        this.newItem.ip = "";
        this.newItem.email = "";
        this.newItem.password = "";
        this.newItem.device_type = "";
        this.$nextTick(() => {
          this.$refs.modal.hide();
        });
      }
    },
    created() {
      this.$on('add_device', () => {

        this.$refs.modal.show()
      })
    },
  }
</script>

<style scoped>

</style>
