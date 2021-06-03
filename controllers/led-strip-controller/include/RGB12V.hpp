#pragma once

#include <Arduino.h>

#include <tuple>

#include "color.hpp"

struct rgb12_config
{
    uint8_t red_channel_pin;
    uint8_t green_channel_pin;
    uint8_t blue_channel_pin;

    uint8_t red_pwm_channel;
    uint8_t green_pwm_channel;
    uint8_t blue_pwm_channel;
};

class RGB12V
{

public:
    RGB12V(uint8_t red_pin, uint8_t green_pin, uint8_t blue_pin, uint8_t red_channel, uint8_t green_channel, uint8_t blue_channel)
        : red_pin(red_pin), green_pin(green_pin), blue_pin(blue_pin), red_channel(red_channel), green_channel(green_channel), blue_channel(blue_channel)
    {
        ledcSetup(red_channel, freq, resolution);
        ledcAttachPin(red_pin, red_channel);

        ledcSetup(green_channel, freq, resolution);
        ledcAttachPin(green_pin, green_channel);

        ledcSetup(blue_channel, freq, resolution);
        ledcAttachPin(blue_pin, blue_channel);
    }
    ~RGB12V() {}

    void set_color(int red, int green, int blue)
    {
        this->red = red;
        this->green = green;
        this->blue = blue;
    }

    void show()
    {
        ledcWrite(red_channel, red);
        ledcWrite(green_channel, green);
        ledcWrite(blue_channel, blue);
    }


private:
   uint8_t red_pin;
    uint8_t green_pin;
    uint8_t blue_pin;
    uint8_t red_channel;
    uint8_t green_channel;
    uint8_t blue_channel;

    int red;
    int green;
    int blue;

    const int freq = 20000;
    const int resolution = 8;
};