package led_strip

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"github.com/rubenwo/home-automation/services/led-strip-service/pkg/colorconv"
	"log"
	"sync"
	"time"
)

const (
	AnnouncementMqttPath = "leds/announcement"

	QosAtMostOnce  = 0
	QosAtLeastOnce = 1
	QosExactlyOnce = 2
)

type Client struct {
	mqttClient mqtt.Client

	// TODO: add internal cooldown between requests to the same id
	connectedLedStrips map[string]string

	onLedStripOnline func(string, string)
}

func NewClient(host string, retry int, storedDevices map[string]string) (*Client, error) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("%s:1883", host))
	opts.SetClientID(uuid.New().String())
	client := mqtt.NewClient(opts)

	var err error
	for i := 0; i < retry; i++ {
		token := client.Connect()
		token.Wait()
		err = token.Error()
		if err == nil {
			break
		}
		time.Sleep(time.Second * 1)
	}
	if err != nil {
		return nil, err
	}

	ledStripClient := &Client{mqttClient: client, connectedLedStrips: map[string]string{}}
	if storedDevices != nil {
		ledStripClient.connectedLedStrips = storedDevices
	}

	client.Subscribe(AnnouncementMqttPath, QosAtLeastOnce, func(cl mqtt.Client, msg mqtt.Message) {
		msg.Ack()
		payload := msg.Payload()
		var am AnnouncementMessage
		if err := json.Unmarshal(payload, &am); err != nil {
			log.Fatal(err)
		}
		if !am.IsHealthy {
			log.Printf("received message from %s - %s saying it's not healthy\n", am.DeviceName, am.DeviceId)
			return
		}
		ledStripClient.connectedLedStrips[am.DeviceId] = am.DeviceName

		if ledStripClient.onLedStripOnline != nil {
			go ledStripClient.onLedStripOnline(am.DeviceId, am.DeviceName)
		}
	})

	return ledStripClient, nil
}

// ConnectLedStrips returns the currently known connected led strips
func (c *Client) ConnectLedStrips() map[string]string {
	return c.connectedLedStrips
}

func (c *Client) SetOnLedStripOnlineCallback(f func(string, string)) {
	c.onLedStripOnline = f
}

// InformationById returns the information of a led strip with the given id
// returns an error in case of failure
func (c *Client) InformationById(id string) (Information, error) {
	// TODO: Change this when the C++ code is updated
	if _, exists := c.connectedLedStrips[id]; !exists {
		return Information{}, ErrIdUnknown
	}

	requestPath := fmt.Sprintf("leds/%s/information", id)
	responsePath := fmt.Sprintf("leds/%s/response", id)

	var im Information
	responseChan := make(chan bool, 1)

	c.mqttClient.Subscribe(responsePath, QosAtLeastOnce, func(cl mqtt.Client, msg mqtt.Message) {
		msg.Ack()
		payload := msg.Payload()
		if err := json.Unmarshal(payload, &im); err != nil {
			log.Fatal(err)
		}
		responseChan <- true
		close(responseChan)
	})

	defer c.mqttClient.Unsubscribe(responsePath)

	token := c.mqttClient.Publish(requestPath, QosAtLeastOnce, false, []byte(""))
	token.Wait()

	select {
	case <-responseChan:
		return im, nil
	case <-time.After(time.Second * 10):
		return Information{}, ErrTimeout
	}
}

// SetSolidColor sets the color for all known led strips
// returns an error in case of failure
func (c *Client) SetSolidColor(col Color) error {
	var errs error
	var wg sync.WaitGroup
	for id, _ := range c.connectedLedStrips {
		wg.Add(1)
		go func(id string) {
			err := c.SetSolidColorById(col, id)
			if err != nil {
				errs = err
			}
			wg.Done()
		}(id)
	}

	wg.Wait()

	return errs
}

// SetSolidColorById sets the color of a single led strip with the given id
// returns an error in case of failure
func (c *Client) SetSolidColorById(col Color, id string) error {
	msg := SolidColorMessage{
		Mode: SolidColorMode,
		R:    col.R,
		G:    col.G,
		B:    col.B,
	}
	data, err := json.Marshal(&msg)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("leds/%s/control", id)
	token := c.mqttClient.Publish(path, QosAtLeastOnce, false, data)
	token.Wait()

	return nil
}

