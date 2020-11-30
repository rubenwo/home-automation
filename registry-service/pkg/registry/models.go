package registry

type HealthzModel struct {
	IsHealthy    bool   `json:"is_healthy"`
	ErrorMessage string `json:"error_message"`
}

type DeviceInfo struct {
	Name          string `json:"name"`
	DeviceType    string `json:"device_type"`
	DeviceCompany string `json:"device_company"`
}
