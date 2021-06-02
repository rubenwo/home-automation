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
    RGB12V(rgb12_config config)
        : red_channel_pin(config.red_channel_pin), green_channel_pin(config.green_channel_pin), blue_channel_pin(config.blue_channel_pin),
          red_pwm_channel(config.red_pwm_channel), green_pwm_channel(config.green_pwm_channel), blue_pwm_channel(config.blue_pwm_channel)
    {
        ledcSetup(red_pwm_channel, freq, resolution);
        ledcSetup(green_pwm_channel, freq, resolution);
        ledcSetup(blue_pwm_channel, freq, resolution);

        ledcAttachPin(red_channel_pin, red_pwm_channel);
        ledcAttachPin(green_channel_pin, green_pwm_channel);
        ledcAttachPin(blue_channel_pin, blue_pwm_channel);
    }
    ~RGB12V() {}

    void set_color(rgb rgb) { this->color.set_rgb(rgb); }
    void set_color_rgb(int r, int g, int b)
    {
        rgb c = {r, g, b};
        this->color.set_rgb(c);
    }

    void set_color(hsv hsv)
    {
        this->color.set_hsv(hsv);
    }
    void set_color_hsv(int h, int s, int v)
    {
        hsv c = {
            h,
            s,
            v,
        };
        this->color.set_hsv(c);
    }

    rgb get_color_rgb()
    {
        return this->color.get_rgb();
    }
    hsv get_color_hsv()
    {
        return this->color.get_hsv();
    }

    void show()
    {
        rgb c = this->color.get_rgb();
        if (c.r > 255)
            c.r = 255;
        if (c.b > 255)
            c.b = 255;
        if (c.g > 255)
            c.g = 255;
        if (c.r < 0)
            c.r = 0;
        if (c.b < 0)
            c.b = 0;
        if (c.g < 0)
            c.g = 0;

        ledcWrite(this->red_channel_pin, c.r);
        ledcWrite(this->green_channel_pin, c.b);
        ledcWrite(this->blue_channel_pin, c.g);
    }

private:
    const int freq = 5000;
    const int resolution = 8;

    Color color;

    uint8_t red_channel_pin;
    uint8_t green_channel_pin;
    uint8_t blue_channel_pin;

    uint8_t red_pwm_channel;
    uint8_t green_pwm_channel;
    uint8_t blue_pwm_channel;
};