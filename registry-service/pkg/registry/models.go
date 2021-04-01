package registry

import (
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

// createSchema creates database schema for the models
func createSchema(db *pg.DB) error {
	models := []interface{}{
		(*DeviceInfo)(nil),
		(*Product)(nil),
		(*SensorDevice)(nil),
	}

	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

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

type SensorDevice struct {
	Id             int64       `json:"id"`
	Name           string      `json:"name"`
	SensorType     string      `json:"sensor_type"`
	ConnectionData interface{} `json:"connection_data"`
}
