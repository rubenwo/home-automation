package service

import (
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

// createSchema creates database schema for the models
func createSchema(db *pg.DB) error {
	tables := []interface{}{
		(*LedDeviceModel)(nil),
		(*HealthzModel)(nil),
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

type LedDeviceModel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	//NumLeds        int      `json:"num_leds"`
	//SupportedModes []string `json:"supported_modes"`
	//CurrentMode    string   `json:"current_mode"`
	//IPAddress      string   `json:"ip_address"`
	//Data           interface{} `json:"data"`
}

type JsonError struct {
	Code         int    `json:"code"`
	ErrorMessage string `json:"error_message"`
}
