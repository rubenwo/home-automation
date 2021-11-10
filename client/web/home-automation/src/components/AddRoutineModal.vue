<template>
    <b-modal
            size="xl"
            id="routineModal"
            ref="routineModal"
            @ok="handleOk"
            @cancel="handleCancel"
            @close="handleCancel"
    >
        <b-row>
            <div>
                <p>Routine:</p>
            </div>
        </b-row>
        <b-form-group>
            <b-input-group>
                <b-container fluid>
                    <b-row class="my-1">
                        <b-col sm="4">
                            <label>Name:</label>
                        </b-col>
                        <b-col sm="8">
                            <input
                                    size="sm"
                                    class="mx-1"
                                    placeholder="Routine Name"
                                    v-model="routine.name"
                            />
                        </b-col>
                        <b-col sm="4">
                            <label>Trigger Type:</label>
                        </b-col>
                        <b-col sm="8">
                            <b-form-select
                                    v-model="chosen_trigger_type"
                                    :options="trigger_type_options"
                            />
                        </b-col>
                        <div v-if="isTimerTriggerType()">
                            <b-input-group>
                                <b-container fluid>
                                    <b-row class="my-1">
                                        <b-col sm="4">
                                            <label>Schedule:</label>
                                        </b-col>
                                        <b-col>
                                            <input
                                                    size="sm"
                                                    class="mx-1"
                                                    placeholder="Cron expr"
                                                    v-model="routine.trigger.cron_expr"
                                            />
                                        </b-col>
                                    </b-row>
                                </b-container>
                            </b-input-group>
                        </div>
                        <div v-if="isEventTriggerType()">
                            <b-input-group>
                                <b-container fluid>
                                    <b-row class="my-1">
                                        <b-col sm="4">
                                            <label>Event:</label>
                                        </b-col>
                                        <b-col>
                                            <input
                                                    size="sm"
                                                    class="mx-1"
                                                    placeholder="ON_EVENT_NAME"
                                                    v-model="routine.trigger.on_event"
                                            />
                                        </b-col>
                                    </b-row>
                                </b-container>
                            </b-input-group>
                        </div>
                    </b-row>
                    <b-row>
                        <div>
                            <p>Actions:</p>
                        </div>
                    </b-row>
                    <b-row
                            v-bind:key="'action-' + index"
                            v-for="(action, index) in routine.actions"
                    >
                        <b-input-group>
                            <b-col sm="4">
                                <label>Device:</label>
                            </b-col>
                            <b-col sm="8">
                                <b-form-select
                                        v-model="action.deviceId"
                                        :options="computedDevices"
                                />
                                <div v-if="action.deviceId !== ''">
                                    Available methods:
                                    <div>
                                        <div v-bind:key="index + method.key"
                                             v-for="method in computeMethodsForDevice(action.deviceId)">
                                            <b-button variant="outline-success"
                                                      @click="selectMethodForAction(action.deviceId, method)">
                                                {{method.name}}
                                            </b-button>
                                        </div>
                                        <span v-if="action.name !== ''">Selected: <strong>{{action.name}}</strong></span>
                                        <span>{{action.deviceId}}</span>
                                    </div>
                                </div>
                            </b-col>
                        </b-input-group>
                        <b-button v-b-toggle.collapse-2 class="m-1">Advanced</b-button>
                        <b-collapse id="collapse-2">

                            <b-input-group>
                                <b-container fluid>
                                    <b-row class="my-1">
                                        <b-col sm="4">
                                            <label>Action Addr:</label>
                                        </b-col>
                                        <b-col sm="8">
                                            <b-form-input
                                                    size="sm"
                                                    class="mx-1"
                                                    placeholder="url"
                                                    v-model="action.addr"
                                            />
                                        </b-col>
                                    </b-row>
                                </b-container>
                            </b-input-group>
                            <b-input-group>
                                <b-container fluid>
                                    <b-row class="my-1">
                                        <b-col sm="4">
                                            <label>Action Method:</label>
                                        </b-col>
                                        <b-col sm="8">
                                            <b-form-input
                                                    size="sm"
                                                    class="mx-1"
                                                    placeholder="method"
                                                    v-model="action.method"
                                            />
                                        </b-col>
                                    </b-row>
                                </b-container>
                            </b-input-group>
                            <b-input-group>
                                <b-container fluid>
                                    <b-row class="my-1">
                                        <b-col sm="4">
                                            <label>Action Data:</label>
                                        </b-col>
                                        <b-col sm="8">
                                            <v-jsoneditor
                                                    v-model="action.data"
                                                    :options="jsonEditorOptions"
                                                    :plus="false"
                                                    :height="'300'"
                                                    @error="onError"
                                            />
                                        </b-col>
                                    </b-row>
                                </b-container>
                            </b-input-group>
                            <b-input-group>
                                <b-container fluid>
                                    <b-row class="my-1">
                                        <b-col sm="4">
                                            <label>Action Script:</label>
                                        </b-col>
                                        <b-col sm="8">
                                            <!--                                        <b-form-textarea-->
                                            <!--                                                v-model="action.script"-->
                                            <!--                                                placeholder="Enter some code in javascript..."-->
                                            <!--                                                rows="7"-->
                                            <!--                                        />-->
                                            <codemirror v-model="action.script" :options="cmOptions"/>
                                        </b-col>
                                    </b-row>
                                </b-container>
                            </b-input-group>
                        </b-collapse>
                    </b-row>

                    <b-button variant="success" @click="increaseActions">+</b-button>
                    <b-button variant="danger" @click="decreaseActions">-</b-button>
                </b-container>
            </b-input-group>
        </b-form-group>
    </b-modal>
