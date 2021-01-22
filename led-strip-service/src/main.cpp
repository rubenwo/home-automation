#include <Arduino.h>
#include <WiFi.h>
#include <WebServer.h>
#include <ArduinoJson.h>
#include <Adafruit_NeoPixel.h>

#define NUM_OF_LEDS 8
#define PIN 4

const char *ssid = "";
const char *password = "";

WebServer server(80);

Adafruit_NeoPixel pixels(NUM_OF_LEDS, PIN, NEO_GRB + NEO_KHZ800);

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

    if (on)
    {
      digitalWrite(2, HIGH);
      on = false;
    }
    else
    {
      digitalWrite(2, LOW);
      on = true;
    }
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
    jsonDocument["device_info"]["R"] = 127;
    jsonDocument["device_info"]["G"] = 0;
    jsonDocument["device_info"]["B"] = 255;

    serializeJson(jsonDocument, buffer);
    server.send(200, "application/json", buffer);
  });

  server.begin();
  // init led strip
}

void loop()
{
  server.handleClient();
}