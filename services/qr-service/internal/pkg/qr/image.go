package qr

import (
	"github.com/makiuchi-d/gozxing"
	"image"
	"image/color"
	"image/draw"
	"math"
)

var (
	white = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	black = color.RGBA{R: 0, G: 0, B: 0, A: 255}
)

// convertImageToRGBA converts a generic image.Image to a specific image.RGBA format for performance.
func convertImageToRGBA(img image.Image) *image.RGBA {
	if rgba, ok := img.(*image.RGBA); ok {
		return rgba
	}
	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, img.Bounds().Min, draw.Src)
	return rgba
}

// mat2d is a simple matrix data-structure. Contains the width and height and the data. data is a slice of color.RGBA with a size of w*h.
type mat2d struct {
	w, h int
	data []color.RGBA
}

// mat2dFromRGBA converts an image.RGBA to our own mat2d structure.
func mat2dFromRGBA(img *image.RGBA) *mat2d {
	w := img.Bounds().Dx()
	h := img.Bounds().Dy()
	mat := &mat2d{
		w:    w,
		h:    h,
		data: make([]color.RGBA, w*h),
	}

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			mat.data[(h*y)+x] = img.RGBAAt(x, y)
		}
	}

	return mat
}

// mat2dFromBitMatrix converts the given BitMatrix to our own mat2d.
func mat2dFromBitMatrix(m *gozxing.BitMatrix) *mat2d {
	w := m.GetWidth()
	h := m.GetHeight()
	mat := &mat2d{
		w:    w,
		h:    h,
		data: make([]color.RGBA, w*h),
	}

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			col := white
			if m.Get(x, y) {
				col = black
			}
			mat.data[(h*y)+x] = col
		}
	}
	return mat
}

// mat2dFromBitMatrix converts the given BitMatrix to our own mat2d. It scales the output matrix using the scaleF parameter.
func mat2dFromBitMatrixWithScale(m *gozxing.BitMatrix, scaleF float64) *mat2d {
	w := m.GetWidth()
	h := m.GetHeight()

	if scaleF <= 1 {
		// Since a scale lower than 1 does not make any sense (BitMatrix is already the smallest possible) we skip the scaling
		return mat2dFromBitMatrix(m)
	}

	// We ceil the values and convert the scale into an integer for more efficient scaling
	scale := int(math.Ceil(scaleF))

	scaleW := w * scale
	scaleH := h * scale

	data := make([]color.RGBA, scaleW*scaleH)
	// memset repeat algorithm (log(n))
	data[0] = white
	for bp := 1; bp < len(data); bp *= 2 {
		copy(data[bp:], data[:bp])
	}

	// Create a row of black pixels of size 'scale'. This is the horizontal part of the scaled pixel.
	// Benchmarks show a performance increase even at 2x scaling (which is the minimum requirement for this code to be run)
	row := make([]color.RGBA, scale)
	for i := range row {
		row[i] = black
	}

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if !m.Get(x, y) {
				continue
			}

			// calculate the place where the black row should be set
			yStart := y * scale
			xStart := x * scale
			xEnd := xStart + scale

			// Set the pixels in the data for the scaled-up block to black
			for yy := yStart; yy < yStart+scale; yy++ {
				// Use copy to set multiple pixels to the same (black) value. This is more important on larger images.
				copy(data[(scaleH*yy)+xStart:(scaleH*yy)+xEnd], row)
			}
		}
	}

	return &mat2d{
		w:    scaleW,
		h:    scaleH,
		data: data,
	}
}

// overlayWithAlphaBlending blends the colors of m1 and m2  where the overlay. This is done based on the 'alpha' value
// of the color.(startX, startY) are the upper-left is the upper-left corner where the function will start to overlay m2 onto m1.
func overlayWithAlphaBlending(m1 *mat2d, m2 *mat2d, startX, startY int) {
	const m = 1<<16 - 1
	for y := 0; y < m2.h; y++ {
		for x := 0; x < m2.w; x++ {
			sr, sg, sb, sa := m2.AtC(x, y)
			dr, dg, db, da := m1.AtC(x+startX, y+startY)
			a := m - sa
			m1.data[(m1.h*(y+startY) + (x + startX))] = color.RGBA{
				R: uint8(uint16((dr*a + sr*m) / m)),
				G: uint8(uint16((dg*a + sg*m) / m)),
				B: uint8(uint16((db*a + sb*m) / m)),
				A: uint8(uint16((da*a + sa*m) / m)),
			}
		}
	}
}

// overlay does not use alpha blending to maximise performance. This is useful for overlays with opaque backgrounds.
// (startX, startY) is the upper-left corner where the function will start to copy m2 into m1.
func overlay(m1 *mat2d, m2 *mat2d, startX, startY int) {
	for y := 0; y < m2.h; y++ {
		// copy m2 into m1 row-wise.
		copy(m1.data[(m1.h*(y+startY))+startX:(m1.h*(y+startY)+m2.w+startX)], m2.data[m2.h*y:m2.h*y+m2.w])
	}
}

func (img *mat2d) ColorModel() color.Model {
	return color.RGBAModel
}

func (img *mat2d) Bounds() image.Rectangle {
	return image.Rect(0, 0, img.w, img.h)
}

func (img *mat2d) AtC(x, y int) (r uint32, g uint32, b uint32, a uint32) {
	c := img.data[img.h*y+x]
	r = uint32(c.R)
	r |= r << 8
	g = uint32(c.G)
	g |= g << 8
	b = uint32(c.B)
	b |= b << 8
	a = uint32(c.A)
	a |= a << 8
	return
}

func (img *mat2d) At(x, y int) color.Color {
	return img.data[(img.h*y)+x]
}

//func (img *mat2d) ToPNG() []byte {
//	src24Bit := make([]uint8, img.w*img.h*3)
//	for y := 0; y < img.h; y++ {
//		for x := 0; x < img.w; x++ {
//			c := img.data[(img.h*y)+x]
//			src24Bit[img.h*y+x+0] = c.R
//			src24Bit[img.h*y+x+1] = c.G
//			src24Bit[img.h*y+x+2] = c.B
//		}
//	}
//
//	return lib.EncodeFromMat2d(src24Bit, uint32(img.w), uint32(img.h), img.w)
//}
