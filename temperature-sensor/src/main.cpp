#include <Arduino.h>
#include <WiFi.h>
#include "PubSubClient.h"
#include <ArduinoJson.h>
#include <Wire.h>
#include <Adafruit_BMP085.h>
#include <Adafruit_BMP280.h>

// Replace the next variables with your SSID/Password combination
const char *ssid = "";
const char *password = "";

// Add your MQTT Broker IP address, example:
const char *mqtt_server = "192.168.2.135";

Adafruit_BMP280 bmp280; // I2C
Adafruit_BMP085 bmp180;

struct Button
{
  const uint8_t PIN;
  uint32_t numberKeyPresses;
  bool pressed;
};

Button button1 = {12, 0, false};
Button button2 = {13, 0, false};
Button button3 = {14, 0, false};

void IRAM_ATTR button1_pressed()
{
  button1.numberKeyPresses += 1;
  button1.pressed = true;
}
void IRAM_ATTR button2_pressed()
{
  button2.numberKeyPresses += 1;
  button2.pressed = true;
}
void IRAM_ATTR button3_pressed()
{
  button3.numberKeyPresses += 1;
  button3.pressed = true;
}

struct bmp180_sensor_data
{
  float altitude;
  float temperature;
  int32_t pressure;
  uint32_t raw_pressure;
  uint16_t raw_temperature;
  int32_t sea_level_pressure;

  void Print() { Serial.printf("bmp180_sensor_data-> altitude: %f, temperature: %f, pressure: %d, raw_pressure: %d, raw_temperature: %d, sea_level_pressure: %d\n",
                               altitude, temperature, pressure, raw_pressure, raw_temperature, sea_level_pressure); }
};

struct bmp280_sensor_data
{
  float temperature;
  float pressure;
  float altitude;

  void Print() { Serial.printf("bmp280_sensor_data-> temperature: %f, pressure: %f, altitude: %f\n",
                               temperature, pressure, altitude); }
};

TaskHandle_t collect_sensor_data_task_handle;
QueueHandle_t bmp180_sensor_data_queue, bmp280_sensor_data_queue;

void collect_sensor_data_task(void *pvParameters)
{
  bmp180_sensor_data bmp180_sensor_data;
  bmp280_sensor_data bmp280_sensor_data;
  for (;;)
  {

    bmp180_sensor_data.altitude = bmp180.readAltitude();
    bmp180_sensor_data.temperature = bmp180.readTemperature();
    bmp180_sensor_data.pressure = bmp180.readPressure();
    bmp180_sensor_data.raw_pressure = bmp180.readRawPressure();
    bmp180_sensor_data.raw_temperature = bmp180.readRawTemperature();
    bmp180_sensor_data.sea_level_pressure = bmp180.readSealevelPressure();

    bmp280_sensor_data.temperature = bmp280.readTemperature();
    bmp280_sensor_data.pressure = bmp280.readPressure();
    bmp280_sensor_data.altitude = bmp280.readAltitude();

    xQueueSend(bmp180_sensor_data_queue, &bmp180_sensor_data, portMAX_DELAY);
    xQueueSend(bmp280_sensor_data_queue, &bmp280_sensor_data, portMAX_DELAY);

    delay(250);
  }
  vTaskDelete(NULL);
}

WiFiClient wifi_client;
PubSubClient mqtt_client(wifi_client);

void setup_wifi()
{
  delay(10);
  // We start by connecting to a WiFi network
  Serial.println();
  Serial.print("Connecting to ");
  Serial.println(ssid);

  WiFi.begin(ssid, password);

  while (WiFi.status() != WL_CONNECTED)
  {
    delay(500);
    Serial.print(".");
  }

  Serial.println("");
  Serial.println("WiFi connected");
  Serial.println("IP address: ");
  Serial.println(WiFi.localIP());
}

