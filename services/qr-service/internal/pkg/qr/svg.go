package qr

import (
	"bytes"
	"fmt"
	"github.com/makiuchi-d/gozxing"
	"io"
)

// generateSVG writes an svg of the qr-code, scaling it to size and adds the ideal logo with the given logoSizePercentage
func generateSVG(w *bytes.Buffer, mat *gozxing.BitMatrix, size int, logoSizePercentage float64) {
	// Write the header for svg, this contains the svg width and height
	svgHeader(w, size, size)
	// Fill the background as a white rectangle
	w.WriteString("\t<rect width=\"100%\" height=\"100%\" fill=\"#FFFFFF\"/>\n")
	// Since the bitMatrix is a square we can take either dimension (w or h)
	_, _ = fmt.Fprintf(w, `<g transform="scale(%f)">`, float64(size)/float64(mat.GetWidth()))
	// Find and write the rectangles that make up the qr-code
	rects := FindRects(mat)
	for _, rect := range rects {
		// Write the black rectangle
		_, _ = fmt.Fprintf(w, `<rect x="%d" y="%d" width="%d" height="%d" style="fill:black;stroke:none"/>`, rect.X1, rect.Y1, rect.X2-rect.X1+1, rect.Y2-rect.Y1+1)
	}
	// Close the svg group for the black rectangles
	w.WriteString("</g>")

	// Calculate the scale of the logo. This scale is used to scale up or down the logo to logoSizePercentage of the total image
	maxSize := logoSizePercentage * float64(size)
	scaleX := maxSize / float64(idealLogoWidth)
	scaleY := maxSize / float64(idealLogoWidth)
	// Add the ideal Logo to the center of the screen with the desired scale
	idealLogoSVG(w, scaleX, scaleY, size/2-(int(idealLogoWidth*scaleX)/2), size/2-(int(idealLogoHeight*scaleY)/2))
	w.WriteString(`
</svg>
`)
}

// svgHeader writes and svg header to w containing the width and height for the view-port
func svgHeader(w io.Writer, width, height int) {
	_, _ = fmt.Fprintf(w, `<?xml version="1.0"?>
<svg width="%d" height="%d" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink">`, width, height)
}

// static values that came from the original svg ideal logo. These values are used to properly scale the logo with the qr-code.
const (
	idealLogoWidth  = 306
	idealLogoHeight = 269
)

// idealLogoSVG writes the ideal logo with the given scale and center translation to w
func idealLogoSVG(w io.Writer, scaleX, scaleY float64, centerX, centerY int) {
	_, _ = fmt.Fprintf(w, `
<style type="text/css">
	.st0{fill:#FFFFFF;}
	.st1{fill:#CC0066;}
</style>
<g transform="translate(%d,%d) scale(%f, %f)">
	<g>
		<path class="st0" d="M0,20v229.8c0,11,9,20,20,20h137.3c103.8,0,148.8-58.1,148.8-135.2C306.1,57.9,261.1,0,157.3,0H20
			C9,0,0,9,0,20z"/>
		<path class="st1" d="M91.9,56.4v169.8h73.9c67.1,0,96.2-37.9,96.2-91.5c0-51.3-29.1-91.1-96.2-91.1h-61.1
			C97.6,43.6,91.9,49.4,91.9,56.4z"/>
		<g>
			<g>
				<path d="M157.3,251.5H37.9c-10.6,0-19.2-8.6-19.2-19.2V37.6c0-10.6,8.6-19.2,19.2-19.2h119.4c113.3,0,130.2,72.9,130.2,116.3
					C287.5,210,241.2,251.5,157.3,251.5z M37.9,24.8c-7.1,0-12.8,5.7-12.8,12.8v194.7c0,7.1,5.7,12.8,12.8,12.8h119.4
					c79.8,0,123.8-39.2,123.8-110.4c0-95.6-77.6-109.9-123.8-109.9H37.9z"/>
			</g>
		</g>
	</g>
	<g>
		<path class="st0" d="M117.9,111.8c2.6,0,5,0.4,7.3,1.2c2.3,0.8,4.2,2.1,5.9,3.7c1.6,1.7,2.9,3.8,3.9,6.2c0.9,2.5,1.4,5.4,1.4,8.8
			c0,3-0.4,5.7-1.1,8.2c-0.8,2.5-1.9,4.7-3.4,6.5c-1.5,1.8-3.4,3.2-5.7,4.3c-2.3,1-5,1.6-8.1,1.6h-17.5v-40.6H117.9z M117.3,144.9
			c1.3,0,2.5-0.2,3.8-0.6c1.2-0.4,2.3-1.1,3.2-2.1c0.9-1,1.7-2.2,2.3-3.8c0.6-1.6,0.9-3.4,0.9-5.7c0-2-0.2-3.9-0.6-5.5
			c-0.4-1.6-1.1-3.1-2-4.2s-2.1-2.1-3.6-2.7c-1.5-0.6-3.3-0.9-5.5-0.9h-6.4v25.6H117.3z"/>
		<path class="st0" d="M172.5,111.8v7.5h-21.4v8.7h19.7v6.9h-19.7v9.9H173v7.5h-30.8v-40.6H172.5z"/>
		<path class="st0" d="M203.1,111.8l15.2,40.6H209l-3.1-9h-15.2l-3.2,9h-9l15.3-40.6H203.1z M203.6,136.7l-5.1-14.9h-0.1l-5.3,14.9
			H203.6z"/>
		<path class="st0" d="M232.8,111.8v33.1h19.8v7.5h-28.7v-40.6H232.8z"/>
	</g>
	<g>
		<circle cx="58.5" cy="132.1" r="18.7"/>
	</g>
	<path d="M72.6,226.2L72.6,226.2c-15.7,0-28.3-12.7-28.3-28.3v-22.1c0-7.8,6.3-14.2,14.2-14.2h0c7.8,0,14.2,6.3,14.2,14.2V226.2z"/>
</g>
`, centerX, centerY, scaleX, scaleY)
}

