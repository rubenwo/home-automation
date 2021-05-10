package models

type WebhookRequest struct {
	ResponseId                  string      `json:"responseId"`
	Session                     string      `json:"session"`
	QueryResult                 QueryResult `json:"queryResult"`
	OriginalDetectIntentRequest interface{} `json:"originalDetectIntentRequest"`
}

type QueryResult struct {
	QueryText                 string               `json:"queryText"`
	Parameters                map[string]string    `json:"parameters"`
	AllRequiredParamsPresent  bool                 `json:"allRequiredParamsPresent"`
	FulfillmentText           string               `json:"fulfillmentText"`
	FulfillmentMessages       []FulfillmentMessage `json:"fulfillmentMessages"`
	OutputContexts            []OutputContext      `json:"outputContexts"`
	Intent                    Intent               `json:"intent"`
	IntentDetectionConfidence float32              `json:"intentDetectionConfidence"`
	DiagnosticInfo            interface{}          `json:"diagnosticInfo"`
	LanguageCode              string               `json:"languageCode"`
}

type OutputContext struct {
	Name          string `json:"name"`
	LifespanCount int    `json:"lifespanCount"`
	//Parameters    map[string]string `json:"parameters"`
}
type Intent struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
}

type FulfillmentMessage struct {
	Text struct {
		Text []string `json:"text"`
	} `json:"text"`
}
