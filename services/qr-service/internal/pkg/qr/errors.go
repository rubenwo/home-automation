package qr

import "fmt"

type ErrFailedToCreateQRCode struct {
	Err error
}

func (e *ErrFailedToCreateQRCode) Error() string {
	return fmt.Sprintf("could not create qr-code: %s", e.Err.Error())
}

type ErrEncodingPNG struct {
	Err error
}

func (e *ErrEncodingPNG) Error() string {
	return fmt.Sprintf("could not encoding to PNG: %s", e.Err.Error())
}

type ErrEncodingJPEG struct {
	Err error
}

func (e *ErrEncodingJPEG) Error() string {
	return fmt.Sprintf("could not encoding to JPEG: %s", e.Err.Error())
}

type ErrEncodingSVG struct {
	Err error
}

func (e *ErrEncodingSVG) Error() string {
	return fmt.Sprintf("could not encoding to SVG: %s", e.Err.Error())
}
