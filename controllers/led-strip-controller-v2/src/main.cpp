#include <Arduino.h>
#include <WiFi.h>
#include <ArduinoJson.h>
#include <Adafruit_NeoPixel.h>
#include <map>
#include <vector>
#include <functional>

#include "RGB12V.hpp"
#include "PubSubClient.h"
#include "strings.hpp"
#include "Modes.hpp"

#define ONBOARD_LED_PIN 2

#define LED_COUNT 150
#define LED_STRIP_DATA_PIN 4

#define RED_12V_PIN 12
#define GREEN_12V_PIN 13
#define BLUE_12V_PIN 14
#define RED_12V_PWM_CHANNEL 0
#define GREEN_12V_PWM_CHANNEL 1
#define BLUE_12V_PWM_CHANNEL 2

const char *ssid = "";
const char *password = "";
const char *mqtt_broker = "";
const char *device_name = "";
const char *device_id = "";

struct Color
{
  int R, G, B;
};

const Color default_color = {225, 125, 15};

std::map<std::string, std::function<void(byte *message, unsigned int length)>> topic_router;

WiFiClient wifi_client;
PubSubClient mqtt_client(wifi_client);

RGB12V *led_strip_12v;
Adafruit_NeoPixel led_strip_5v_addressable(LED_COUNT, LED_STRIP_DATA_PIN, NEO_GRB + NEO_KHZ800);

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

void callback(char *topic, byte *message, unsigned int length)
{

  digitalWrite(ONBOARD_LED_PIN, HIGH);
  Serial.println(length);

  auto func = topic_router[topic];
  if (func)
    func(message, length);

  digitalWrite(ONBOARD_LED_PIN, LOW);
}

void reconnect()
{
  // Loop until we're reconnected
  while (!mqtt_client.connected())
  {
    Serial.println("Attempting MQTT connection...");
    // Attempt to connect
    if (mqtt_client.connect(device_id))
    {
      Serial.println("connected");
      // Subscribe

      std::map<std::string, std::function<void(byte * message, unsigned int length)>>::iterator it;
      for (it = topic_router.begin(); it != topic_router.end(); it++)
      {
        auto key = it->first;
        mqtt_client.subscribe(key.c_str());

        Serial.println(key.c_str());
      }
    }
    else
    {
      Serial.print("failed, rc=");
      Serial.print(mqtt_client.state());
      Serial.println(" try again in 5 seconds");
      // Wait 5 seconds before retrying
      delay(5000);
    }
  }
}

void write_announcement_message()
{
  reconnect();

  DynamicJsonDocument jsonDocument(1024);
  char buffer[1024];

  jsonDocument.clear();

  jsonDocument["device_name"] = device_name;
  jsonDocument["device_id"] = device_id;
  jsonDocument["is_healthy"] = true;
  jsonDocument["error_message"] = "";

  serializeJson(jsonDocument, buffer);

  mqtt_client.publish("leds/announcement", buffer);
}

void get_information(byte *message, unsigned int length)
{
  DynamicJsonDocument jsonDocument(sizeof(byte) * length);
  deserializeJson(jsonDocument, message, length);

  DynamicJsonDocument responseDocument(1024);
  char buffer[1024];

  responseDocument.clear();

  responseDocument["device_name"] = device_name;
  responseDocument["device_id"] = device_id;
  responseDocument["is_healthy"] = true;
  responseDocument["error_message"] = "";

  serializeJson(responseDocument, buffer);

  mqtt_client.publish(string_format("leds/%s/response", device_id).c_str(), buffer);
}
int red, green, blue;
bool rgb_updated = false;
std::vector<Color> animation_colors;
unsigned int current_color_index;
Mode mode = Mode::SINGLE_COLOR_RGB;

int animation_speed = 10;
unsigned long timer;

void control_leds(byte *message, unsigned int length)
{
  DynamicJsonDocument jsonDocument(65535 * sizeof(byte));
  deserializeJson(jsonDocument, message, length);

  const std::string &json_mode = jsonDocument["mode"];
  mode = from_string(json_mode);

  if (mode == Mode::SINGLE_COLOR_RGB)
  {

    red = jsonDocument["r"];
    green = jsonDocument["g"];
    blue = jsonDocument["b"];

    rgb_updated = true;
  }
  else if (mode == Mode::GRADIENT_RGB)
  {
  }
  else if (mode == Mode::ANIMATION_RGB)
  {
    Serial.println("got an animation request");

    animation_speed = jsonDocument["animation_speed"];

    const JsonArray &json_config = jsonDocument["config"];
    animation_colors.clear();

    for (const JsonObject &value : json_config)
    {
      int r = value["r"];
      int g = value["g"];
      int b = value["b"];
      // Serial.printf("R: %d, G: %d, B: %d\n", r, g, b);
      animation_colors.push_back(Color{
          r,
          g,
          b,
      });
    }

    current_color_index = 0;
  }

  DynamicJsonDocument responseDocument(1024);
  char buffer[1024];

  responseDocument.clear();

  responseDocument["device_name"] = device_name;
  responseDocument["device_id"] = device_id;
  responseDocument["is_healthy"] = true;
  responseDocument["error_message"] = "";

  serializeJson(responseDocument, buffer);

  mqtt_client.publish(string_format("leds/%s/response", device_id).c_str(), buffer);
}

void setup()
{
  Serial.begin(460800);
  delay(1000);
  pinMode(ONBOARD_LED_PIN, OUTPUT);
  // Turn on onboard LED to indicate the setup
  digitalWrite(ONBOARD_LED_PIN, HIGH);

  led_strip_12v = new RGB12V(12, 13, 14, 0, 1, 2);

  led_strip_5v_addressable.begin();
  led_strip_5v_addressable.fill(led_strip_5v_addressable.gamma32(led_strip_5v_addressable.Color(default_color.R, default_color.G, default_color.B)), 0, led_strip_5v_addressable.numPixels() - 1);
  led_strip_5v_addressable.show();

  led_strip_12v->set_color(default_color.R, default_color.G, default_color.B);
  led_strip_12v->show();

  topic_router[string_format("leds/%s/information", device_id)] = get_information;
  topic_router[string_format("leds/%s/control", device_id)] = control_leds;

  setup_wifi();
  mqtt_client.setServer(mqtt_broker, 1883);
  mqtt_client.setCallback(callback);

  write_announcement_message();

  digitalWrite(ONBOARD_LED_PIN, LOW);

  // put your setup code here, to run once:
}

void loop()
{
  if (!mqtt_client.connected())
  {
    reconnect();
  }
  mqtt_client.loop();

  if (mode == Mode::SINGLE_COLOR_RGB && rgb_updated)
  {
    led_strip_5v_addressable.fill(led_strip_5v_addressable.gamma32(led_strip_5v_addressable.Color(red, green, blue)), 0, led_strip_5v_addressable.numPixels() - 1);
    led_strip_5v_addressable.show();

    led_strip_12v->set_color(red, green, blue);
    led_strip_12v->show();
  }

  if (mode == Mode::ANIMATION_RGB && animation_colors.size() > 0)
  {
    if (millis() - timer > animation_speed)
    {
      Color current_color = animation_colors[current_color_index % animation_colors.size()];

      led_strip_5v_addressable.fill(led_strip_5v_addressable.gamma32(led_strip_5v_addressable.Color(current_color.R, current_color.G, current_color.B)), 0, led_strip_5v_addressable.numPixels());
      led_strip_5v_addressable.show();

      led_strip_12v->set_color(current_color.R, current_color.G, current_color.B);
      led_strip_12v->show();

      current_color_index++;
      timer = millis();
    }
  }
}