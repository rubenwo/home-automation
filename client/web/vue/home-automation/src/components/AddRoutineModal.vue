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
                    <codemirror v-model="action.script" :options="cmOptions" />
                  </b-col>
                </b-row>
              </b-container>
            </b-input-group>
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
import { codemirror } from "vue-codemirror";
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
        { value: 0, text: "TimerTriggerType" },
        // {value: 1, text: "WebhookTrigger"},
      ],
      routine: {
        name: "",
        trigger: {
          type: -1,
          cron_expr: "",
        },
        actions: [
          {
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
      console.log("error");
    },
    isTimerTriggerType() {
      return this.chosen_trigger_type === 0;
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
    handleCancel() {},
    async handleSubmit() {
      this.routine.trigger.type = this.chosen_trigger_type;
      this.chosen_trigger_type = -1;

      console.log(this.routine);

      const res = await RoutineService.addRoutine(this.routine);
      console.log(res);

      this.routine = {
        trigger: {
          type: -1,
          cron_expr: "",
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
  created() {
    this.$on("add_routine", () => {
      console.log("Got Event in Routine");
      this.$refs.routineModal.show();
    });
  },
};
</script>

<style scoped></style>
