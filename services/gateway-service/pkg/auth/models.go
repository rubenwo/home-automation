package auth

import (
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/rubenwo/home-automation/services/gateway-service/pkg/auth/models"
)

// createSchema creates database schema for the models
func createSchema(db *pg.DB) error {
	tables := []interface{}{
		(*models.Claims)(nil),
		(*models.Authorization)(nil),
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
