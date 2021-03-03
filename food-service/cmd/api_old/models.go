package main

type Recipe struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Img         string       `json:"img"`
	Ingredients []Ingredient `json:"ingredients"`
	Steps       []Step       `json:"steps"`
}

type Ingredient struct {
	Name   string `json:"name"`
	Amount string `json:"amount"`
}

type Step struct {
	Instruction string `json:"instruction"`
}
