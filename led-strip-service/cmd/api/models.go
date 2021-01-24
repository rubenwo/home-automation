package main

type HealthzModel struct {
	IsHealthy    bool   `json:"is_healthy"`
	ErrorMessage string `json:"error_message"`
}

type RegisterLedDeviceModel struct {
	IPAddress string `json:"ip_address"`
}

type LedDeviceModel struct {
	ID             string      `json:"id"`
	Name           string      `json:"name"`
	NumLeds        int         `json:"num_leds"`
	SupportedModes []string    `json:"supported_modes"`
	CurrentMode    string      `json:"current_mode"`
	IPAddress      string      `json:"ip_address"`
	Data           interface{} `json:"data"`
}

type JsonError struct {
	Code         int    `json:"code"`
	ErrorMessage string `json:"error_message"`
}

type LedControllerInfo struct {
	DeviceName string `json:"device_name"`
	DeviceType string `json:"device_type"`
	DeviceInfo struct {
		SupportedModes []string    `json:"supported_modes"`
		CurrentMode    string      `json:"current_mode"`
		Data           interface{} `json:"data"`
	}
}

type CommandResponseModel struct {
	Message string `json:"message"`
}
