package app

import "github.com/rubenwo/home-automation/services/tradfri-service/internal/entity"

type RegistrySyncerService interface {
	PublishDevice(device entity.TradfriDevice) error
}

type Services struct {
	RegistrySyncerService
}
