package main

import (
	"context"
	"database/sql"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

const getPaymentFromQueueQuery = `UPDATE simple_queue SET lock = current_timestamp
WHERE id = (
	SELECT id FROM simple_queue
	WHERE (completed_at IS NULL) AND lock IS NULL
	LIMIT 1
	FOR UPDATE SKIP LOCKED)
RETURNING id`

const insertSimpleQueue = `INSERT INTO simple_queue (completed_at) VALUES (NULL);`

const saveSimpleQueue = `UPDATE simple_queue SET completed_at = current_timestamp, lock = NULL WHERE id = $1;`

func (r *Repository) BeginTx(ctx context.Context, txOptions *sql.TxOptions) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, txOptions)
}

func (r *Repository) DequeueFromSimpleQueue(ctx context.Context, tx *sql.Tx) (*SimpleQueue, error) {
	row := tx.QueryRow(getPaymentFromQueueQuery)

	var simpleQueue SimpleQueue

	err := row.Scan(&simpleQueue.ID)
	if err != nil {
		return nil, err
	}

	return &simpleQueue, err
}

func (r *Repository) SaveSimpleQueue(ctx context.Context, tx *sql.Tx, simpleQueue *SimpleQueue) error {
	_, err := tx.Exec(saveSimpleQueue, simpleQueue.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) EnqueueIntoSimpleQueue() error {
	_, err := r.db.Exec(insertSimpleQueue)
	return err
}
