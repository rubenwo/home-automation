package registry

type HealthzModel struct {
	IsHealthy    bool   `json:"is_healthy"`
	ErrorMessage string `json:"error_message"`
}

type DeviceInfo struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Product  Product `json:"product"`
}
type Product struct {
	Company string `json:"company"`
	Type    string `json:"type"`
}
