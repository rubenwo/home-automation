package usecases

import (
	"github.com/rubenwo/home-automation/services/tradfri-service/internal/app"
	"github.com/rubenwo/home-automation/services/tradfri-service/internal/entity"
)

type DaoTradfri interface {
	app.TradfriDao
}

type ServicesTradfri interface {
}

func NewTradfriUsecases(dao DaoTradfri, services ServicesTradfri) *TradfriUsecases {
	return &TradfriUsecases{dao: dao, services: services}
}

type TradfriUsecases struct {
	dao      DaoTradfri
	services ServicesTradfri
}

func (u *TradfriUsecases) FetchAllDevices() ([]entity.TradfriDevice, error) {
	return []entity.TradfriDevice{}, nil
}

func (u *TradfriUsecases) FetchDevice(deviceId string) (entity.TradfriDevice, error) {
	return entity.TradfriDevice{}, nil
}

func (u *TradfriUsecases) CommandDevice(deviceId string) error {
	return nil
}

func (u *TradfriUsecases) FetchAllGroups() ([]entity.TradfriGroup, error) {
	return []entity.TradfriGroup{}, nil
}

func (u *TradfriUsecases) FetchGroup(groupId string) (entity.TradfriGroup, error) {
	return entity.TradfriGroup{}, nil
}

func (u *TradfriUsecases) CommandGroup(groupId string) error {
	return nil
}
