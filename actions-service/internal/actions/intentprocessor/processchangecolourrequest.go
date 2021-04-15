package intentprocessor

import (
	"fmt"
	"github.com/rubenwo/home-automation/actions-service/internal/actions/models"
)

type ProcessChangeColourRequest struct{}

func (p *ProcessChangeColourRequest) ProcessIntent(params map[string]string) (string, error) {

	deviceRespone, err := getDevices()
	if err != nil {
		return "", err
	}

	requestedDevice := params["device"]
	requestColour := params["colour"]

	device := models.Device{}

	for _, ddevice := range deviceRespone.Devices {
		if ddevice.Name == requestedDevice {
			device = ddevice
			break
		}
	}
	if device.Id == "" {
		return "", fmt.Errorf("device: %s not found in registry", requestedDevice)
	}

	fmt.Println(device)
	//
	//client := &http.Client{}
	//
	//
	//
	//req, err := http.NewRequest("POST", "led-strip.default.svc.cluster.local/leds/devices/{id}/command", bytes.NewBuffer(b))

	return fmt.Sprintf("successfully changed colour to: %s for device: %s", requestColour, requestedDevice), nil
}
