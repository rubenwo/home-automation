package led_strip

// TODO: Change this information. Also do this in the C++ code
type Information struct {
	DeviceName   string `json:"device_name"`
	DeviceId     string `json:"device_id"`
	IsHealthy    bool   `json:"is_healthy"`
	ErrorMessage string `json:"error_message"`
}
