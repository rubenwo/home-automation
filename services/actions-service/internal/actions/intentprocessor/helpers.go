package intentprocessor

import (
	"encoding/json"
	"fmt"
	"github.com/rubenwo/home-automation/services/actions-service/internal/actions/models"
	"io/ioutil"
	"log"
	"net/http"
)

func getDevices() (models.DevicesResponse, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "http://registry.default.svc.cluster.local/devices", nil)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	var deviceRespone models.DevicesResponse
	if err := json.NewDecoder(resp.Body).Decode(&deviceRespone); err != nil {
		log.Fatal(err)
	}
	if err := resp.Body.Close(); err != nil {
		return models.DevicesResponse{}, err
	}

	return deviceRespone, nil
}

func commandTapoDevice(deviceId string, command string, brightness int) error {
	client := &http.Client{}

	url := fmt.Sprintf("http://tapo.default.svc.cluster.local/tapo/lights/%s?command=%s&brightness=%d",
		deviceId, command, brightness)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	_, _ = ioutil.ReadAll(resp.Body)
	if err := resp.Body.Close(); err != nil {
		return err
	}
	return nil
}
