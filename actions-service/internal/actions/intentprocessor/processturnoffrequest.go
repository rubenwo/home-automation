package intentprocessor

import (
	"fmt"
	"strings"
)

type ProcessTurnOffRequest struct{}

func (p *ProcessTurnOffRequest) ProcessIntent(params map[string]string) (string, error) {
	deviceRespone, err := getDevices()
	if err != nil {
		return "", err
	}

	requestedDevice := params["device"]

	if requestedDevice == "all" {
		var errors []error

		for _, device := range deviceRespone.Devices {
			switch device.Product.Company {
			case "tp-link":
				if err := commandTapoDevice(device.Id, "off", 0); err != nil {
					errors = append(errors, err)
				}
			}
		}

		if errors != nil {
			return "", fmt.Errorf("%v", errors)
		}
		return fmt.Sprint("turned on all supported devices successfully"), nil
	}

	var (
		deviceId      = ""
		deviceCompany = ""
	)

	for _, device := range deviceRespone.Devices {
		if strings.ToLower(requestedDevice) == strings.ToLower(device.Name) {
			deviceId = device.Id
			deviceCompany = device.Product.Company
			break
		}
	}

	if deviceId == "" {
		return "", fmt.Errorf("device %s not found", requestedDevice)
	}
	switch strings.ToLower(deviceCompany) {
	case "tp-link":
		if err := commandTapoDevice(deviceId, "off", 0); err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf("turn off command for company: %s is not supported", deviceCompany)
	}

	return fmt.Sprintf("turned off device %s successfully", requestedDevice), nil
}