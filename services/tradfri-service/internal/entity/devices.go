package entity

type TradfriDevice struct {
	Id         string
	Name       string
	Category   string
	DeviceType string
}

type DeviceType string

const (
	LIGHT DeviceType = "light"
)

type DeviceCommand struct {
	DeviceType           DeviceType
	DimmableLightCommand *DimmableLightCommand
}

type DimmableLightCommand struct {
	Power      int // 0 or 1
	Brightness int // 0 to 255
}
