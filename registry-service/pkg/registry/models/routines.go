package models

import "time"

type Routine struct {
	Id      int64    `json:"id"`
	Trigger Trigger  `json:"trigger"`
	Actions []Action `json:"actions"`
}

type Trigger struct {
	Type   string    `json:"type"`
	Repeat bool      `json:"repeat"`
	When   time.Time `json:"when"`
}

type Action struct {
	Addr   string                 `json:"addr"`
	Method string                 `json:"method"`
	Data   map[string]interface{} `json:"data"`
}
