#include <Arduino.h>
#include <WiFi.h>
#include "PubSubClient.h"
#include <ArduinoJson.h>

RTC_DATA_ATTR int bootCount = 0;

// Replace the next variables with your SSID/Password combination
const char *ssid = "";
const char *password = "";

// Add your MQTT Broker IP address, example:
const char *mqtt_server = "192.168.2.135";

WiFiClient wifi_client;
PubSubClient mqtt_client(wifi_client);

void print_wakeup_reason();
void write_to_mqtt();

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

#define LED 2

void setup()
{
  Serial.begin(460800);
  delay(1000);

  pinMode(LED, OUTPUT);

  digitalWrite(LED, HIGH);

  //Increment boot number and print it every reboot
  ++bootCount;
  Serial.println("Boot number: " + String(bootCount));

  //Print the wakeup reason for ESP32
  print_wakeup_reason();

  setup_wifi();
  mqtt_client.setServer(mqtt_server, 1883);

  write_to_mqtt();

  //Configure GPIO33 as ext0 wake up source for HIGH logic level
  esp_sleep_enable_ext0_wakeup(GPIO_NUM_33, 1);

  delay(4000);

  digitalWrite(LED, LOW);

  //Go to sleep now
  esp_deep_sleep_start();
}

void loop()
{
  // int signal = digitalRead(GPIO_NUM_33);
  // Serial.printf("Input: %d\n", signal);
  // delay(2000);
}

// JSON data buffer
StaticJsonDocument<1000> jsonDocument;
char buffer[1000];

void write_to_mqtt()
{
  while (!mqtt_client.connected())
  {
    mqtt_client.connect("ESP32 temperature sensor");
    delay(100);
  }

  jsonDocument.clear();
  jsonDocument["message_type"] = "MOTION_DETECTED";
  jsonDocument["data"] = jsonDocument.createNestedObject();
  jsonDocument["data"]["msg"] = "Detected motion";
  jsonDocument["data"]["ip_addr"] = WiFi.localIP();

  serializeJson(jsonDocument, buffer);
  mqtt_client.publish("/sensors/motion", buffer);
}

//Function that prints the reason by which ESP32 has been awaken from sleep
void print_wakeup_reason()
{
  esp_sleep_wakeup_cause_t wakeup_reason;
  wakeup_reason = esp_sleep_get_wakeup_cause();
  switch (wakeup_reason)
  {
  case 1:
    Serial.println("Wakeup caused by external signal using RTC_IO");
    break;
  case 2:
    Serial.println("Wakeup caused by external signal using RTC_CNTL");
    break;
  case 3:
    Serial.println("Wakeup caused by timer");
    break;
  case 4:
    Serial.println("Wakeup caused by touchpad");
    break;
  case 5:
    Serial.println("Wakeup caused by ULP program");
    break;
  default:
    Serial.println("Wakeup was not caused by deep sleep");
    break;
  }
}