void setup()
{
  Serial.begin(460800);
  pinMode(button1.PIN, INPUT_PULLUP);
  attachInterrupt(button1.PIN, button1_pressed, FALLING);
  pinMode(button2.PIN, INPUT_PULLUP);
  attachInterrupt(button2.PIN, button2_pressed, FALLING);
  pinMode(button3.PIN, INPUT_PULLUP);
  attachInterrupt(button3.PIN, button3_pressed, FALLING);
  if (!bmp280.begin(0x76))
  {
    Serial.println(F("Could not find a valid BMP280 sensor, check wiring!"));
    while (1)
      ;
  }

  if (!bmp180.begin())
  {
    Serial.println("Could not find a valid BMP085 sensor, check wiring!");
    while (1)
    {
    }
  }

  /* Default settings from datasheet. */
  bmp280.setSampling(Adafruit_BMP280::MODE_NORMAL,     /* Operating Mode. */
                     Adafruit_BMP280::SAMPLING_X2,     /* Temp. oversampling */
                     Adafruit_BMP280::SAMPLING_X16,    /* Pressure oversampling */
                     Adafruit_BMP280::FILTER_X16,      /* Filtering. */
                     Adafruit_BMP280::STANDBY_MS_500); /* Standby time. */

  setup_wifi();
  mqtt_client.setServer(mqtt_server, 1883);

  bmp180_sensor_data_queue = xQueueCreate(1, sizeof(bmp180_sensor_data));
  bmp280_sensor_data_queue = xQueueCreate(1, sizeof(bmp280_sensor_data));

  xTaskCreatePinnedToCore(
      collect_sensor_data_task,
      "collect_sensor_data_task",
      10000,
      NULL,
      0,
      &collect_sensor_data_task_handle,
      0);
  Serial.println("Task created...");
}
// JSON data buffer
StaticJsonDocument<1000> jsonDocument;
char buffer[1000];
void loop()
{
  bmp180_sensor_data bmp180_sensor_data;
  bmp280_sensor_data bmp280_sensor_data;

  while (!mqtt_client.connected())
  {
    mqtt_client.connect("ESP32 temperature sensor");
    delay(100);
  }
  if (xQueueReceive(bmp180_sensor_data_queue, &bmp180_sensor_data, 0))
  {
    bmp180_sensor_data.Print();

    jsonDocument.clear();
    jsonDocument["message_type"] = "BMP180_SENSOR_DATA";
    jsonDocument["data"] = jsonDocument.createNestedObject();
    jsonDocument["data"]["altitude"] = bmp180_sensor_data.altitude;
    jsonDocument["data"]["temperature"] = bmp180_sensor_data.temperature;
    jsonDocument["data"]["pressure"] = bmp180_sensor_data.pressure;
    jsonDocument["data"]["raw_pressure"] = bmp180_sensor_data.raw_pressure;
    jsonDocument["data"]["raw_temperature"] = bmp180_sensor_data.raw_temperature;
    jsonDocument["data"]["sea_level_pressure"] = bmp180_sensor_data.sea_level_pressure;

    serializeJson(jsonDocument, buffer);
    mqtt_client.publish("/esp32/toggle", buffer);
  }
  if (xQueueReceive(bmp280_sensor_data_queue, &bmp280_sensor_data, 0))
  {
    bmp280_sensor_data.Print();

    jsonDocument.clear();
    jsonDocument["message_type"] = "BMP280_SENSOR_DATA";
    jsonDocument["data"] = jsonDocument.createNestedObject();
    jsonDocument["data"]["temperature"] = bmp280_sensor_data.temperature;
    jsonDocument["data"]["pressure"] = bmp280_sensor_data.pressure;
    jsonDocument["data"]["altitude"] = bmp280_sensor_data.altitude;

    serializeJson(jsonDocument, buffer);
    mqtt_client.publish("/esp32/toggle", buffer);
  }

  if (button1.pressed)
  {
    Serial.printf("Button 1 has been pressed %u times\n", button1.numberKeyPresses);
    button1.pressed = false;
  }

  if (button2.pressed)
  {
    Serial.printf("Button 2 has been pressed %u times\n", button2.numberKeyPresses);
    button2.pressed = false;
  }

  if (button3.pressed)
  {
    Serial.printf("Button 3 has been pressed %u times\n", button3.numberKeyPresses);
    button3.pressed = false;
  }
}