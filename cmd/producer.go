package main

import (
	"time"
)

type Producer struct {
	repository *Repository
}

func NewProducer(repository *Repository) *Producer {
	return &Producer{
		repository: repository,
	}
}

func (p *Producer) Produce() {
	for {
		p.repository.EnqueueIntoSimpleQueue()
		time.Sleep(250 * time.Millisecond)
	}
}
