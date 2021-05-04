package models

type Routine struct {
	Id      int64    `json:"id"`
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
	Schedule Schedule    `json:"schedule,omitempty"`
	Webhook  string      `json:"webhook,omitempty"`
}


//Schedule supports from every month at a specific date and time to every minute of every day
type Schedule struct {
	DayOfWeek   DayOfWeek `json:"day_of_week"`
	MonthOfYear string `json:"month_of_year"`

	HourOfDay   string `json:"hour_of_day"`
	MinuteOfDay string `json:"minute_of_day"`
}

type Action struct {
	Addr   string                 `json:"addr"`
	Method string                 `json:"method"`
	Data   map[string]interface{} `json:"data"`
}
