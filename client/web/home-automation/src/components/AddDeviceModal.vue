<template>
    <b-modal
            size="xl"
            id="deviceModal"
            ref="deviceModal"
            @ok="handleOk"
            @cancel="handleCancel"
            @close="handleCancel"
    >
        <b-form-group>
            <b-input-group>
                <b-container fluid>
                    <b-row class="my-1">
                        <b-col sm="4">
                            <label>Device Type:</label>
                        </b-col>
                        <b-col sm="8">
                            <b-form-select
                                    v-model="device_type"
                                    :options="device_type_options"
                            ></b-form-select>
                        </b-col>
                    </b-row>
                </b-container>
            </b-input-group>
            <div v-if="isTapoDevice()">
                <b-input-group>
                    <b-container fluid>
                        <b-row class="my-1">
                            <b-col sm="4">
                                <label>IP Address:</label>
                            </b-col>
                            <b-col sm="8">
                                <input
                                        size="sm"
                                        class="mx-1"
                                        placeholder="ip address"
                                        v-model="newItem.ip"
                                />
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
                                <input
                                        size="sm"
                                        class="mx-1"
                                        placeholder="email"
                                        v-model="newItem.email"
                                />
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
                                <input
                                        type="password"
                                        size="sm"
                                        class="mx-1"
                                        placeholder="password"
                                        v-model="newItem.password"
                                />
                            </b-col>
                        </b-row>
                    </b-container>
                </b-input-group>
            </div>
            <div v-if="isLEDStripDevice()">
                <b-input-group>
                    <b-container fluid>
                        <b-row class="my-1">
                            <b-col sm="4">
                                <label>IP Address:</label>
                            </b-col>
                            <b-col sm="8">
                                <input
                                        size="sm"
                                        class="mx-1"
                                        placeholder="ip address"
                                        v-model="newItem.ip"
                                />
                            </b-col>
                        </b-row>
                    </b-container>
                </b-input-group>
            </div>
        </b-form-group>
    </b-modal>
</template>

<script>
  import {mapActions, mapState} from "vuex";

  export default {
    name: "AddDeviceModal",
    data: () => ({
      device_type: "",
      newItem: {
        ip: "",
        email: "",
        password: "",
        device_type: "",
      },
      device_type_options: [
        {value: "L510E", text: "L510E"},
        {value: "P100", text: "P100"},
        {value: "LED_STRIP", text: "LED Strip"},
      ],
    }),
    computed: {
      ...mapState("devices", {
        devices: (state) => state.devices,
      }),
    },
    methods: {
      ...mapActions("devices", ["addNewDevice"]),
      isTapoDevice() {
        return this.device_type === "L510E" || this.device_type === "P100";
      },
      isLEDStripDevice() {
        return this.device_type === "LED_STRIP";
      },
      handleOk(evt) {
        evt.preventDefault();
        console.log(this.newItem);
        this.handleSubmit();
      },
      handleCancel() {
      },
      async handleSubmit() {
        let input = {};
        if (this.isTapoDevice()) {
          input = {
            device_type: "tapo",
            data: {
              ip_address: this.newItem.ip,
              email: this.newItem.email,
              password: this.newItem.password,
              device_type: this.device_type,
            },
          };
        } else if (this.isLEDStripDevice()) {
          input = {
            device_type: "LED_STRIP",
            data: {
              ip_address: this.newItem.ip,
            },
          };
        }
        await this.addNewDevice(input);
        this.device_type = "";
        this.newItem.ip = "";
        this.newItem.email = "";
        this.newItem.password = "";
        this.newItem.device_type = "";
        this.$nextTick(() => {
          this.$refs.deviceModal.hide();
        });
      },
    },
    created() {
      this.$on("add_device", () => {
        console.log("Got Event in Device");
        this.$refs.deviceModal.show();
      });
    },
  };
</script>

<style scoped></style>
