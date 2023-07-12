package model

import "errors"

var ErrSubscriberAlreadyExist = errors.New("subscriber already exists in the file")
var ErrSubscriberFileIsEmpty = errors.New("there are no subscribers in file")

type Subscriber struct {
	email string
}

func NewSubscriber(email string) *Subscriber {
	return &Subscriber{
		email: email,
	}
}

func (s Subscriber) GetEmail() string {
	return s.email
}
