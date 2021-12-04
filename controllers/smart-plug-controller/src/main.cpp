#include <Arduino.h>
#include <WiFi.h>
#include <ArduinoJson.h>
#include <map>
#include <vector>
#include <functional>

#include "PubSubClient.h"
#include "strings.hpp"

#define ONBOARD_LED_PIN 2

std::map<std::string, std::function<void(byte *message, unsigned int length)>> topic_router;
std::map<int32_t, int32_t> pins;

WiFiClient wifi_client;
PubSubClient mqtt_client(wifi_client);

const char *ssid = "";
const char *password = "";
const char *mqtt_broker = "";
const char *device_name = "";
const char *device_id = "";

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

    mqtt_client.publish(string_format("smartplug/%s/response", device_id).c_str(), buffer);
}

void control_plug(byte *message, unsigned int length)
{
    DynamicJsonDocument jsonDocument(length * sizeof(byte));
    deserializeJson(jsonDocument, message, length);

    const std::string &json_mode = jsonDocument["mode"];
    const int32_t &json_pin = jsonDocument["pin"];

    auto pin = pins[json_pin];
    if (pin)
    {
        if (json_mode == "on")
        {
            digitalWrite(pin, HIGH);
        }
        else if (json_mode == "off")
        {
            digitalWrite(pin, LOW);
        }
    }
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

void setup()
{
    Serial.begin(460800);
    delay(1000);
    pinMode(ONBOARD_LED_PIN, OUTPUT);
    // Turn on onboard LED to indicate the setup
    digitalWrite(ONBOARD_LED_PIN, HIGH);

    topic_router[string_format("smartplug/%s/information", device_id)] = get_information;
    topic_router[string_format("smartplug/%s/control", device_id)] = control_plug;

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
}