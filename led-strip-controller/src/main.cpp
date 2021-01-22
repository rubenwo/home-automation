#include <Arduino.h>
#include <WiFi.h>
#include <HTTPClient.h>
#include <WebServer.h>
#include <ArduinoJson.h>
#include <Adafruit_NeoPixel.h>

const char *ssid = "";
const char *password = "";

WebServer server(80);

Adafruit_NeoPixel led_strip(150, 4, NEO_GRB + NEO_KHZ800);

// JSON data buffer
StaticJsonDocument<500> jsonDocument;
char buffer[500];
bool on = false;
void setup()
{
  Serial.begin(460800);
  Serial.printf("Connecting to: %s\n", ssid);

  // setup wifi
  WiFi.begin(ssid, password);

  while (WiFi.status() != WL_CONNECTED)
  {
    Serial.print(".");
    delay(500);
  }
  Serial.printf(" Connected. IP: %s\n", WiFi.localIP().toString().c_str());

  // setup routing web server
  server.on("/healthz", HTTP_GET, []() {
    Serial.println("/healthz");
    jsonDocument.clear();
    jsonDocument["isHealthy"] = true;
    jsonDocument["error_message"] = "";
    serializeJson(jsonDocument, buffer);
    server.send(200, "application/json", buffer);
  });
  pinMode(2, OUTPUT);
  server.on("/led", HTTP_POST, []() {
    Serial.println("/led");
    if (server.hasArg("plain") == false)
    {
      Serial.println("ERROR in /led");
      jsonDocument.clear();
      jsonDocument["error_message"] = "body is not valid";
      serializeJson(jsonDocument, buffer);
      server.send(422, "application/json", buffer);
    }

    String body = server.arg("plain");
    deserializeJson(jsonDocument, body);

    // Get RGB components
    int red = jsonDocument["red"];
    int green = jsonDocument["green"];
    int blue = jsonDocument["blue"];
    Serial.printf("R: %d, G: %d, B: %d\n", red, green, blue);

    for (auto i = 0; i < led_strip.numPixels(); i++)
    {
      led_strip.setPixelColor(i, red, green, blue);
    }
    led_strip.show();
    jsonDocument.clear();
    jsonDocument["message"] = "successfully changed colour of the led";
    serializeJson(jsonDocument, buffer);
    server.send(200, "application/json", buffer);
  });

  server.on("/info", HTTP_GET, []() {
    jsonDocument.clear();
    jsonDocument["device_id"] = "123-456-789-xyz";
    jsonDocument["device_name"] = "led_strip_a";
    jsonDocument["device_type"] = "RGB_LED_STRIP";
    jsonDocument["device_info"] = jsonDocument.createNestedObject();
    jsonDocument["device_info"]["mode"] = "single"; // implement gradient, animation, etc.
    jsonDocument["device_info"]["R"] = 127;
    jsonDocument["device_info"]["G"] = 0;
    jsonDocument["device_info"]["B"] = 255;

    serializeJson(jsonDocument, buffer);
    server.send(200, "application/json", buffer);
  });

  server.begin();

  // init led strip
  led_strip.begin();
  int r, g, b;
  r = 245;
  g = 149;
  b = 24;
  for (auto i = 0; i < led_strip.numPixels(); i++)
  {
    led_strip.setPixelColor(i, r, g, b);
  }
  led_strip.show();

  // announce online
  HTTPClient http;
  http.begin("http://192.168.2.135/api/v1/new_id");
  int httpResponseCode = http.GET();
  String payload = "{}";
  if (httpResponseCode > 0)
  {
    Serial.printf("HTTP Response code: %d\n", httpResponseCode);
    payload = http.getString();
  }
  else
  {
    Serial.printf("Error code: %d\n", httpResponseCode);
  }
  http.end();
  Serial.println(payload);

  deserializeJson(jsonDocument, payload);
  String val = jsonDocument["id"];
  Serial.println(val);
}

void loop()
{
  server.handleClient();
}