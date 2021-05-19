package methods

type Method interface {
	Method() string
	Params() interface{}
}
