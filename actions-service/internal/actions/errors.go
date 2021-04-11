package actions

import "errors"

var (
	ConfigNilError = errors.New("config is nil")
	ConfigValidationError = errors.New("error validating config")
)
