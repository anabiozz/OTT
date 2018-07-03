package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/adjust/rmq"
	"github.com/anabiozz/OTT/utils"
)

const (
	numDeliveries = 100000000
	batchSize     = 1000
	unackedLimit  = 1000
	numConsumers  = 10
)

var shutdown = make(chan interface{})

func publish(done chan interface{}, taskQueue rmq.Queue, generatedStringChannel <-chan string) {
	var before time.Time
	for {
		select {
		case <-done:
			close(shutdown)
		case s := <-generatedStringChannel:
			taskQueue.Publish(s)
			duration := time.Now().Sub(before)
			before = time.Now()
			perSecond := time.Second / (duration / batchSize)
			log.Printf("produced %d", perSecond)
		}
	}
}

func consume(done chan interface{}, taskQueue rmq.Queue) {
	taskQueue.StartConsuming(unackedLimit, 500*time.Millisecond)
	for i := 0; i < numConsumers; i++ {
		name := fmt.Sprintf("consumer %d", i)
		taskQueue.AddConsumer(name, NewConsumer(i))
	}
	select {}
}

func main() {

	done := make(chan interface{})

	connection := rmq.OpenConnection("tasks", "tcp", "localhost:6379", 1)
	taskQueue := connection.OpenQueue("tasks")

	go publish(done, taskQueue, utils.Generator())
	go consume(done, taskQueue)

	consumer := &Consumer{}
	delivery := rmq.NewTestDeliveryString("task payload")
	consumer.Consume(delivery)

	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		sig := <-gracefulStop
		fmt.Printf("caught sig: %+v\n", sig)
		// fmt.Println("Wait for 2 second to finish processing")
		// time.Sleep(2 * time.Second)
		done <- true
	}()

	<-shutdown
}
