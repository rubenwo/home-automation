package qr

// The following code was copied and modified from the Go source code. We now use a *mat2d data structure that holds
// the color values. This is a lot faster than calling interface functions. Total time reduction of >2x

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
import (
	"bufio"
	"compress/zlib"
	"encoding/binary"
	"hash/crc32"
	"io"
)

const pngHeader = "\x89PNG\r\n\x1a\n"

// PNGEncoder configures encoding PNG images.
type PNGEncoder struct {
	CompressionLevel CompressionLevel

	// BufferPool optionally specifies a buffer pool to get temporary
	// EncoderBuffers when encoding an image.
	BufferPool EncoderBufferPool
}

// EncoderBufferPool is an interface for getting and returning temporary
// instances of the EncoderBuffer struct. This can be used to reuse buffers
// when encoding multiple images.
type EncoderBufferPool interface {
	Get() *EncoderBuffer
	Put(*EncoderBuffer)
}

// EncoderBuffer holds the buffers used for encoding PNG images.
type EncoderBuffer pngEncoder

type pngEncoder struct {
	enc     *PNGEncoder
	w       io.Writer
	m       *mat2d
	err     error
	header  [8]byte
	footer  [4]byte
	tmp     [4 * 256]byte
	cr      [nFilter][]uint8
	pr      []uint8
	zw      *zlib.Writer
	zwLevel int
	bw      *bufio.Writer
}

// CompressionLevel indicates the compression level.
type CompressionLevel int

const (
	DefaultCompression CompressionLevel = 0
	NoCompression      CompressionLevel = -1
	BestSpeed          CompressionLevel = -2
	BestCompression    CompressionLevel = -3

	// Positive CompressionLevel values are reserved to mean a numeric zlib
	// compression level, although that is not implemented yet.
)

// Color type, as per the PNG spec.
const (
	ctTrueColor = 2
)

// Filter type, as per the PNG spec.
const (
	ftNone  = 0
	nFilter = 5
)

func (e *pngEncoder) writeChunk(b []byte, name string) {
	if e.err != nil {
		return
	}
	n := uint32(len(b))

	binary.BigEndian.PutUint32(e.header[:4], n)
	e.header[4] = name[0]
	e.header[5] = name[1]
	e.header[6] = name[2]
	e.header[7] = name[3]
	crc := crc32.NewIEEE()
	_, _ = crc.Write(e.header[4:8])
	_, _ = crc.Write(b)
	binary.BigEndian.PutUint32(e.footer[:4], crc.Sum32())

	_, e.err = e.w.Write(e.header[:8])
	if e.err != nil {
		return
	}
	_, e.err = e.w.Write(b)
	if e.err != nil {
		return
	}
	_, e.err = e.w.Write(e.footer[:4])
}

func (e *pngEncoder) writeIHDR() {
	b := e.m.Bounds()
	binary.BigEndian.PutUint32(e.tmp[0:4], uint32(b.Dx()))
	binary.BigEndian.PutUint32(e.tmp[4:8], uint32(b.Dy()))
	// Set bit depth and color type.
	e.tmp[8] = 8
	e.tmp[9] = ctTrueColor
	e.tmp[10] = 0 // default compression method
	e.tmp[11] = 0 // default filter method
	e.tmp[12] = 0 // non-interlaced
	e.writeChunk(e.tmp[:13], "IHDR")
}

// An encoder is an io.Writer that satisfies writes by writing PNG IDAT chunks,
// including an 8-byte header and 4-byte CRC checksum per Write call. Such calls
// should be relatively infrequent, since writeIDATs uses a bufio.Writer.
//
// This method should only be called from writeIDATs (via writeImage).
// No other code should treat an encoder as an io.Writer.
func (e *pngEncoder) Write(b []byte) (int, error) {
	e.writeChunk(b, "IDAT")
	if e.err != nil {
		return 0, e.err
	}
	return len(b), nil
}

func zeroMemory(v []uint8) {
	v[0] = 0
	for bp := 1; bp < len(v); bp *= 2 {
		copy(v[bp:], v[:bp])
	}
}

