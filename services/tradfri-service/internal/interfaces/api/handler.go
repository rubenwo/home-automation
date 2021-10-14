package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/rubenwo/home-automation/services/tradfri-service/internal/entity"
	"github.com/rubenwo/home-automation/services/tradfri-service/internal/interfaces/api/scheme"
	"github.com/rubenwo/home-automation/services/tradfri-service/internal/usecases"
	"net/http"
	"strings"
)

func RegisterHandler(usecases *usecases.TradfriUsecases, router chi.Router) {
	handler := Handler{usecases: usecases}

	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			next.ServeHTTP(w, r)
		})
	})

	router.Get("/tradfri/devices", handler.getTradfriDevices)
	router.Get("/tradfri/devices/{deviceId}", handler.getTradfriDevice)
	router.Post("/tradfri/devices/command", handler.postDevicesCommand)
	router.Post("/tradfri/devices/{deviceId}/command", handler.postDeviceCommand)

	router.Get("/tradfri/groups", handler.getTradfriGroups)
	router.Get("/tradfri/groups/{groupId}", handler.getTradfriGroup)
	router.Post("/tradfri/groups/command", handler.postCommandInAllGroups)
	router.Post("/tradfri/groups/{groupId}/command", handler.postCommandInGroup)
}

type Handler struct {
	usecases *usecases.TradfriUsecases
}

func (h *Handler) getTradfriDevices(w http.ResponseWriter, r *http.Request) {
	devices, err := h.usecases.FetchAllDevices()
	if err != nil {
		errorController(w, err)
		return
	}

	schemeDevices := make([]scheme.Device, len(devices))
	for i, device := range devices {
		schemeDevices[i] = scheme.Device{
			Id:         device.Id,
			Name:       device.Name,
			Category:   device.Category,
			DeviceType: device.DeviceType,
		}
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(&schemeDevices)
}

func (h *Handler) getTradfriDevice(w http.ResponseWriter, r *http.Request) {
	deviceId := chi.URLParam(r, "deviceId")
	device, err := h.usecases.FetchDevice(deviceId)
	if err != nil {
		errorController(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(&scheme.Device{
		Id:         device.Id,
		Name:       device.Name,
		Category:   device.Category,
		DeviceType: device.DeviceType,
	})
}

func commandSchemeToEntity(cmd scheme.Command) (entity.DeviceCommand, error) {
	switch strings.ToLower(cmd.DeviceType) {
	case "light":
		if cmd.DimmableLightCommand == nil {
			return entity.DeviceCommand{}, errors.New("device type is 'light', but dimmable light command is nil")
		}
		return entity.DeviceCommand{
			DeviceType: entity.LIGHT,
			DimmableLightCommand: &entity.DimmableLightCommand{
				Power:      cmd.DimmableLightCommand.Power,
				Brightness: cmd.DimmableLightCommand.Brightness,
			},
		}, nil
	}

	return entity.DeviceCommand{}, nil
}

func (h *Handler) postDevicesCommand(w http.ResponseWriter, r *http.Request) {
	devices, err := h.usecases.FetchAllDevices()
	if err != nil {
		errorController(w, err)
		return
	}
	var command scheme.Command
	if err := json.NewDecoder(r.Body).Decode(&command); err != nil {
		errorController(w, err)
		return
	}

	entityCommand, err := commandSchemeToEntity(command)
	if err != nil {
		errorController(w, err)
		return
	}

	for _, device := range devices {
		err := h.usecases.CommandDevice(device.Id, entityCommand)
		if err != nil {
			errorController(w, err)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) postDeviceCommand(w http.ResponseWriter, r *http.Request) {
	deviceId := chi.URLParam(r, "deviceId")
	var command scheme.Command
	if err := json.NewDecoder(r.Body).Decode(&command); err != nil {
		errorController(w, err)
		return
	}

	entityCommand, err := commandSchemeToEntity(command)
	if err != nil {
		errorController(w, err)
		return
	}

	if err := h.usecases.CommandDevice(deviceId, entityCommand); err != nil {
		errorController(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) getTradfriGroups(w http.ResponseWriter, r *http.Request) {
	groups, err := h.usecases.FetchAllGroups()
	if err != nil {
		errorController(w, err)
		return
	}
	fmt.Println(groups)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) getTradfriGroup(w http.ResponseWriter, r *http.Request) {
	groupId := chi.URLParam(r, "groupId")
	group, err := h.usecases.FetchGroup(groupId)
	if err != nil {
		errorController(w, err)
		return
	}
	fmt.Println(group)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) postCommandInAllGroups(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) postCommandInGroup(w http.ResponseWriter, r *http.Request) {
	groupId := chi.URLParam(r, "groupId")
	fmt.Println(groupId)
	w.WriteHeader(http.StatusOK)
}
