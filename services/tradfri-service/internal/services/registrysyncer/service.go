package registrysyncer

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/rubenwo/home-automation/services/tradfri-service/internal/entity"
	"io/ioutil"
	"net/http"
)

type message struct {
	device entity.TradfriDevice
	err    chan error
}

type service struct {
	httpClient *http.Client

	buffer chan message

	db *sql.DB
}

func NewService(db *sql.DB) *service {
	return &service{
		httpClient:     &http.Client{},
		buffer:         make(chan message, 100),
		db:             db,
	}
}

func (s *service) PublishDevice(device entity.TradfriDevice) error {
	errChan := make(chan error, 1)
	s.buffer <- message{
		device: device,
		err:    errChan,
	}
	return <-errChan
}

func (s *service) Run(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-s.buffer:
				if err := s.publishDevice(msg.device); err != nil {
					msg.err <- err
				}
				close(msg.err)
			}
		}
	}()
}

func (s *service) publishDevice(device entity.TradfriDevice) error {
	var id string
	query := `SELECT id FROM ids_tradfriids WHERE tradfri_id = $1`
	if err := s.db.QueryRow(query, device.Id).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			id = uuid.New().String()
			query = `INSERT INTO ids_tradfriids (id, tradfri_id) VALUES ($1, $2)`
			if _, err := s.db.Exec(query, id, device.Id); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	req, err := http.NewRequest("DELETE", "http://registry.default.svc.cluster.local/devices/"+id, nil)
	if err != nil {
		return err
	}
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err := resp.Body.Close(); err != nil {
		return err
	}

	var data struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Category string `json:"category"`
		Product  struct {
			Company string `json:"company"`
			Type    string `json:"type"`
		} `json:"product"`
	}
	data.ID = id
	data.Name = device.Name
	data.Category = device.Category
	data.Product.Company = "IKEA"
	data.Product.Type = device.DeviceType
	jsonData, err := json.Marshal(&data)
	if err != nil {
		return err
	}

	req, err = http.NewRequest("POST", "http://registry.default.svc.cluster.local/devices", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	resp, err = s.httpClient.Do(req)
	if err != nil {
		return err
	}
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err := resp.Body.Close(); err != nil {
		return err
	}
	return nil
}
