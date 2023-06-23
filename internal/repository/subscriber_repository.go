package repository

import (
	"bitcoin-exchange-rate/internal/model"
	"bufio"
	"errors"
	"log"
	"os"
	"strings"
)

var ErrEmailAlreadyExist = errors.New("subscriber already exists in the file")

type SubscriberFileRepository struct {
	filePath string
}

func NewSubscriberFileRepository(filePath string) *SubscriberFileRepository {
	return &SubscriberFileRepository{
		filePath: filePath,
	}
}

func (r *SubscriberFileRepository) SetFilePath(filePath string) {
	r.filePath = filePath
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

func (r *SubscriberFileRepository) Create(subscriber *model.Subscriber) error {
	content, err := os.ReadFile(r.filePath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")

	for _, line := range lines {
		if line == subscriber.GetEmail() {
			log.Printf("subscriber '%s' already exists", subscriber.GetEmail())
			return ErrEmailAlreadyExist
		}
	}

	file, err := os.OpenFile(r.filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("error closing file: %s", err)
		}
	}()

	if _, err := file.WriteString("\n" + subscriber.GetEmail()); err != nil {
		return err
	}

	log.Printf("subscriber '%s' added successfully", subscriber.GetEmail())
	return nil
}

func (r *SubscriberFileRepository) ClearFile() error {
	file, err := os.OpenFile(r.filePath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("error closing file: %s", err)
		}
	}()

	return nil
}

func (r *SubscriberFileRepository) isFileEmpty(filePath string) (bool, error) {
	// Get file information
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return false, err
	}

	// Check if the file size is zero
	if fileInfo.Size() == 0 {
		return true, nil
	}

	return false, nil
}
