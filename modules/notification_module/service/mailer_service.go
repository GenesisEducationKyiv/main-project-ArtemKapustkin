package service

import (
	"bitcoin-exchange-rate/modules/notification_module/model"
	"fmt"
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

type Logger interface {
	Error(message string)
	Info(message string)
}

type MailerService struct {
	subscriptionRepository SubscriptionRepository
	exchangeRateService    RateService
	mailer                 Mailer
	baseMessageToSend      string
	logger                 Logger
}

func NewMailerService(subscriptionRepository SubscriptionRepository, exchangeRateService RateService, mailer Mailer, logger Logger) *MailerService {
	return &MailerService{
		subscriptionRepository: subscriptionRepository,
		exchangeRateService:    exchangeRateService,
		mailer:                 mailer,
		baseMessageToSend:      "Subject: BTCUAH Exchange Rate Update\n\nDear subscriber,\n\nHere is current BTCUAH exchange rate: %s\n\nSincerely,\nArtem Kapustkin Mailer",
		logger:                 logger,
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
			specifiedErr := fmt.Errorf("failed to send email to %s: %s", subscriber.GetEmail(), err)
			return specifiedErr
		}

		log := fmt.Sprintf("email sent successfully to %s", subscriber.GetEmail())
		s.logger.Info(log)
	}

	return nil
}

func (s *MailerService) SendExchangeRate() error {
	value, err := s.exchangeRateService.GetRate()
	if err != nil {
		s.logger.Error(err.Error())
		return err
	}

	err = s.sendValueToAllEmails(model.NewEmailMessage(strconv.FormatFloat(value, 'f', 2, 64)))
	if err != nil {
		s.logger.Error(err.Error())
		return err
	}

	return nil
}

func (s *MailerService) Subscribe(subscriber *model.Subscriber) error {

	err := s.subscriptionRepository.Create(subscriber)
	if err != nil {
		s.logger.Error(err.Error())
		return err
	}
	log := fmt.Sprintf("subscriber %s created successfully", subscriber.GetEmail())
	s.logger.Info(log)
	return nil
}
