package models

type DeviceConnectionInfo struct {
	DeviceId string `json:"device_id"` // The same ID as the TapoDevice

	DeviceType string `json:"device_type"`
	IpAddress  string `json:"ip_address"`
	Email      string `json:"email"`
	Password   string `json:"password"`
}

type TapoDevice struct {
	DeviceId   string                 `json:"device_id"`
	DeviceName string                 `json:"device_name"`
	DeviceType string                 `json:"device_type"`
	DeviceInfo map[string]interface{} `json:"device_info"`
}
