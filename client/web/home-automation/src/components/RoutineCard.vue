<template>
    <b-card
            v-bind:sub-title="name"
            style="
      max-width: 300px;
      min-width: 300px;
      min-height: 250px;
      max-height: 250px;
        background-color: rgba(255, 255, 255, 0.25);
        backdrop-filter: blur(5px);    "
            class="mb-2"
    >
        <div>Type: {{ routine.trigger.type }}</div>
        <div>Schedule: {{ routine.trigger.cron_expr }}</div>
        <div># of actions: {{ routine.actions.length }}</div>

        <b-checkbox
                v-model="routine.is_active"
                @change="set_active()"
        >: Active
        </b-checkbox>

        <div slot="footer">
            <b-button style="background-color: #4287f5" v-bind:to="navigate()"
            >Information
            </b-button>
            <b-button variant="danger" @click="deleteRoutine()">X</b-button>
        </div>
    </b-card>
</template>

<script>
  import RoutineService from "../services/routines.service";

  export default {
    name: "RoutineCard",
    props: {
      id: {
        type: Number,
        default: -1,
      },
      name: {
        type: String,
        default: "No name provided",
      },
    },
    data() {
      return {
        routine: {
          name: "",
          is_active: true,
          trigger: {
            type: -1,
            cron_expr: "",
          },
          actions: [],
        },
      };
    },
    methods: {
      navigate() {
        return "routine/" + this.id;
      },
      async deleteRoutine() {
        const res = await RoutineService.deleteRoutine(this.id);
        console.log(res);
      },
      async set_active() {
        console.log("clicked");
        // this.routine.is_active = !this.routine.is_active;
        this.updateRoutine();
      },
      async updateRoutine() {
        const res = await RoutineService.updateRoutine(this.id, this.routine);
        console.log(res);
        this.routine = res.routine;
      },
    },
    async mounted() {
      const routine = await RoutineService.fetchRoutine(this.id);
      console.log(routine);
      this.routine = routine.routine;
    },
  };
</script>

<style scoped></style>
