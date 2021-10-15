package scheme

type Device struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Category   string `json:"category"`
	DeviceType string `json:"device_type"`

	DimmableDeviceData *DimmableDeviceData `json:"dimmable_device_data,omitempty"`
	RemoteData         *RemoteData         `json:"remote_data,omitempty"`
}

type DimmableDeviceData struct {
	Power      *int `json:"power,omitempty"`      // 0 or 1
	Brightness *int `json:"brightness,omitempty"` // 0 to 255
}

type RemoteData struct {
	BatteryLevel int `json:"battery_level"`
}
