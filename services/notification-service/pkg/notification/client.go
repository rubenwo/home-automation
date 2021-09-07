package notification

import (
	"context"
	"firebase.google.com/go/v4/messaging"
	"fmt"
	"github.com/go-pg/pg/v10"
	"log"
	"sync"
)

type Manager struct {
	notifications    chan Notification
	msgClient        *messaging.Client
	db               *pg.DB
	lock             sync.Mutex
	subscribeClients []NotificationSubscriber
}

func NewManager(msgClient *messaging.Client, db *pg.DB) *Manager {
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

func (m *Manager) NotifyDatasetChanged() {}

func (m *Manager) SendNotification(notification Notification) {
	m.notifications <- notification
}
