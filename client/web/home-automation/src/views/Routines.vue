<template>
    <div>
        <br>
        <input class="form-control" type="text" placeholder="Search" aria-label="Search"/>
        <div class="row" style="margin-top: 10px">
            <b-col
                    cols="4"
                    sm="3"
                    md="3"
                    lg="2"
                    xl="2"
                    v-bind:key="routine.id"
                    v-for="routine in this.routines"
                    style="margin-right: 5em; margin-left: 5em"
            >
                <routine-card v-bind:id="routine.id" v-bind:name="routine.name"/>
            </b-col>
        </div>
    </div>
</template>

<script>
  import RoutineCard from "../components/RoutineCard";
  import RoutineService from "../services/routines.service";

  export default {
    name: "Routines",
    components: {RoutineCard},
    data() {
      return {
        routines: [],
      };
    },
    methods: {},
    async mounted() {
      this.routines = await RoutineService.fetchRoutines(this.id);
    },
  };
</script>

<style scoped></style>
