package main

import (
	"fmt"
	"log"
	"time"

	"github.com/adjust/rmq"
)

// Consumer ...
type Consumer struct {
	name   string
	count  int
	before time.Time
}

// NewConsumer ...
func NewConsumer(tag int) *Consumer {
	return &Consumer{
		name:   fmt.Sprintf("consumer%d", tag),
		count:  0,
		before: time.Now(),
	}
}

// Consume ...
func (consumer *Consumer) Consume(delivery rmq.Delivery) {
	consumer.count++
	if consumer.count%batchSize == 0 {
		duration := time.Now().Sub(consumer.before)
		consumer.before = time.Now()
		perSecond := time.Second / (duration / batchSize)
		log.Printf("%s consumed %d %s %d", consumer.name, consumer.count, delivery.Payload(), perSecond)
	}
	time.Sleep(time.Millisecond)
	if consumer.count%batchSize == 0 {
		delivery.Reject()
	} else {
		delivery.Ack()
	}
}
