package models

import "time"

type RoutineLog struct {
	RoutineId int64     `json:"routine_id"`
	LoggedAt  time.Time `pg:"default:now()" json:"logged_at"`
	Message   string    `json:"message"`
}
