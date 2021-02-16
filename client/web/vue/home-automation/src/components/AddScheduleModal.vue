<template>
    <b-modal
            size="xl"
            id="scheduleModal"
            ref="scheduleModal"
            @ok="handleOk"
            @cancel="handleCancel"
            @close="handleCancel"
    >
        <b-row>
            <div>
                <p>Schedule:</p>
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
                                    placeholder="Schedule Name"
                                    v-model="schedule.name"
                            />
                        </b-col>
                    </b-row>
                </b-container>
            </b-input-group>
        </b-form-group>
    </b-modal>
</template>

<script>
  export default {
    name: "AddScheduleModal",
    data() {
      return {
        schedule: {
          name: ""
        }
      }
    },
    methods: {
      handleOk() {
      },
      handleCancel() {
        this.handleSubmit()
      },
      handleSubmit() {
        this.$nextTick(() => {
          this.$refs.scheduleModal.hide();
        });
      }
    },
    created() {
      this.$on("add_schedule", () => {
        console.log("Got Event in Schedule");
        this.$refs.scheduleModal.show();
      });
    }
  }
</script>

<style scoped>

</style>
