package models

type Routine struct {
	Id       int64    `json:"id"`
	Name     string   `json:"name"`
	IsActive bool     `json:"is_active"`
	Trigger  Trigger  `json:"trigger"`
	Actions  []Action `json:"actions"`
}

type TriggerType uint8

const (
	TimerTriggerType TriggerType = iota
	MqttEventTriggerType
)

type Trigger struct {
	Type     TriggerType `json:"type"`
	CronExpr string      `json:"cron_expr,omitempty"`
	Webhook  string      `json:"webhook,omitempty"`
}

type Action struct {
	Script string `json:"script"`

	Addr   string                 `json:"addr"`
	Method string                 `json:"method"`
	Data   map[string]interface{} `json:"data"`
}
