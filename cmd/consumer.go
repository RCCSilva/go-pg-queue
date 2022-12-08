package main

import (
	"context"
	"database/sql"
	"log"
	"time"
)

type Consumer struct {
	repository *Repository
}

func NewConsumer(repository *Repository) *Consumer {
	return &Consumer{
		repository: repository,
	}
}

const maxAttempts = 100

func (c *Consumer) Consume(ctx context.Context, consumerId int) error {
	log.Printf("[consumer - %d] starting consumer %d ID", consumerId, consumerId)

	for {
		var simpleQueue *SimpleQueue
		var tx *sql.Tx

		err := Retry(ctx, maxAttempts, 5*time.Second, time.Minute, func() error {
			var err error
			tx, err = c.repository.BeginTx(ctx, &sql.TxOptions{})

			if err != nil {
				return err
			}

			simpleQueue, err = c.repository.DequeueFromSimpleQueue(ctx, tx)
			if err != nil {
				tx.Rollback()
				return err
			}

			return nil
		})

		if err != nil {
			log.Printf("[consumer - %d] id consumer failed", consumerId)
			return err
		}

		log.Printf("[consumer - %d] id %d - processing", consumerId, simpleQueue.ID)
		time.Sleep(2 * time.Second)

		err = c.repository.SaveSimpleQueue(ctx, tx, simpleQueue)
		if err != nil {
			log.Printf("[consumer - %d] failed to save simple queue: %v", consumerId, err)

			tx.Rollback()
			continue
		}

		err = tx.Commit()
		if err != nil {
			log.Printf("[consumer - %d] failed to commit", consumerId)

			tx.Rollback()
			continue
		}

		log.Printf("[consumer - %d] id %d - completed", consumerId, simpleQueue.ID)
	}
}
