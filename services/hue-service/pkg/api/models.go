package api

import (
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/rubenwo/home-automation/services/hue-service/pkg/hue"
)

func createSchema(db *pg.DB) error {
	tables := []interface{}{
		(*hue.Bridge)(nil),
	}

	for _, model := range tables {
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

type JsonError struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
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
