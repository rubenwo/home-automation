#include <Arduino.h>
#include <WiFi.h>
#include <WebServer.h>

const char *ssid = "";
const char *password = "";

WebServer server(80);

void setup()
{
  Serial.begin(460800);
  Serial.printf("Connecting to: %s\n", ssid);

  

}

void loop()
{
  // put your main code here, to run repeatedly:
}