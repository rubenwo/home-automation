package usecases

import (
	"fmt"
	"github.com/eriklupander/tradfri-go/tradfri"
	"github.com/rubenwo/home-automation/services/tradfri-service/internal/app"
	"github.com/rubenwo/home-automation/services/tradfri-service/internal/entity"
	"os"
	"strconv"
)

type DaoTradfri interface {
	app.TradfriDao
}

type ServicesTradfri interface {
	app.RegistrySyncerService
}

func NewTradfriUsecases(dao DaoTradfri, services ServicesTradfri) *TradfriUsecases {
	return &TradfriUsecases{
		dao:      dao,
		services: services,
		client: tradfri.NewTradfriClient(
			os.Getenv("TRADFRI_GATEWAY_ADDRESS"),
			os.Getenv("TRADFRI_GATEWAY_CLIENT_ID"),
			os.Getenv("TRADFRI_GATEWAY_PSK"),
		),
	}
}

type TradfriUsecases struct {
	dao      DaoTradfri
	services ServicesTradfri
	client   *tradfri.Client
}

func (u *TradfriUsecases) FetchAllDevices() ([]entity.TradfriDevice, error) {
	tradfriDevices, err := u.client.ListDevices()
	if err != nil {
		return nil, err
	}

	devices := make([]entity.TradfriDevice, len(tradfriDevices))

	for i := range tradfriDevices {
		devices[i] = entity.TradfriDevice{
			Id:         fmt.Sprintf("%d", tradfriDevices[i].DeviceId),
			Name:       tradfriDevices[i].Name,
			Category:   fmt.Sprintf("%d", tradfriDevices[i].Type),
			DeviceType: fmt.Sprintf("%d", tradfriDevices[i].Type),
		}
	}

	return devices, nil
}

func (u *TradfriUsecases) FetchDevice(deviceId string) (entity.TradfriDevice, error) {
	id, err := strconv.Atoi(deviceId)
	if err != nil {
		return entity.TradfriDevice{}, err
	}
	tradfriDevice, err := u.client.GetDevice(id)
	if err != nil {
		return entity.TradfriDevice{}, err
	}

	return entity.TradfriDevice{
		Id:         fmt.Sprintf("%d", tradfriDevice.DeviceId),
		Name:       tradfriDevice.Name,
		Category:   fmt.Sprintf("%d", tradfriDevice.Type),
		DeviceType: fmt.Sprintf("%d", tradfriDevice.Type),
	}, nil
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