func (e *pngEncoder) writeImage(w io.Writer, m *mat2d, level int) error {
	if e.zw == nil || e.zwLevel != level {
		zw, err := zlib.NewWriterLevel(w, level)
		if err != nil {
			return err
		}
		e.zw = zw
		e.zwLevel = level
	} else {
		e.zw.Reset(w)
	}
	defer e.zw.Close()

	bitsPerPixel := 24

	// cr[*] and pr are the bytes for the current and previous row.
	// cr[0] is unfiltered (or equivalently, filtered with the ftNone filter).
	// cr[ft], for non-zero filter types ft, are buffers for transforming cr[0] under the
	// other PNG filter types. These buffers are allocated once and re-used for each row.
	// The +1 is for the per-row filter type, which is at cr[*][0].
	b := m.Bounds()
	sz := 1 + (bitsPerPixel*b.Dx()+7)/8
	for i := range e.cr {
		if cap(e.cr[i]) < sz {
			e.cr[i] = make([]uint8, sz)
		} else {
			e.cr[i] = e.cr[i][:sz]
		}
		e.cr[i][0] = uint8(i)
	}
	cr := e.cr
	if cap(e.pr) < sz {
		e.pr = make([]uint8, sz)
	} else {
		e.pr = e.pr[:sz]
		zeroMemory(e.pr)
	}
	pr := e.pr

	for y := b.Min.Y; y < b.Max.Y; y++ {
		// Convert from colors to bytes.
		i := 1

		// We have previously verified that the alpha value is fully opaque.
		cr0 := cr[0]

		for x := b.Min.X; x < b.Max.X; x++ {
			c := m.data[m.h*y+x]
			cr0[i+0] = c.R
			cr0[i+1] = c.G
			cr0[i+2] = c.B
			i += 3
		}

		// Apply the filter.
		// Skip filter for NoCompression and paletted images (cbP8) as
		// "filters are rarely useful on palette images" and will result
		// in larger files (see http://www.libpng.org/pub/png/book/chapter09.html).
		f := ftNone

		// Write the compressed bytes.
		if _, err := e.zw.Write(cr[f]); err != nil {
			return err
		}

		// The current row for y is the previous row for y+1.
		pr, cr[0] = cr[0], pr
	}
	return nil
}

// Write the actual image data to one or more IDAT chunks.
func (e *pngEncoder) writeIDATs() {
	if e.err != nil {
		return
	}
	if e.bw == nil {
		e.bw = bufio.NewWriterSize(e, 1<<15)
	} else {
		e.bw.Reset(e)
	}
	e.err = e.writeImage(e.bw, e.m, levelToZlib(e.enc.CompressionLevel))
	if e.err != nil {
		return
	}
	e.err = e.bw.Flush()
}

// This function is required because we want the zero value of
// Encoder.CompressionLevel to map to zlib.DefaultCompression.
func levelToZlib(l CompressionLevel) int {
	switch l {
	case DefaultCompression:
		return zlib.DefaultCompression
	case NoCompression:
		return zlib.NoCompression
	case BestSpeed:
		return zlib.BestSpeed
	case BestCompression:
		return zlib.BestCompression
	default:
		return zlib.DefaultCompression
	}
}

func (e *pngEncoder) writeIEND() { e.writeChunk(nil, "IEND") }

// Encode writes the Image(*mat2d) m to w in PNG format.
func (enc *PNGEncoder) Encode(w io.Writer, m *mat2d) error {
	var e *pngEncoder
	if enc.BufferPool != nil {
		buffer := enc.BufferPool.Get()
		e = (*pngEncoder)(buffer)

	}
	if e == nil {
		e = &pngEncoder{}
	}
	if enc.BufferPool != nil {
		defer enc.BufferPool.Put((*EncoderBuffer)(e))
	}

	e.enc = enc
	e.w = w
	e.m = m

	_, e.err = io.WriteString(w, pngHeader)
	e.writeIHDR()
	e.writeIDATs()
	e.writeIEND()
	return e.err
}
