package models

type Routine struct {
	Id      int64    `json:"id"`
	Name    string   `json:"name"`
	Trigger Trigger  `json:"trigger"`
	Actions []Action `json:"actions"`
}

type TriggerType uint8

const (
	TimerTriggerType TriggerType = iota
	//OnWebhook called, but that is a problem for another time
)

type Trigger struct {
	Type     TriggerType `json:"type"`
	CronExpr string      `json:"cron_expr,omitempty"`
	Webhook  string      `json:"webhook,omitempty"`
}

type Action struct {
	Addr   string                 `json:"addr"`
	Method string                 `json:"method"`
	Data   map[string]interface{} `json:"data"`
}
