<template>
  <b-modal
    size="xl"
    id="sensorModal"
    ref="sensorModal"
    @ok="handleOk"
    @cancel="handleCancel"
    @close="handleCancel"
  >
    <p>Hello Sensor modal</p>
  </b-modal>
</template>

<script>
export default {
  name: "AddSensorModal",
  methods: {
    handleOk() {},
    handleCancel() {},
    handleSubmit() {
      this.$nextTick(() => {
        this.$refs.sensorModal.hide();
      });
    },
  },
  created() {
    this.$on("add_sensor", () => {
      this.$refs.sensorModal.show();
    });
  },
};
</script>

<style scoped></style>
