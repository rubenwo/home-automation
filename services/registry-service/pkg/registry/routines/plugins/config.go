package plugins

import "errors"

var (
	ErrConfigCastingFailed = errors.New("wrong config provided for the plugin")
)

type Config interface {
	Assert() error
}