</template>

<script>
  // import VueCronEditorBuefy from "vue-cron-editor-buefy";
  import VJsoneditor from "v-jsoneditor/src/index";
  import RoutineService from "../services/routines.service";
  import DevicesService from "../services/devices.service";
  import {codemirror} from "vue-codemirror";
  // import language js
  import "codemirror/mode/javascript/javascript.js";

  // import theme style
  import "codemirror/theme/base16-dark.css";
  // import base style
  import "codemirror/lib/codemirror.css";

  export default {
    name: "AddRoutineModal",
    components: {
      // VueCronEditorBuefy,
      VJsoneditor,
      codemirror,
    },
    data() {
      return {
        devices: [],
        cmOptions: {
          tabSize: 4,
          mode: "text/javascript",
          theme: "base16-dark",
          lineNumbers: true,
          line: true,
          // more CodeMirror options...
        },
        jsonEditorOptions: {
          mode: "code",
        },
        chosen_trigger_type: -1,
        trigger_type_options: [
          {value: 0, text: "TimerTriggerType"},
          {value: 1, text: "EventTriggerType"}
          // {value: 1, text: "WebhookTrigger"},
        ],
        routine: {
          name: "",
          trigger: {
            type: -1,
            cron_expr: "",
            on_event: ""
          },
          actions: [
            {
              name: "",
              deviceId: "",
              addr: "",
              method: "",
              data: null,
              script: "",
            },
          ],
        },
      };
    },
    methods: {
      onError() {
      },
      isTimerTriggerType() {
        return this.chosen_trigger_type === 0;
      },
      isEventTriggerType() {
        return this.chosen_trigger_type === 1;
      },
      increaseActions() {
        this.routine.actions.push({
          addr: "",
          method: "",
          data: null,
          script: "",
        });
      },
      decreaseActions() {
        this.routine.actions.pop();
      },
      handleOk() {
        this.handleSubmit();
      },
      selectMethodForAction(deviceId, method) {
        this.routine.actions.forEach(action => {
          if (action.deviceId === deviceId) {
            action.addr = method.addr;
            action.method = method.method;
            action.data = method.data;
            action.name = method.name;
          }
        })
      },
      computeMethodsForDevice(deviceId) {
        let methods = [];
        if (!deviceId) {
          return methods;
        }
        let device = this.devices.find(device => device.id === deviceId);
        if (!device) {
          return methods;
        }
        if (device.product.company === "IKEA") {
          if (device.product.type === "light") {
            // Turn off device
            methods.push({
              pressed: false,
              deviceId: deviceId,
              key: deviceId + "-turnoff",
              name: "Turn Off",
              addr: `http://tradfri.default.svc.cluster.local/tradfri/devices/${deviceId}/command`,
              method: "POST",
              data: {device_type: "light", dimmable_light_command: {power: 0}}
            });
            // Turn on device
            methods.push({
              pressed: false,
              deviceId: deviceId,
              key: deviceId + "-turnon",
              name: "Turn On",
              addr: `http://tradfri.default.svc.cluster.local/tradfri/devices/${deviceId}/command`,
              method: "POST",
              data: {device_type: "light", dimmable_light_command: {power: 1}}
            });
            // Dimm
            // Turn on device
            methods.push({
              pressed: false,
              deviceId: deviceId,
              key: deviceId + "-dimm",
              name: "Dimm",
              addr: `http://tradfri.default.svc.cluster.local/tradfri/devices/${deviceId}/command`,
              method: "POST",
              data: {device_type: "light", dimmable_light_command: {power: 1, brightness: 127}}
            });
          }
        } else if (device.product.company === "esp32") {
          if (device.category === "led-strip") {
            // Set color
            methods.push({
              pressed: false,
              deviceId: deviceId,
              key: deviceId + "-setcolor",
              name: "Set Color",
              addr: "",
              method: "",
              data: {},
              optional: {r: 0, g: 0, b: 0}
            });
            // Set colorcycle
            methods.push({
              pressed: false,
              deviceId: deviceId,
              key: deviceId + "-setcolorcycle",
              name: "Set Colorcycle",
              addr: "",
              method: "",
              data: {},
            });
            // Set breathing
            methods.push({
              pressed: false,
              deviceId: deviceId,
              key: deviceId + "-setbreathing",
              name: "Set Breathing",
              addr: "",
              method: "",
              data: {},
              optional: {r: 0, g: 0, b: 0}
            });
          }
        } else if (device.product.company === "tp-link") {
          if (device.category === "light" && device.product.type === "L510E") {
            // Turn off device
            methods.push({
              pressed: false,
              deviceId: deviceId,
              key: deviceId + "-turnoff",
              name: "Turn Off",
              addr: `http://tapo.default.svc.cluster.local/tapo/lights/${deviceId}?command=off&brightness=1`,
              method: "GET",
              data: null
            });
            // Turn on device
            methods.push({
              pressed: false,
              deviceId: deviceId,
              key: deviceId + "-turnon",
              name: "Turn On",
              addr: `http://tapo.default.svc.cluster.local/tapo/lights/${deviceId}?command=on&brightness=100`,
              method: "GET",
              data: null
            });
            // Dimm
            // Turn on device
            methods.push({
              pressed: false,
              deviceId: deviceId,
              key: deviceId + "-dimm",
              name: "Dimm",
              addr: `http://tapo.default.svc.cluster.local/tapo/lights/${deviceId}?command=on&brightness=50`,
              method: "GET",
              data: null,
              optional: {brightness: 0}
            });

          } else if (device.category === "plug" && device.product.type === "P100") {
            // Turn off device
            methods.push({
              pressed: false,
              deviceId: deviceId,
              key: deviceId + "-turnoff",
              name: "Turn Off",
              addr: `http://tapo.default.svc.cluster.local/tapo/lights/${deviceId}?command=off&brightness=1`,
              method: "GET",
              data: null
            });
            // Turn on device
            methods.push({
              pressed: false,
              deviceId: deviceId,
              key: deviceId + "-turnon",
              name: "Turn On",
              addr: `http://tapo.default.svc.cluster.local/tapo/lights/${deviceId}?command=on&brightness=100`,
              method: "GET",
              data: null
            });
          }
        }

        return methods;
      },
      handleCancel() {
      },
      async handleSubmit() {
        this.routine.trigger.type = this.chosen_trigger_type;
        this.chosen_trigger_type = -1;

        await RoutineService.addRoutine(this.routine);

        this.routine = {
          trigger: {
            type: -1,
            cron_expr: "",
            on_event: "",
          },
          actions: [
            {
              script: "",
              addr: "",
              method: "",
              data: null,
            },
          ],
        };

        this.$nextTick(() => {
          this.$refs.routineModal.hide();
        });
      },
    },
    computed: {
      computedDevices() {
        let types = [];
        this.devices.forEach(device => {
          types.push({value: device.id, text: `${device.product.company} - ${device.category} -  ${device.name}`})
        });
        return types;
      }
    },
    async created() {
      this.devices = (await DevicesService.fetchDevices()).devices;
      console.log(this.devices);

      this.$on("add_routine", () => {
        this.$refs.routineModal.show();
      });
    },
  };
</script>

<style scoped></style>
