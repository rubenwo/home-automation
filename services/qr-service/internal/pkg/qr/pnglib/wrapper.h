// wrapper.h
#ifdef __cplusplus
extern "C"
{
#endif
#include <stdlib.h>
#include <stdint.h>

    unsigned char *fpng_encode_with_data(uint32_t *in, uint32_t w, uint32_t h, int size, size_t *outSize);
    unsigned char *fpng_encode_with_data8(uint8_t *in, uint32_t w, uint32_t h, int size, size_t *outSize);
    void fpng_init_wrap();

#ifdef __cplusplus
}
#endif
