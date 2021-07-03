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

        <div>
            <h4>Logs</h4>
            <ul v-if="logs !== null">
                <li v-bind:key="'logs-' +index" v-for="(log, index) in logs">
                    [{{log.logged_at}}] - [{{log.message}}]
                </li>
            </ul>
        </div>
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
        },
        logs: [
          {
            logged_at: "",
            message: ""
          }
        ]
      }
    },
    async created() {
      this.id = this.$route.params.id;

      const res = await RoutinesService.fetchRoutine(this.id);
      console.log(res);
      this.routine = res.routine;
      console.log(this.routine);

      const logsRes = await RoutinesService.fetchLogsForId(this.id);
      console.log(logsRes);
      this.logs = logsRes.logs;
      console.log(this.logs);
    },
    mounted() {
    },
  }
</script>

<style scoped>
    .routine {
        margin-top: 10px;
        text-align: center;
        background-color: rgba(255, 255, 255, 0.25);
        backdrop-filter: blur(5px);
    }
</style>
