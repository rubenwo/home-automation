#pragma once

#include <stdio.h>
#include <string>

enum Mode
{
    INVALID = -1,
    SINGLE_COLOR_RGB = 0,
    GRADIENT_RGB = 1,
    ANIMATION_RGB = 2,
};


Mode from_string(const std::string &str)
{
    if (str == "SINGLE_COLOR_RGB")
        return Mode::SINGLE_COLOR_RGB;
    else if (str == "GRADIENT_RGB")
        return Mode::GRADIENT_RGB;
    else if (str == "ANIMATION_RGB")
        return Mode::ANIMATION_RGB;

    return Mode::INVALID;
}

std::string mode_to_string(const Mode &m)
{
    switch (m)
    {
    case Mode::SINGLE_COLOR_RGB:
        return "SINGLE_COLOR_RGB";
    case Mode::GRADIENT_RGB:
        return "GRADIENT_RGB";
    case Mode::ANIMATION_RGB:
        return "ANIMATION_RGB";
    case Mode::INVALID:
        return "INVALID";
    }
    return "INVALID";
}
