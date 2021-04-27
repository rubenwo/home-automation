package models

type SensorDevice struct {
	Id             int64          `json:"id"`
	Name           string         `json:"name"`
	SensorType     string         `json:"sensor_type"`
	ConnectionData ConnectionData `json:"connection_data"`
}

type ConnectionData struct {
	IpAddr string `json:"ip_addr"`
}
