package models

type TextResponse struct {
	FulfillmentMessages []Message `json:"fulfillmentMessages"`
}

type Message struct {
	Text Text `json:"text"`
}
type Text struct {
	Text []string `json:"text"`
}
type CardResponse struct {
	FulfillmentMessages []struct {
	} `json:"fulfillmentMessages"`
}

type GoogleAssistantResponse struct {
}
