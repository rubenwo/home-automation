package mqtt

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"net/http"
)

type Client struct {
	clients map[string]mqtt.Client
}

func New() *Client {
	return &Client{clients: make(map[string]mqtt.Client)}
}

func (c *Client) Register(path, host string) error {
	if _, ok := c.clients[path]; !ok {
		fmt.Println("registering:", path, host)
		opts := mqtt.NewClientOptions()
		opts.AddBroker(fmt.Sprintf("%s:1883", host))
		opts.SetClientID(uuid.New().String())
		client := mqtt.NewClient(opts)
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			return token.Error()
		}
		c.clients[path] = client
	}
	fmt.Println(c.clients)
	return nil
}

func (c *Client) BrokerMQTTRequest(w http.ResponseWriter, r *http.Request) {
	raw, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Println(string(raw))
	if err := c.post(raw, r.URL.Path); err != nil {
		log.Println(err)
	}

	w.WriteHeader(http.StatusOK)
}

func (c *Client) post(data []byte, path string) error {
	fmt.Println(c.clients)

	client, ok := c.clients[path]
	if !ok {
		return fmt.Errorf("client for path: %s, was not registered", path)
	}
	token := client.Publish(path, 0, false, data)
	token.Wait()
	fmt.Println("POSTED")
	return token.Error()
}
