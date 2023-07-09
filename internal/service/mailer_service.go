package service

import (
	"bitcoin-exchange-rate/internal/model"
	"errors"
	"fmt"
	"log"
)

var ErrSubscriberFileIsEmpty = errors.New("there are no subscribers in file")

type SubscriberRepository interface {
	GetAll() ([]*model.Subscriber, error)
}

type Mailer interface {
	SendEmail(email, value string) error
}

type MailerService struct {
	subscriberRepository SubscriberRepository
	mailer               Mailer
	baseMessageToSend    string
}

func NewMailerService(subscriberRepository SubscriberRepository, mailer Mailer) *MailerService {
	return &MailerService{
		subscriberRepository: subscriberRepository,
		mailer:               mailer,
		baseMessageToSend:    "Subject: BTCUAH Exchange Rate Update\n\nDear subscriber,\n\nHere is current BTCUAH exchange rate: %s\n\nSincerely,\nArtem Kapustkin Mailer",
	}
}

func (s *MailerService) SendValueToAllEmails(emailMessage model.EmailMessage) error {
	subscribers, err := s.subscriberRepository.GetAll()
	if err != nil {
		return err
	}

	if len(subscribers) == 0 {
		return ErrSubscriberFileIsEmpty
	}

	for _, subscriber := range subscribers {
		message := fmt.Sprintf(s.baseMessageToSend, emailMessage)

		err := s.mailer.SendEmail(subscriber.GetEmail(), message)
		if err != nil {
			log.Printf("failed to send email to %s: %s", subscriber.GetEmail(), err)
			return err
		}

		log.Printf("email sent successfully to %s", subscriber.GetEmail())
	}

	return nil
}
