package led_strip

const (
	SolidColorMode     = "SINGLE_COLOR_RGB"
	GradientColorMode  = "GRADIENT_RGB"
	AnimationColorMode = "ANIMATION_RGB"
)

type AnnouncementMessage struct {
	DeviceName   string `json:"device_name"`
	DeviceId     string `json:"device_id"`
	IsHealthy    bool   `json:"is_healthy"`
	ErrorMessage string `json:"error_message"`
}

type AnimationMessage struct {
	Mode           string  `json:"mode"`
	AnimationSpeed int     `json:"animation_speed"`
	Config         []Color `json:"config"`
}

type SolidColorMessage struct {
	Mode string `json:"mode"`

	R int `json:"r"`
	G int `json:"g"`
	B int `json:"b"`
}
