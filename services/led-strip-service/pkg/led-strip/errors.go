package led_strip

import "errors"

var (
	ErrIdUnknown = errors.New("id is not known")
	ErrTimeout   = errors.New("timeout on request occurred")
)
