package scheme

type Device struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Category   string `json:"category"`
	DeviceType string `json:"device_type"`

	DimmableDeviceData *DimmableDeviceData `json:"dimmable_device_data"`
	RemoteData         *RemoteData         `json:"remote_data"`
}

type DimmableDeviceData struct {
	Power      int `json:"power"`      // 0 or 1
	Brightness int `json:"brightness"` // 0 to 255
}

type RemoteData struct {
	BatteryLevel int `json:"battery_level"`
}
