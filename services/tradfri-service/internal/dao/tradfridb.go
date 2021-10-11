package dao

import "database/sql"

func NewTradfriDB(db *sql.DB) *tradfriDB {
	return &tradfriDB{db: db}
}

type tradfriDB struct {
	db *sql.DB
}
