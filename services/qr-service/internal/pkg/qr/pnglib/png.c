#include "png.h"
#include "wrapper.h"

#include <stdio.h>
#include <stdlib.h>

unsigned char *fpng_encode_with_data_wrapper(uint32_t *in, uint32_t w, uint32_t h, int size, size_t *outSize)
{
    return fpng_encode_with_data(in, w, h, size, outSize);
}

unsigned char *fpng_encode_with_data_wrapper8(uint8_t *in, uint32_t w, uint32_t h, int size, size_t *outSize)
{
    return fpng_encode_with_data8(in, w, h, size, outSize);
}

void fpng_init_wrapper()
{
    fpng_init_wrap();
}