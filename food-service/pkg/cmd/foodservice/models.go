package foodservice

import (
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

// createSchema creates database schema for the Item, ... models
func createSchema(db *pg.DB) error {
	models := []interface{}{
		(*Recipe)(nil),
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

type Recipe struct {
	Id          int64        `json:"id"`
	Name        string       `json:"name"`
	Img         string       `json:"img"`
	Ingredients []Ingredient `json:"ingredients",pg:"rel:has-many"`
	Steps       []Step       `json:"steps",pg:"rel:has-many"`
}

type Ingredient struct {
	Id     int64  `json:"id"`
	Name   string `json:"name"`
	Amount string `json:"amount"`
}

type Step struct {
	Id          int64  `json:"id"`
	Instruction string `json:"instruction"`
}
