package api

type HealthzModel struct {
	IsHealthy    bool   `json:"is_healthy"`
	ErrorMessage string `json:"error_message"`
}

type JsonError struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

type DeviceInfo struct {
	DeviceID int    `json:"device_id"`
	Status   string `json:"status"`
}
