package app

type TradfriDao interface {
}

// DataAccessObject is used by usecases (and services) to load and save entities.
type DataAccessObject struct {
	TradfriDao
}
