package entity

type DeviceType string

const (
	Light  DeviceType = "light"
	Remote DeviceType = "remote"
)

type TradfriDevice struct {
	Id                 string
	Name               string
	Category           string
	DeviceType         DeviceType
	DimmableDeviceData *DimmableDeviceData
	RemoteData         *RemoteData
}

type DimmableDeviceData struct {
	Power      int // 0 or 1
	Brightness int // 0 to 255
}

type RemoteData struct {
	BatteryLevel int
}

type DeviceCommand struct {
	DeviceType           DeviceType
	DimmableLightCommand *DimmableLightCommand
}

type DimmableLightCommand struct {
	Power      *int // 0 or 1
	Brightness *int // 0 to 255
}
