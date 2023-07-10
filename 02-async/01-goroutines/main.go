package main

import (
	"log"
	"time"
)

type User struct {
	Email string
}

type UserRepository interface {
	CreateUserAccount(u User) error
}

type NotificationsClient interface {
	SendNotification(u User) error
}

type NewsletterClient interface {
	AddToNewsletter(u User) error
}

type Handler struct {
	repository          UserRepository
	newsletterClient    NewsletterClient
	notificationsClient NotificationsClient
}

func NewHandler(
	repository UserRepository,
	newsletterClient NewsletterClient,
	notificationsClient NotificationsClient,
) Handler {
	return Handler{
		repository:          repository,
		newsletterClient:    newsletterClient,
		notificationsClient: notificationsClient,
	}
}

func (h Handler) SignUp(u User) error {
	timesToRepeat := 5
	createUserFunc := func() error {
		return h.repository.CreateUserAccount(u)
	}

	userCreatedCh := make(chan struct{})
	go func() {
		runRepeatedly(createUserFunc, timesToRepeat)
		close(userCreatedCh)
	}()
	<-userCreatedCh

	addToNewsletterFunc := func() error {
		return h.newsletterClient.AddToNewsletter(u)
	}
	go func() {
		runRepeatedly(addToNewsletterFunc, timesToRepeat)
	}()

	sendNotificationFunc := func() error {
		return h.notificationsClient.SendNotification(u)
	}
	go func() {
		runRepeatedly(sendNotificationFunc, timesToRepeat)
	}()

	return nil
}

func runRepeatedly(runFunc func() error, timesToRepeat int) {
	for i := 0; i < timesToRepeat; i++ {
		if err := runFunc(); err != nil {
			log.Printf("\nrunning func: %v\n", err)
			time.Sleep(1 * time.Second)
			continue
		}

		break
	}
}
