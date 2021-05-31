#pragma once

struct rgb
{
    int r, g, b;
};

struct hsv
{
    int h, s, v;
};

class Color
{
public:
    Color() {}
    ~Color() {}

    void set_rgb(rgb rgb) {}
    void set_hsv(hsv hsv) {}

// TODO: 
    rgb get_rgb() { return this->rgb_color; }
    hsv get_hsv() { return this->hsv_color; }

private:
    rgb rgb_color;
    hsv hsv_color;
};