package service

import (
	"bitcoin-exchange-rate/internal/model"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockSubscriberRepository struct {
	Subscribers []*model.Subscriber
	Err         error
}

func (m *MockSubscriberRepository) GetAll() ([]*model.Subscriber, error) {
	return m.Subscribers, m.Err
}

type MockMailer struct {
	EmailsSent int
	Err        error
}

func (m *MockMailer) SendEmail(email, value string) error {
	m.EmailsSent++
	return m.Err
}

func TestMailerService_SendValueToAllEmails_Success(t *testing.T) {
	subscribers := []*model.Subscriber{}
	subscribers = append(subscribers, model.NewSubscriber("subscriber1@example.com"), model.NewSubscriber("subscriber2@example.com"))

	mockRepo := &MockSubscriberRepository{Subscribers: subscribers}
	mockMailer := &MockMailer{}

	mailerService := NewMailerService(mockRepo, mockMailer)

	err := mailerService.SendValueToAllEmails("10000")
	assert.NoError(t, err)
	assert.Equal(t, 2, mockMailer.EmailsSent)
}

func TestMailerService_SendValueToAllEmails_SubscriberFileIsEmpty(t *testing.T) {
	mockRepo := &MockSubscriberRepository{}
	mockMailer := &MockMailer{}

	mailerService := NewMailerService(mockRepo, mockMailer)

	err := mailerService.SendValueToAllEmails("10000")
	assert.Equal(t, ErrSubscriberFileIsEmpty, err)
	assert.Zero(t, mockMailer.EmailsSent)
}

func TestMailerService_SendValueToAllEmails_GetAllError(t *testing.T) {
	mockRepo := &MockSubscriberRepository{Err: errors.New("repository error")}
	mockMailer := &MockMailer{}

	mailerService := NewMailerService(mockRepo, mockMailer)

	err := mailerService.SendValueToAllEmails("10000")
	assert.EqualError(t, err, "repository error")
	assert.Zero(t, mockMailer.EmailsSent)
}

func TestMailerService_SendValueToAllEmails_SendEmailError(t *testing.T) {
	subscribers := []*model.Subscriber{}
	subscribers = append(subscribers, model.NewSubscriber("subscriber1@example.com"))

	mockRepo := &MockSubscriberRepository{Subscribers: subscribers}
	mockMailer := &MockMailer{Err: errors.New("mailer error")}

	mailerService := NewMailerService(mockRepo, mockMailer)

	err := mailerService.SendValueToAllEmails("10000")
	assert.EqualError(t, err, "mailer error")
	assert.Equal(t, 1, mockMailer.EmailsSent)
}
