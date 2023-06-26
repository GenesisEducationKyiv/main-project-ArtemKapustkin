package repository

import (
	"bitcoin-exchange-rate/internal/model"
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

var ErrEmailAlreadyExist = errors.New("subscriber already exists in the file")
var ErrFileIsNotEmpty = errors.New("file is not empty")

type SubscriberFileRepository struct {
	filePath string
}

func NewSubscriberFileRepository(filePath string) *SubscriberFileRepository {
	return &SubscriberFileRepository{
		filePath: filePath,
	}
}

func (r *SubscriberFileRepository) GetAll() ([]*model.Subscriber, error) {
	file, err := os.Open(r.filePath)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("error closing file: %s", err)
		}
	}()

	var subscribers []*model.Subscriber

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		subscribers = append(subscribers, model.NewSubscriber(strings.TrimSpace(scanner.Text())))
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return subscribers, nil
}

func (r *SubscriberFileRepository) isSubscriberExists(subscribers []string, subscriber *model.Subscriber) error {
	for _, sub := range subscribers {
		if sub == subscriber.GetEmail() {
			log.Printf("subscriber '%s' already exists", subscriber.GetEmail())
			return fmt.Errorf("%w, subscriber's email: %s", ErrEmailAlreadyExist, subscriber.GetEmail())
		}
	}
	return nil
}

func (r *SubscriberFileRepository) Create(subscriber *model.Subscriber) error {
	readFile, err := os.ReadFile(r.filePath)
	if err != nil {
		return err
	}

	subscribers := strings.Split(string(readFile), "\n")

	if err = r.isSubscriberExists(subscribers, subscriber); err != nil {
		return err
	}

	writeFile, err := os.OpenFile(r.filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	defer func() {
		if err := writeFile.Close(); err != nil {
			log.Printf("error closing file: %s", err)
		}
	}()

	if _, err := writeFile.WriteString("\n" + subscriber.GetEmail()); err != nil {
		return err
	}
	log.Printf("subscriber '%s' added successfully", subscriber.GetEmail())

	return nil
}

func ClearFile(filePath string) error {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	log.Printf("file '%s' cleared successfully", filePath)
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("error closing file: %s", err)
		}
	}()

	return nil
}
