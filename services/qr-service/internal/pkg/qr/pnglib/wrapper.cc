// wrapper.cc
#include "wrapper.h"
#include <iostream>
#include "fpng.h"

typedef std::vector<uint8_t> uint8_vec;

unsigned char *fpng_encode_with_data(uint32_t *in, uint32_t w, uint32_t h, int size, size_t *outSize)
{
    std::vector<uint8_t> *out = new std::vector<uint8_t>();

    uint8_vec srcImageBuffer24(w * h * 3);
    for (uint32_t i = 0; i < w * h; i++)
    {
        srcImageBuffer24[i * 3 + 0] = in[i] >> 24;
        srcImageBuffer24[i * 3 + 1] = in[i] >> 16;
        srcImageBuffer24[i * 3 + 2] = in[i] >> 8;
    }

    const uint8_t *pImage = srcImageBuffer24.data();

    if (!fpng::fpng_encode_image_to_memory(pImage, w, h, 3, *out, 0))
    {
        std::cout << "Failed" << std::endl;
    }

    *outSize = out->size();

    return out->data();
}

unsigned char *fpng_encode_with_data8(uint8_t *in, uint32_t w, uint32_t h, int size, size_t *outSize)
{
    std::vector<uint8_t> *out = new std::vector<uint8_t>();

    if (!fpng::fpng_encode_image_to_memory(in, w, h, 3, *out, 0))
    {
        std::cout << "Failed" << std::endl;
    }

    *outSize = out->size();

    return out->data();
}

void fpng_init_wrap()
{
    fpng::fpng_init();
}