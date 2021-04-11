<template>
  <div>
    <p style="color: aliceblue">BMP180 Temperature: {{ bmp180_temperature }}</p>
    <p style="color: aliceblue">BMP280 Temperature: {{ bmp280_temperature }}</p>
    <p style="color: aliceblue"># of Satellites: {{ number_of_satellites }}</p>
  </div>
</template>

<script>
export default {
  name: "Recipes",
  components: {},
  data() {
    return {
      recipes: [],
      connection: null,
      bmp180_temperature: 0.0,
      bmp280_temperature: 0.0,
      number_of_satellites: 0
    };
  },
  methods: {
    onMessage(event) {
      let data = JSON.parse(event.data);
      console.log(data);
      switch (data.message_type) {
        case "BMP180_SENSOR_DATA": {
          this.bmp180_temperature = data.data.temperature;
          break;
        }
        case "BMP280_SENSOR_DATA": {
          this.bmp280_temperature = data.data.temperature;
          break;
        }
        case "GPS_DATA": {
          this.number_of_satellites = data.data.gps_sat_amount;
          break;
        }
      }
    }
  },
  async mounted() {},
  created() {
    console.log("Starting connection to WebSocket Server");
    console.log(process.env.VUE_APP_BACKEND_URL);

    let url = "wss://";
    if (process.env.VUE_APP_BACKEND_URL === "/") url += window.location.host;
    else url += process.env.VUE_APP_BACKEND_URL;
    url += "/api/v1/esp32/toggle";

    console.log(url);

    this.connection = new WebSocket(url);
    this.connection.onmessage = this.onMessage;

    this.connection.onopen = function(event) {
      console.log(event);
      console.log("Successfully connected to the echo websocket server...");
    };
  },
  destroyed() {
    console.log("closed websocket");
    this.connection.close();
  }
};
</script>

<style></style>
