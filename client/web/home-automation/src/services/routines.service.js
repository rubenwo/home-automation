import ApiService from "./api.service";

export default {
    async fetchRoutines() {
        const res = await ApiService()
            .get("/api/v1/routines")
            .catch(() => {
                return null;
            });
        console.log(res);
        return res.data;
    },
    async addRoutine(data) {
        const res = await ApiService()
            .post("/api/v1/routines", data)
            .catch(() => {
                return null;
            });
        console.log(res);
        return res.data;
    },
    async updateRoutine(id, data) {
        const res = await ApiService()
            .put("/api/v1/routines/" + id, data)
            .catch(() => {
                return null;
            });
        console.log(res);
        return res.data;
    },
    async fetchRoutine(routineId) {
        const res = await ApiService()
            .get("/api/v1/routines/" + routineId)
            .catch(() => {
                return null;
            });
        console.log(res);
        return res.data;
    },
    async deleteRoutine(recipeId) {
        const res = await ApiService()
            .delete("/api/v1/routines/" + recipeId)
            .catch(() => {
                return null;
            });
        console.log(res);
        return res.data;
    },
    async fetchLogs() {
        const res = await ApiService()
            .get("/api/v1/routines/logs")
            .catch(() => {
                return null;
            });
        console.log(res);
        return res.data;
    },
    async fetchLogsForId(routineId) {
        const res = await ApiService()
            .get("/api/v1/routines/logs/" + routineId)
            .catch(() => {
                return null;
            });
        console.log(res);
        return res.data;
    }
};
