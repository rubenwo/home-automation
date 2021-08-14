#include <Arduino.h>
#include <WiFi.h>
#include <HTTPClient.h>
#include <WebServer.h>
#include <ArduinoJson.h>
#include <Adafruit_NeoPixel.h>
#include <array>
#include <vector>
#include "color.hpp"
#include "RGB12V.hpp"

#define HTTP_PORT 80
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

String device_name = "rgb esp32 led-strip PCB";

WebServer server(HTTP_PORT);

RGB12V *led_strip_12v;

Adafruit_NeoPixel led_strip(LED_COUNT, LED_STRIP_DATA_PIN, NEO_GRB + NEO_KHZ800);

int red, green, blue;
int hue, value, saturation;
std::vector<rgb> gradient_rgb_config;
std::vector<hsv> gradient_hsv_config;
enum Mode
{
  INVALID = -1,
  SINGLE_COLOR_RGB = 0,
  SINGLE_COLOR_HSV = 1,
  GRADIENT_RGB = 2,
  GRADIENT_HSV = 3,
  ANIMATION_RGB = 4,
  ANIMATION_HSV = 5,
};

Mode mode = Mode::SINGLE_COLOR_RGB;

Mode from_string(const String &str)
{
  if (str == "SINGLE_COLOR_RGB")
    return Mode::SINGLE_COLOR_RGB;
  else if (str == "SINGLE_COLOR_HSV")
    return Mode::SINGLE_COLOR_HSV;
  else if (str == "GRADIENT_RGB")
    return Mode::GRADIENT_RGB;
  else if (str == "GRADIENT_HSV")
    return Mode::GRADIENT_HSV;
  else if (str == "ANIMATION_RGB")
    return Mode::ANIMATION_RGB;
  else if (str == "ANIMATION_HSV")
    return Mode::ANIMATION_HSV;

  return Mode::INVALID;
}

String mode_to_string(Mode m)
{
  switch (m)
  {
  case Mode::SINGLE_COLOR_RGB:
    return "SINGLE_COLOR_RGB";
  case Mode::SINGLE_COLOR_HSV:
    return "SINGLE_COLOR_HSV";
  case Mode::GRADIENT_RGB:
    return "GRADIENT_RGB";
  case Mode::GRADIENT_HSV:
    return "GRADIENT_HSV";
  case Mode::ANIMATION_RGB:
    return "ANIMATION_RGB";
  case Mode::ANIMATION_HSV:
    return "ANIMATION_HSV";
  case Mode::INVALID:
    return "INVALID";
  }
  return "INVALID";
}

unsigned long timer;
int animation_speed = 10;

// JSON data buffer
StaticJsonDocument<1000> jsonDocument;
char buffer[1000];

