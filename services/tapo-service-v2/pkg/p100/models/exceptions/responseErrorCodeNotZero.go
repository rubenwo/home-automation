package exceptions

import "fmt"

type ResponseErrorCodeNotZero struct {
	responseCode int
}

func (r *ResponseErrorCodeNotZero) Error() string {
	return fmt.Sprintf("ResponseErrorCodeNotZero: %d", r.responseCode)
}
