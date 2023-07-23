package repository

import (
	"bitcoin-exchange-rate/modules/notification_module/model"
	"database/sql"
	"fmt"
	"log"
)

type SubscriberDBRepository struct {
	db *sql.DB
}

func NewSubscriberDBRepository(db *sql.DB) *SubscriberDBRepository {
	return &SubscriberDBRepository{db: db}
}

func (r *SubscriberDBRepository) Create(subscriber *model.Subscriber) error {
	var email string
	query := "INSERT INTO subscriber_db.public.subscribers (email) VALUES ($1) RETURNING (email)"
	row := r.db.QueryRow(query, subscriber.GetEmail())
	if row != nil {
		if err := row.Scan(&email); err != nil {
			return fmt.Errorf("error inserting subscriber: %s", err)
		}
		log.Println(email)
	}
	return nil
}

func (r *SubscriberDBRepository) GetAll() ([]*model.Subscriber, error) {
	var subscribers []*model.Subscriber

	query := "SELECT email FROM subscriber_db.public.subscribers"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error fetching emails: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var subscriber string
		if err := rows.Scan(&subscriber); err != nil {
			return nil, fmt.Errorf("error scanning email: %w", err)
		}
		subscribers = append(subscribers, model.NewSubscriber(subscriber))
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return subscribers, nil
}
