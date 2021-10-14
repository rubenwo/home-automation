package scheme

type Command struct {
	DeviceType           string                `json:"device_type"`
	DimmableLightCommand *DimmableLightCommand `json:"dimmable_light_command"`
}

type DimmableLightCommand struct {
	Power      int `json:"power"`
	Brightness int `json:"brightness"`
}