// SetAnimationBreathing sets the mode of all known led strips to the animation breathing cycle,
// increasing and decreasing the brightness of the provided color
// returns an error in case of failure
func (c *Client) SetAnimationBreathing(col Color) error {
	var errs error
	var wg sync.WaitGroup
	for id, _ := range c.connectedLedStrips {
		wg.Add(1)
		go func(id string) {
			err := c.SetAnimationBreathingById(col, id)
			if err != nil {
				errs = err
			}
			wg.Done()
		}(id)
	}
	wg.Wait()
	return errs
}

// SetAnimationBreathingById sets the mode of the led strip for the provided id to the animation breathing cycle,
// increasing and decreasing the brightness of the provided color
// returns an error in case of failure
func (c *Client) SetAnimationBreathingById(col Color, id string) error {
	msg := AnimationMessage{
		Mode:           AnimationColorMode,
		AnimationSpeed: 10,
		Config:         []Color{},
	}

	for i := 0; i < 255; i++ {
		r := i * col.R / 255
		g := i * col.G / 255
		b := i * col.B / 255
		msg.Config = append(msg.Config, Color{R: r, G: g, B: b})
	}

	for i := 255; i > 0; i-- {
		r := i * col.R / 255
		g := i * col.G / 255
		b := i * col.B / 255
		msg.Config = append(msg.Config, Color{R: r, G: g, B: b})
	}

	data, err := json.Marshal(&msg)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("leds/%s/control", id)
	token := c.mqttClient.Publish(path, QosAtLeastOnce, false, data)
	token.Wait()

	return token.Error()
}

// SetAnimationColorCycle sets the mode of all known led strips to the animation color cycle
// returns an error in case of failure
func (c *Client) SetAnimationColorCycle() error {
	var errs error
	var wg sync.WaitGroup
	for id, _ := range c.connectedLedStrips {
		wg.Add(1)
		go func(id string) {
			err := c.SetAnimationColorCycleById(id)
			if err != nil {
				errs = err
			}
			wg.Done()
		}(id)
	}
	wg.Wait()
	return errs
}

// SetAnimationColorCycleById sets the mode of the led strip for the given id to the animation color cycle
// returns an error in case of failure
func (c *Client) SetAnimationColorCycleById(id string) error {
	msg := AnimationMessage{
		Mode:           AnimationColorMode,
		AnimationSpeed: 10,
		Config:         []Color{},
	}

	for i := 0; i < 360; i++ {
		hsv := colorconv.HSV{
			H: float64(i) / 360.0,
			S: 1.0,
			V: 1.0,
		}

		rgb := hsv.RGB()

		msg.Config = append(msg.Config, Color{R: int(rgb.R), G: int(rgb.G), B: int(rgb.B)})
	}

	data, err := json.Marshal(&msg)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("leds/%s/control", id)
	token := c.mqttClient.Publish(path, QosAtLeastOnce, false, data)
	token.Wait()

	return token.Error()
}

// SetCustomData pushes the given data directly to all known led strips.
// !Important: can mess up the led strips
// returns an error in case of failure
func (c *Client) SetCustomData(data []byte) error {
	var errs error
	var wg sync.WaitGroup
	for id, _ := range c.connectedLedStrips {
		wg.Add(1)
		go func(id string) {
			err := c.SetCustomDataById(data, id)
			if err != nil {
				errs = err
			}
			wg.Done()
		}(id)
	}
	wg.Wait()
	return errs
}

// SetCustomDataById pushes the given data directly to the led strip for the given id
// !Important: can mess up the led strips
// returns an error in case of failure
func (c *Client) SetCustomDataById(data []byte, id string) error {
	path := fmt.Sprintf("leds/%s/control", id)
	token := c.mqttClient.Publish(path, QosAtLeastOnce, false, data)
	token.Wait()

	return token.Error()
}
