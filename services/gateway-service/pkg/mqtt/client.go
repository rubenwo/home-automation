package mqtt

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Client struct {
	clients  map[string]mqtt.Client
	upgrader websocket.Upgrader
}

func New() *Client {
	return &Client{
		clients: make(map[string]mqtt.Client),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}}
}

func (c *Client) Register(path, host string, retry int) error {
	if _, ok := c.clients[path]; !ok {
		log.Println("registering:", path, host)
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
			return err
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

func (c *Client) SocketMQTTRequest(w http.ResponseWriter, r *http.Request) {
	conn, err := c.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	path := r.URL.Path
	client, ok := c.clients[path]
	if !ok {
		http.Error(w, fmt.Errorf("client for path: %s, was not registered", path).Error(), http.StatusBadRequest)
		return
	}
	go func() {
		client.Subscribe(path, 0, func(cl mqtt.Client, msg mqtt.Message) {
			msg.Ack()
			payload := msg.Payload()
			fmt.Println(string(payload))
			// Forward content of MQTT message to the websocket
			if err := conn.WriteMessage(websocket.TextMessage, payload); err != nil {
				log.Println(err)
				client.Unsubscribe(path)
			}
		})
	}()

	// read from websocket connection to not block the client
	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read: ", err)
			break
		}
		log.Printf("recv: %s\n", message)
		fmt.Println(mt)
	}
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