void setup()
{
  Serial.begin(115200);
  Serial.printf("Connecting to: %s\n", ssid);

  led_strip_12v = new RGB12V(12, 13, 14, 0, 1, 2);

  // setup wifi
  WiFi.begin(ssid, password);

  while (WiFi.status() != WL_CONNECTED)
  {
    Serial.print(".");
    delay(500);
  }
  Serial.printf(" Connected. IP: %s\n", WiFi.localIP().toString().c_str());

  // setup routing web server
  server.on("/healthz", HTTP_GET, []()
            {
              Serial.println("/healthz");
              jsonDocument.clear();
              jsonDocument["is_healthy"] = true;
              jsonDocument["error_message"] = "";
              serializeJson(jsonDocument, buffer);
              server.send(200, "application/json", buffer);
            });
  server.on("/led", HTTP_POST, []()
            {
              Serial.println("/led");
              if (server.hasArg("plain") == false)
              {
                Serial.println("ERROR in /led");
                jsonDocument.clear();
                jsonDocument["error_message"] = "body is not valid";
                serializeJson(jsonDocument, buffer);
                server.send(422, "application/json", buffer);
                return;
              }

              String body = server.arg("plain");
              deserializeJson(jsonDocument, body);

              // Get RGB components
              String json_mode = jsonDocument["mode"];
              mode = from_string(json_mode);
              JsonArray json_config;

              switch (mode)
              {
              case Mode::SINGLE_COLOR_RGB:
                red = jsonDocument["red"];
                green = jsonDocument["green"];
                blue = jsonDocument["blue"];

                led_strip.fill(led_strip.gamma32(led_strip.Color(red, green, blue)), 0, led_strip.numPixels() - 1);
                led_strip.show();

                led_strip_12v->set_color(red, green, blue);
                led_strip_12v->show();

                break;
              case Mode::SINGLE_COLOR_HSV:
                hue = jsonDocument["hue"];
                saturation = jsonDocument["saturation"];
                value = jsonDocument["value"];

                led_strip.fill(led_strip.gamma32(led_strip.ColorHSV(hue, saturation, value)), 0, led_strip.numPixels() - 1);
                led_strip.show();
                break;
              case Mode::GRADIENT_RGB:
              {
                json_config = jsonDocument["config"];
                gradient_rgb_config.clear();
                for (auto i = 0; i < json_config.size(); i++)
                {
                  gradient_rgb_config.push_back(rgb{
                      json_config[i]["red"],
                      json_config[i]["green"],
                      json_config[i]["blue"],
                  });
                }
                int step_size = LED_COUNT / gradient_rgb_config.size();
                int rgb_index = 0;
                for (auto i = 0; i < led_strip.numPixels(); i++)
                {
                  if (i / (rgb_index + 1) >= step_size)
                    rgb_index++;
                  uint32_t col = led_strip.gamma32(led_strip.Color(
                      gradient_rgb_config[rgb_index].r,
                      gradient_rgb_config[rgb_index].g,
                      gradient_rgb_config[rgb_index].b));
                  led_strip.setPixelColor(i, col);

                  // Serial.printf("Index: %d, RGB_INDEX: %d, Step Size: %d, R: %d, G: %d, B: %d\n", i, rgb_index, step_size, gradient_rgb_config[rgb_index].r, gradient_rgb_config[rgb_index].g, gradient_rgb_config[rgb_index].b);
                }
                led_strip.show();
              }
              break;
              case Mode::GRADIENT_HSV:
              {
                json_config = jsonDocument["config"];
                gradient_hsv_config.clear();
                for (auto i = 0; i < json_config.size(); i++)
                {
                  gradient_hsv_config.push_back(hsv{
                      json_config[i]["hue"],
                      json_config[i]["saturation"],
                      json_config[i]["value"],
                  });
                }
                int step_size = LED_COUNT / gradient_hsv_config.size();
                int hsv_index = 0;

                for (auto i = 0; i < led_strip.numPixels(); i++)
                {
                  if (i / (hsv_index + 1) >= step_size)
                    hsv_index++;
                  uint32_t col = led_strip.gamma32(led_strip.ColorHSV(
                      gradient_hsv_config[hsv_index].h,
                      gradient_hsv_config[hsv_index].s,
                      gradient_hsv_config[hsv_index].v));
                  led_strip.setPixelColor(i, col);
                }
                led_strip.show();
              }
              break;
              case Mode::ANIMATION_RGB:
                red = jsonDocument["red"];
                green = jsonDocument["green"];
                blue = jsonDocument["blue"];
                animation_speed = jsonDocument["animation_speed"];
                break;
              case Mode::ANIMATION_HSV:
                hue = jsonDocument["hue"];
                saturation = jsonDocument["saturation"];
                value = jsonDocument["value"];
                animation_speed = jsonDocument["animation_speed"];
                break;
              case Mode::INVALID:
                Serial.println("ERROR in /led");
                jsonDocument.clear();
                jsonDocument["error_message"] = "mode is not valid";
                serializeJson(jsonDocument, buffer);
                server.send(422, "application/json", buffer);
                return;
              default:
                Serial.println("ERROR in /led");
                jsonDocument.clear();
                jsonDocument["error_message"] = "mode is not valid";
                serializeJson(jsonDocument, buffer);
                server.send(422, "application/json", buffer);
                return;
              }

              jsonDocument.clear();
              jsonDocument["message"] = "successfully changed colour of the led";
              serializeJson(jsonDocument, buffer);
              server.send(200, "application/json", buffer);
            });

  server.on("/info", HTTP_GET, []()
            {
              jsonDocument.clear();
              jsonDocument["device_name"] = device_name.c_str();
              jsonDocument["device_type"] = "RGB_LED_STRIP";
              jsonDocument["device_info"] = jsonDocument.createNestedObject();
              jsonDocument["device_info"]["current_mode"] = mode_to_string(mode);
              jsonDocument["device_info"]["supported_modes"] = jsonDocument.createNestedArray();
              jsonDocument["device_info"]["supported_modes"].add(mode_to_string(Mode::SINGLE_COLOR_RGB));
              jsonDocument["device_info"]["supported_modes"].add(mode_to_string(Mode::SINGLE_COLOR_HSV));
              jsonDocument["device_info"]["supported_modes"].add(mode_to_string(Mode::GRADIENT_RGB));
              jsonDocument["device_info"]["supported_modes"].add(mode_to_string(Mode::GRADIENT_HSV));
              jsonDocument["device_info"]["supported_modes"].add(mode_to_string(Mode::ANIMATION_RGB));
              jsonDocument["device_info"]["supported_modes"].add(mode_to_string(Mode::ANIMATION_HSV));
              jsonDocument["device_info"]["data"] = jsonDocument.createNestedObject();
              switch (mode)
              {
              case SINGLE_COLOR_RGB:
              {
                jsonDocument["device_info"]["data"]["red"] = red;
                jsonDocument["device_info"]["data"]["green"] = green;
                jsonDocument["device_info"]["data"]["blue"] = blue;
              }
              break;
              case SINGLE_COLOR_HSV:
              {
                jsonDocument["device_info"]["data"]["hue"] = hue;
                jsonDocument["device_info"]["data"]["saturation"] = saturation;
                jsonDocument["device_info"]["data"]["value"] = value;
              }
              break;
              case ANIMATION_RGB:
              {
                jsonDocument["device_info"]["data"]["red"] = red;
                jsonDocument["device_info"]["data"]["green"] = green;
                jsonDocument["device_info"]["data"]["blue"] = blue;
                jsonDocument["device_info"]["data"]["animation_speed"] = animation_speed;
              }
              break;
              case ANIMATION_HSV:
              {
                jsonDocument["device_info"]["data"]["hue"] = hue;
                jsonDocument["device_info"]["data"]["saturation"] = saturation;
                jsonDocument["device_info"]["data"]["value"] = value;
                jsonDocument["device_info"]["data"]["animation_speed"] = animation_speed;
              }
              break;
              case GRADIENT_RGB:
              {
                jsonDocument["device_info"]["data"]["gradients"] = jsonDocument.createNestedArray();
                for (auto c : gradient_rgb_config)
                {
                  auto obj = jsonDocument["device_info"]["data"]["gradients"].createNestedObject();
                  obj["red"] = c.r;
                  obj["green"] = c.g;
                  obj["blue"] = c.b;
                }
              }
              break;
              case GRADIENT_HSV:
              {
                jsonDocument["device_info"]["data"]["gradients"] = jsonDocument.createNestedArray();
                for (auto c : gradient_hsv_config)
                {
                  auto obj = jsonDocument["device_info"]["data"]["gradients"].createNestedObject();
                  obj["hue"] = c.h;
                  obj["saturation"] = c.s;
                  obj["value"] = c.v;
                }
              }
              break;
              default:
                break;
              }
              serializeJson(jsonDocument, buffer);
              server.send(200, "application/json", buffer);
            });

  server.begin();
  // init led strip
  led_strip.begin();
  led_strip.fill(led_strip.gamma32(led_strip.Color(225, 125, 15)), 0, led_strip.numPixels() - 1);
  led_strip.show();
  led_strip_12v->set_color(225, 125, 15);
  led_strip_12v->show();

  timer = millis();
}

int current_led = 0;

void loop()
{
  server.handleClient();
  switch (mode)
  {
  case Mode::ANIMATION_RGB:
  {
    if (millis() - timer > animation_speed)
    {
      if (current_led > 0)
        led_strip.setPixelColor(current_led - 1, 0, 0, 0);

      if (current_led >= led_strip.numPixels())
        current_led = 0;
      uint32_t col = led_strip.gamma32(led_strip.Color(red, green, blue));
      led_strip.setPixelColor(current_led, col);
      led_strip.show();
      current_led++;

      timer = millis();
    }
  }
  break;
  case Mode::ANIMATION_HSV:
  {
    if (millis() - timer > animation_speed)
    {
      if (current_led > 0)
        led_strip.setPixelColor(current_led - 1, 0, 0, 0);

      if (current_led >= led_strip.numPixels())
        current_led = 0;
      uint32_t col = led_strip.gamma32(led_strip.ColorHSV(hue, saturation, value));
      led_strip.setPixelColor(current_led, col);
      led_strip.show();
      current_led++;

      timer = millis();
    }
  }
  break;
  default:
    break;
  }
}