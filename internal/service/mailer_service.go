package service

import (
	"bitcoin-exchange-rate/internal/model"
	"fmt"
	"log"
	"strconv"
)

type SubscriptionRepository interface {
	GetAll() ([]*model.Subscriber, error)
	Create(subscriber *model.Subscriber) error
}

type Mailer interface {
	SendEmail(email, value string) error
}

type RateService interface {
	GetRate() (float64, error)
}

type MailerService struct {
	subscriptionRepository SubscriptionRepository
	exchangeRateService    RateService
	mailer                 Mailer
	baseMessageToSend      string
}

func NewMailerService(subscriptionRepository SubscriptionRepository, exchangeRateService RateService, mailer Mailer) *MailerService {
	return &MailerService{
		subscriptionRepository: subscriptionRepository,
		exchangeRateService:    exchangeRateService,
		mailer:                 mailer,
		baseMessageToSend:      "Subject: BTCUAH Exchange Rate Update\n\nDear subscriber,\n\nHere is current BTCUAH exchange rate: %s\n\nSincerely,\nArtem Kapustkin Mailer",
	}
}

func (s *MailerService) sendValueToAllEmails(emailMessage model.EmailMessage) error {
	subscribers, err := s.subscriptionRepository.GetAll()
	if err != nil {
		return err
	}

	if len(subscribers) == 0 {
		return model.ErrSubscriberFileIsEmpty
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

func (s *MailerService) SendExchangeRate() error {
	value, err := s.exchangeRateService.GetRate()
	if err != nil {
		return err
	}

	err = s.sendValueToAllEmails(model.NewEmailMessage(strconv.FormatFloat(value, 'f', 2, 64)))
	if err != nil {
		return err
	}

	return nil
}
