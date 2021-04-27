package registry

import (
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/rubenwo/home-automation/registry-service/pkg/registry/models"
)

// createSchema creates database schema for the models
func createSchema(db *pg.DB) error {
	tables := []interface{}{
		(*models.DeviceInfo)(nil),
		(*models.Product)(nil),
		(*models.SensorDevice)(nil),
		(*models.ConnectionData)(nil),
		(*models.Routine)(nil),
		(*models.Trigger)(nil),
		(*models.Action)(nil),
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
