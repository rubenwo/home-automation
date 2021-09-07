package notification

import (
	"context"
	"firebase.google.com/go/v4/messaging"
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/rubenwo/home-automation/servics/notification-service/internal/service"
	"log"
	"sync"
)

type Manager struct {
	notifications    chan Notification
	msgClient        *messaging.Client
	db               *pg.DB
	lock             sync.Mutex
	subscribeClients []service.NotificationSubscriber
}

func NewManager(msgClient *messaging.Client, db *pg.DB) *Manager {
	m := &Manager{
		msgClient:     msgClient,
		db:            db,
		notifications: make(chan Notification, 100),
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
		//token := "fzNfIKbuPJeraNBhx-p9wf:APA91bFSmGLxWk9kKQqeCB37u6elkj4FEGV0ju4rWQFVk3BshhQek0yAmQllp02UGj21BcuAlONNiqGGtjuqTPWGhdEKZ8GwW4pLDZTFVv1Ut-P2FfUTalFvnGg9yJ96tbObVZKEIH1M"
	}
}

func (m *Manager) NotifyDatasetChanged() {}
