package usecases

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/rubenwo/home-automation/services/tradfri-service/internal/app"
	"github.com/rubenwo/home-automation/services/tradfri-service/internal/entity"
	"github.com/rubenwo/home-automation/services/tradfri-service/pkg/tradfri"
	"os"
)

type DaoTradfri interface {
	app.TradfriDao
}

type ServicesTradfri interface {
	app.RegistrySyncerService
}

// TODO: move db code to dao/services

func NewTradfriUsecases(db *sql.DB, dao DaoTradfri, services ServicesTradfri) *TradfriUsecases {
	return &TradfriUsecases{
		dao:      dao,
		services: services,
		client: tradfri.NewTradfriClient(
			os.Getenv("TRADFRI_GATEWAY_ADDRESS"),
			os.Getenv("TRADFRI_GATEWAY_CLIENT_ID"),
			os.Getenv("TRADFRI_GATEWAY_PSK"),
		),
		db: db,
	}
}

type TradfriUsecases struct {
	dao      DaoTradfri
	services ServicesTradfri
	client   *tradfri.Client
	db       *sql.DB
}

func (u *TradfriUsecases) tradfriIdToApiId(tradfriId int) (string, error) {
	var id string
	query := `SELECT id FROM ids_tradfriids WHERE tradfri_id = $1`
	if err := u.db.QueryRow(query, tradfriId).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			id = uuid.New().String()
			query = `INSERT INTO ids_tradfriids (id, tradfri_id) VALUES ($1, $2)`
			if _, err := u.db.Exec(query, id, tradfriId); err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	}
	return id, nil
}

func (u *TradfriUsecases) apiIdToTradfriId(id string) (int, error) {
	var tradfriId int
	query := `SELECT tradfri_id FROM ids_tradfriids WHERE id = $1`
	if err := u.db.QueryRow(query, id).Scan(&tradfriId); err != nil {
		return 0, err
	}
	return tradfriId, nil
}

func (u *TradfriUsecases) FetchAllDevices() ([]entity.TradfriDevice, error) {
	tradfriDevices, err := u.client.ListDevices()
	if err != nil {
		return nil, err
	}

	devices := make([]entity.TradfriDevice, len(tradfriDevices))

	for i, device := range tradfriDevices {
		// Map tradfri id (int) to our uuid. This is done to avoid conflicts in the rest of the application
		id, err := u.tradfriIdToApiId(tradfriDevices[i].DeviceId)
		if err != nil {
			return []entity.TradfriDevice{}, err
		}

		var (
			remoteData         *entity.RemoteData         = nil
			dimmableDeviceData *entity.DimmableDeviceData = nil
		)

		var deviceType entity.DeviceType
		switch device.Type {
		case 0:
			deviceType = entity.Remote
			remoteData = &entity.RemoteData{BatteryLevel: device.Metadata.Battery}
		case 2:
			deviceType = entity.Light
			dimmableDeviceData = &entity.DimmableDeviceData{
				Power:      device.LightControl[0].Power,
				Brightness: device.LightControl[0].Dimmer,
			}
		}

		devices[i] = entity.TradfriDevice{
			Id:                 id,
			Name:               tradfriDevices[i].Name,
			Category:           fmt.Sprintf("%d", tradfriDevices[i].Type),
			DeviceType:         deviceType,
			RemoteData:         remoteData,
			DimmableDeviceData: dimmableDeviceData,
		}
	}

	return devices, nil
}

func (u *TradfriUsecases) FetchDevice(deviceId string) (entity.TradfriDevice, error) {
	tradfriId, err := u.apiIdToTradfriId(deviceId)
	if err != nil {
		return entity.TradfriDevice{}, err
	}

	tradfriDevice, err := u.client.GetDevice(tradfriId)
	if err != nil {
		return entity.TradfriDevice{}, err
	}

	var (
		remoteData         *entity.RemoteData         = nil
		dimmableDeviceData *entity.DimmableDeviceData = nil
	)

	var deviceType entity.DeviceType
	switch tradfriDevice.Type {
	case 0:
		deviceType = entity.Remote
		remoteData = &entity.RemoteData{BatteryLevel: tradfriDevice.Metadata.Battery}
	case 2:
		deviceType = entity.Light
		dimmableDeviceData = &entity.DimmableDeviceData{
			Power:      tradfriDevice.LightControl[0].Power,
			Brightness: tradfriDevice.LightControl[0].Dimmer,
		}
	}
	return entity.TradfriDevice{
		Id:                 deviceId,
		Name:               tradfriDevice.Name,
		Category:           fmt.Sprintf("%d", tradfriDevice.Type),
		DeviceType:         deviceType,
		RemoteData:         remoteData,
		DimmableDeviceData: dimmableDeviceData,
	}, nil
}

func (u *TradfriUsecases) CommandDevice(deviceId string, command entity.DeviceCommand) error {
	tradfriId, err := u.apiIdToTradfriId(deviceId)
	if err != nil {
		return err
	}
	switch command.DeviceType {
	case entity.Light:
		if command.DimmableLightCommand == nil {
			return errors.New("device type is 'Light', but command is nil")
		}
		if command.DimmableLightCommand.Power != nil {
			_, err = u.client.PutDevicePower(tradfriId, *command.DimmableLightCommand.Power)
			if err != nil {
				return err
			}
		}
		if command.DimmableLightCommand.Brightness != nil {
			_, err = u.client.PutDeviceDimming(tradfriId, *command.DimmableLightCommand.Brightness)
			if err != nil {
				return err
			}
		}
	default:
		return errors.New("unsupported command type")
	}

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
