#ifndef _PNG_H
#define _PNG_H
#include <stdlib.h>
#include <stdint.h>

unsigned char *fpng_encode_with_data_wrapper(uint32_t *in, uint32_t w, uint32_t h, int size, size_t *outSize);
unsigned char *fpng_encode_with_data_wrapper8(uint8_t *in, uint32_t w, uint32_t h, int size, size_t *outSize);
void fpng_init_wrapper();

#endif