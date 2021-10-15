package scheme

type Command struct {
	DeviceType           string                `json:"device_type"`
	DimmableLightCommand *DimmableLightCommand `json:"dimmable_light_command,omitempty"`
}

type DimmableLightCommand struct {
	Power      *int `json:"power,omitempty"`
	Brightness *int `json:"brightness,omitempty"`
}
