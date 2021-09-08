package notification

import (
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"

	"context"
	"firebase.google.com/go/v4/messaging"
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
	"log"
	"sync"
	"time"
)

const (
	NotificationsMqttPath = "notifications"

	QosAtMostOnce  = 0
	QosAtLeastOnce = 1
	QosExactlyOnce = 2
)

type Manager struct {
	notifications    chan Notification
	msgClient        *messaging.Client
	db               *pg.DB
	lock             sync.Mutex
	subscribeClients []NotificationSubscriber
}

func NewManager(host string, retry int, msgClient *messaging.Client, db *pg.DB) *Manager {
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
		log.Fatal(err)
	}

	var notificationSubscribers []NotificationSubscriber
	if err := db.Model(&notificationSubscribers).Select(); err != nil {
		log.Fatal(err)
		return nil
	}

	m := &Manager{
		msgClient:        msgClient,
		db:               db,
		notifications:    make(chan Notification, 100),
		subscribeClients: notificationSubscribers,
	}
	go m.run()

	client.Subscribe(NotificationsMqttPath, QosAtLeastOnce, func(cl mqtt.Client, msg mqtt.Message) {
		msg.Ack()
		payload := msg.Payload()
		var nm Notification
		if err := json.Unmarshal(payload, &nm); err != nil {
			log.Fatal(err)
		}

		m.SendNotification(nm)
	})

	return m
}

func (m *Manager) run() {
	for notification := range m.notifications {
		m.lock.Lock()
		for _, sub := range m.subscribeClients {
			// See documentation on defining a message payload.
			message := &messaging.Message{
				Notification: &messaging.Notification{
					Title: notification.Title,
					Body:  notification.Body,
				},
				Token: sub.FirebaseMessagingToken,
			}
			send, err := m.msgClient.Send(context.Background(), message)
			if err != nil {
				log.Fatalf("error sending message: %v\n", err)
			}
			fmt.Println(send)
		}
		m.lock.Unlock()
	}
}

func (m *Manager) NotifyDatasetChanged() {
	var notificationSubscribers []NotificationSubscriber
	if err := m.db.Model(&notificationSubscribers).Select(); err != nil {
		log.Println(err)
		return
	}

	m.lock.Lock()
	m.subscribeClients = notificationSubscribers
	m.lock.Unlock()
}

func (m *Manager) SendNotification(notification Notification) {
	m.notifications <- notification
}
