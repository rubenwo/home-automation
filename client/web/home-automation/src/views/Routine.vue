<template>
    <div class="routine">
        <h3>{{ routine.name }}</h3>
        <p>Is active: {{routine.is_active ? "True":"False"}}</p>
        <p>Trigger: type: {{routine.trigger.type}}, cron_expr: {{routine.trigger.cron_expr}}</p>
        <ul>
            <li v-bind:key="'action-' + index" v-for="(action, index) in routine.actions">
                [{{ action.method }}] - [{{action.addr}}]
                <br>
                Script: [{{action.script}}]
            </li>
        </ul>
    </div>
</template>

<script>
  import RoutinesService from "../services/routines.service";

  export default {
    name: "Routine",
    data() {
      return {
        routine: {
          id: 0,
          name: "",
          is_active: false,
          trigger: {
            cron_expr: "",
            type: 0
          },
          actions: [
            {
              script: "",
              data: null,
              method: "",
              addr: ""
            }
          ]
        }
      }
    },
    async created() {
      this.id = this.$route.params.id;

      const res = await RoutinesService.fetchRoutine(this.id);
      console.log(res);
      this.routine = res.routine;
      console.log(this.routine);
    },
    mounted() {
    },
  }
</script>

<style scoped>
    .routine {
        margin-top: 10px;
        text-align: center;
        background-color: rgba(255, 255, 255, 0.7);
    }
</style>
