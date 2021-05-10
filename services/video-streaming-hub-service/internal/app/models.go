package app

import (
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

// createSchema creates database schema for the Item, ... models
func createSchema(db *pg.DB) error {
	models := []interface{}{
		(*Camera)(nil),
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

type Camera struct {
	Id   int64  `json:"id"`
	Host string `json:"host"`
}
