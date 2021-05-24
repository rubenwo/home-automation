package models

type NewDevice struct {
	Id         int64  `json:"id"`
	IpAddress  string `json:"ip_address"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	DeviceType string `json:"device_type"`
}

type DatabaseDevice struct {
	Id        int64  `json:"id"`
	IpAddress string `json:"ip_address"`
	DeviceId  string `json:"device_id"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type Device struct {
	Id         int64                  `json:"id"`
	DeviceId   string                 `json:"device_id"`
	DeviceName string                 `json:"device_name"`
	DeviceType string                 `json:"device_type"`
	DeviceInfo map[string]interface{} `json:"device_info"`
}
