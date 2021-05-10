package models


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