type InputMatrix interface {
	GetWidth() int
	GetHeight() int
	Get(int, int) bool
}

// matrix is a simple 2d data structure that is twice as efficient as a gozxing.BitMatrix in the FindRects algorithm
type matrix struct {
	w    int
	h    int
	data []bool
}

// newMatrix converts an InputMatrix to our own data structure.
func newMatrix(m InputMatrix) *matrix {
	w := m.GetWidth()
	h := m.GetHeight()
	// instantiate our matrix
	mat := &matrix{w: w, h: h, data: make([]bool, w*h)}
	for x := 0; x < mat.w; x++ {
		for y := 0; y < mat.h; y++ {
			// Deep-copy the values
			mat.data[(w*x)+y] = m.Get(x, y)
		}
	}
	return mat
}

// Rect describes a rectangle where (X1, Y2) is the upper-left corner and (X2, Y2) is the lower-right corner
type Rect struct {
	X1, Y1, X2, Y2 int
}

// FindRects takes an InputMatrix and returns a slice of rectangles found in the given input
func FindRects(m InputMatrix) []Rect {
	var rects []Rect

	mat := newMatrix(m)

	// Calculate the total area covered by rectangles.
	totalRectArea := 0
	for x := 0; x < mat.w; x++ {
		for y := 0; y < mat.h; y++ {
			if mat.data[(mat.w*x)+y] {
				totalRectArea += 1
			}
		}
	}

	rectArea := 0
	// While we have not found each rectangle, keep trying. Eventually we'll find every rectangle
	for rectArea < totalRectArea {
		// Find the next rectangle
		rect := findNextRect(mat)
		// Add it to our slice of rectangles to be returned
		rects = append(rects, rect)
		// mark the found rectangle, so we don't find it again in the next iteration.
		markRect(mat, rect)
		// Add the area of the found rectangle to the rectArea
		rectArea += (rect.X2 - rect.X1 + 1) * (rect.Y2 - rect.Y1 + 1)
	}

	return rects
}

// findNextRect finds a new rectangle in the given matrix.
func findNextRect(mat *matrix) Rect {
	foundCorner := false
	r := Rect{X1: 0, X2: mat.w - 1, Y1: 0, Y2: mat.h - 1}
	// Find the upper-left corner
	for x := 0; x < mat.w; x++ {
		for y := 0; y < mat.h; y++ {
			if mat.data[(mat.w*x)+y] {
				r.X1 = x
				r.Y1 = y
				foundCorner = true
				break
			}
		}
		if foundCorner {
			break
		}
	}

	// Find the bottom-right corner
	for x := r.X1; x <= r.X2; x++ {
		if !mat.data[(mat.w*x)+r.Y1] {
			r.X2 = x - 1
			return r
		}
		for y := r.Y1; y <= r.Y2; y++ {
			if !mat.data[(mat.w*x)+y] {
				r.Y2 = y - 1
				break
			}
		}
	}

	return r
}

// markRect removes the given rectangle r from the matrix.
// It does this by setting the values inside the rectangle inside the matrix to false
func markRect(mat *matrix, r Rect) {
	for x := r.X1; x <= r.X2; x++ {
		for y := r.Y1; y <= r.Y2; y++ {
			mat.data[(mat.w*x)+y] = false
		}
	}
}
