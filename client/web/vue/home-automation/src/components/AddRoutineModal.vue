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
            <b-col sm="4"><label>Schedule:</label></b-col>
            <b-col>
              <VueCronEditorBuefy v-model="cronExpression" />
              {{ cronExpression }}
            </b-col>
          </b-row>
          <b-row> </b-row>
        </b-container>
      </b-input-group>
    </b-form-group>
  </b-modal>
</template>

<script>
import VueCronEditorBuefy from "vue-cron-editor-buefy";

export default {
  name: "AddRoutineModal",
  components: {
    VueCronEditorBuefy,
  },
  data() {
    return {
      cronExpression: "*/1 * * * *",
      routine: {
        trigger: {
          type: "",
          repeat: true,
          when: "",
        },
        actions: [
          {
            addr: "",
            method: "",
            data: {},
          },
        ],
      },
    };
  },
  methods: {
    handleOk() {},
    handleCancel() {
      this.handleSubmit();
    },
    handleSubmit() {
      this.routine = {
        trigger: {
          type: "",
          repeat: true,
          when: "",
        },
        actions: [
          {
            addr: "",
            method: "",
            data: {},
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
