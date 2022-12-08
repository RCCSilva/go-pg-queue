package main

import (
	"context"
	"log"
)

func main() {
	config := NewConfig()

	db, err := ConnectDatabase(config)
	if err != nil {
		log.Panicf("failed to connect to the database: %v", err)
	}

	err = MigrateDatabase(config)
	if err != nil {
		log.Panicf("failed to migrate: %v", err)
	}

	repository := NewRepository(db)

	ctx := context.Background()

	for i := 0; i < 10; i++ {
		consumer := NewConsumer(repository)
		go consumer.Consume(ctx, i+1)
	}

	producer := NewProducer(repository)
	producer.Produce()
}
