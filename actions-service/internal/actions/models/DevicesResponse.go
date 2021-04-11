package models

type DevicesResponse struct {
	Devices []struct {
		Category string `json:"category"`
		Id       string `json:"id"`
		Name     string `json:"name"`
		Product  struct {
			Company string `json:"company"`
			Type    string `json:"type"`
		} `json:"product"`
	} `json:"devices"`
}
