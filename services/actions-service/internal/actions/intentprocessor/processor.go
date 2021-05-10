package intentprocessor

type IntentProcessor interface {
	ProcessIntent(params map[string]string) (string, error)
}
