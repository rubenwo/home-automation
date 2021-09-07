package service

import (
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

// createSchema creates database schema for the models
func createSchema(db *pg.DB) error {
	tables := []interface{}{
		(*NotificationModel)(nil),
		(*NotificationSubscriber)(nil),
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

type NotificationModel struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type NotificationSubscriber struct {
	ClientPlatform         string `json:"client_platform"`
	FirebaseMessagingToken string `json:"firebase_messaging_token"`
}

type JsonError struct {
	Code         int    `json:"code"`
	ErrorMessage string `json:"error_message"`
}
