package qr

import (
	"bytes"
	"fmt"
	"github.com/dgraph-io/ristretto"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	"github.com/makiuchi-d/gozxing/qrcode/decoder"
	"github.com/nfnt/resize"
	"image"
	"io"
)

const (
	// minCacheSize fits 1 2000x2000 image
	minCacheSize = 20000000
)

type Generator interface {
	GenerateQRCode(content string, size int, margin int, outputFormat OutputFormat, opts ...options) (io.WriterTo, error)
}

type generator struct {
	// idealQRLogo is the original logo which can be resized. The result of this resize is stored in the resizeCacheMat2d
	idealQRLogo *image.RGBA
	// resizeCacheMat2d is a ristretto cache with eviction based on cost
	resizeCacheMat2d *ristretto.Cache

	pngEncoder *PNGEncoder
}

// calculateCost returns the size of the mat2d in bytes. This function is used by the cache for eviction purposes.
func calculateCost(m *mat2d) int64 {
	const (
		sInt32 = 4 // size of an int32 in bytes
		sInt64 = 8 // size of an int64 in bytes
	)
	// int alias (m.w and m.h) will be int64 since we compile to a 64-bit binary
	return int64((len(m.data))*sInt32 + 2*sInt64)
}

// NewGenerator will create a new generator that creates iDEAL-qr qr-codes. The given image is the iDEAL logo.
// This image is converted into a *mat2d data-structure to efficiently overlay the logo on top off the QR-code.
// If maxCacheSizeInBytes is smaller than (minCacheSize = 20000000), the minCacheSize is used instead.
func NewGenerator(idealQRLogo image.Image, maxCacheSizeInBytes int64) Generator {
	if maxCacheSizeInBytes < minCacheSize {
		maxCacheSizeInBytes = minCacheSize
	}

	rgbaLogo := convertImageToRGBA(idealQRLogo)
	if idealQRLogo.Bounds().Dx() != idealQRLogo.Bounds().Dy() {
		rgbaLogo = resize.Resize(uint(idealQRLogo.Bounds().Dx()), uint(idealQRLogo.Bounds().Dx()), rgbaLogo, resize.Bilinear).(*image.RGBA)
	}

	logoMatrix := mat2dFromRGBA(rgbaLogo)

	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 2048, // We expect a max of 1800 cache entries. (2000-200)
		MaxCost:     maxCacheSizeInBytes,
		BufferItems: 64, // Use the recommended setting here.
	})

	if err != nil {
		panic(fmt.Sprintf("ristretto cache initialization: %s", err.Error()))
	}

	cache.Set(idealQRLogo.Bounds().Dx()*idealQRLogo.Bounds().Dy(), logoMatrix, calculateCost(logoMatrix))

	return &generator{
		idealQRLogo:      rgbaLogo,
		resizeCacheMat2d: cache,
		pngEncoder:       &PNGEncoder{CompressionLevel: BestSpeed},
	}
}

// GenerateQRCode will generate a qr-code for the ideal-qr spec. This contains a logo and the highest error correction.
// SVGs are about 300x faster to generate than PNGs at 2000x2000. SVGs are generated in constant time depending on the content size.
// PNGs and JPGs scale nearly linearly with the requested size (i.e. 2000x2000 takes about 100 times longer than 200x200).
// When using SVG as OutputFormat the given size is the size of the viewport.
// By returning an io.WriterTo interface we can reduce our allocations by not getting a copy of the buffer's slice.
func (g *generator) GenerateQRCode(content string, size int, margin int, outputFormat OutputFormat, opts ...options) (io.WriterTo, error) {
	cfg := &config{
		logoSize:   DefaultLogoSize,
		blendAlpha: false,
	}

	for _, opt := range opts {
		if opt != nil {
			opt(cfg)
		}
	}

	// By setting the width and height to 0 in the 'Encode' function call we get the smallest qr-code possible.
	// This is fine because we can scale the image efficiently to the desired 'size'.
	bitMatrix, err := qrcode.NewQRCodeWriter().Encode(content, gozxing.BarcodeFormat_QR_CODE, 0, 0, map[gozxing.EncodeHintType]interface{}{
		gozxing.EncodeHintType_ERROR_CORRECTION: decoder.ErrorCorrectionLevel_H,
		gozxing.EncodeHintType_MARGIN:           margin,
	})
	if err != nil {
		return nil, &ErrFailedToCreateQRCode{Err: err}
	}

	// At the moment, our only non-rasterized output format is an SVG, but make it a variable anyway for better readability
	isRasterizedOutput := outputFormat != SVG

	// outputImage is only set if the output is a rasterized image (i.e. PNG or JPG)
	var outputImage *mat2d
	if isRasterizedOutput {
		scale := float64(size) / float64(bitMatrix.GetWidth())
		if scale < 1 {
			scale = 1
			size = bitMatrix.GetHeight()
		}

		// Use a custom scaling algorithm (and our own matrix data-structure) for lower resource usage at the same performance
		// Benchmarked at 1/3 memory usage and 1/2 cpu usage compared to NearestNeighbour from resize.Resize.
		resizedMatrix := mat2dFromBitMatrixWithScale(bitMatrix, scale)

		logoX, logoY := int(float64(resizedMatrix.Bounds().Dx())*cfg.logoSize), int(float64(resizedMatrix.Bounds().Dy())*cfg.logoSize)
		key := logoX * logoY
		logoInterface, ok := g.resizeCacheMat2d.Get(key)
		if !ok {
			// Here we can leverage the resize.Resize function from the library as the Logo contains 'complex' shapes and colours.
			// There is a subtle bug here because of a race-condition. However, it might not be worth it to use locks
			// here since the value for a key won't change after it has been set.
			logo := mat2dFromRGBA(resize.Resize(uint(logoX), uint(logoY), g.idealQRLogo, resize.Bilinear).(*image.RGBA))
			g.resizeCacheMat2d.Set(key, logo, calculateCost(logo))
			logoInterface = logo
		}
		logo := logoInterface.(*mat2d)
		centerX, centerY := resizedMatrix.Bounds().Dx()/2, resizedMatrix.Bounds().Dx()/2
		if cfg.blendAlpha {
			overlayWithAlphaBlending(resizedMatrix, logo, centerX-(logo.w/2), centerY-(logo.h/2))
		} else {
			overlay(resizedMatrix, logo, centerX-(logo.w/2), centerY-(logo.h/2))
		}
		outputImage = resizedMatrix
	}

	out := new(bytes.Buffer)

	switch outputFormat {
	case PNG:
		//We use our own encoder which is based on the stdlib, but removed every piece of code we don't need.
		err := g.pngEncoder.Encode(out, outputImage)
		if err != nil {
			return nil, &ErrEncodingPNG{Err: err}
		}
		//out.Write(outputImage.ToPNG())
	case JPG:
		err := JPEGEncode(out, outputImage, nil)
		if err != nil {
			return nil, &ErrEncodingJPEG{Err: err}
		}
	case SVG:
		generateSVG(out, bitMatrix, size, cfg.logoSize)
	}

	return out, nil
}
