package demo

import (
	"github.com/handsomestWei/go-nsq-connect/mq"
	"time"
	"log"
)

func Run() {
	initConsumer()

	producerAliaName := "testProduce"
	initProducer(producerAliaName)

	for {
		t := time.NewTimer(30 * time.Second)
		<-t.C

		messageBody := []byte(time.Now().Format(time.RFC3339) + " hello")
		topicName := "test-topic"

		// Synchronously publish a single message to the specified topic.
		// Messages can also be sent asynchronously and/or in batches.
		err := mq.GetNsqProducer(producerAliaName).Publish(topicName, messageBody)
		if err != nil {
			log.Fatal(err)
		}
	}

}

func initProducer(producerAliaName string) {
	mq.InitSimpleNsqProducer(producerAliaName, "172.16.21.11:4150")
}

func initConsumer() {
	go mq.InitSimpleNsqConsumer("test-topic", "test-channel", "172.16.21.11:4161", func(msg []byte) error {
		log.Println("processMessage: " + string(msg))
		return nil
	})
}
