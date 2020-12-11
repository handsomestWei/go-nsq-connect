package mq

import (
	"github.com/handsomestWei/go-nsq-connect/util"
	"github.com/nsqio/go-nsq"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var consumers []*nsq.Consumer // 多消费者
var isOpenConsumerListen bool

type messageHandler struct {
	ProcessMessage func(msg []byte) error
}

// HandleMessage implements the Handler interface.
func (h *messageHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		// Returning nil will automatically send a FIN command to NSQ to mark the message as processed.
		// In this case, a message with an empty body is simply ignored/discarded.
		return nil
	}

	// do whatever actual message processing is desired
	err := h.ProcessMessage(m.Body)

	// Returning a non-nil error will automatically send a REQ command to NSQ to re-queue the message.
	return err
}

func InitSimpleNsqConsumer(topic, channel, addr string, fn func(msg []byte) error) {
	InitNsqConsumer(topic, channel, addr, nsq.NewConfig(), fn)
}

func InitNsqConsumer(topic, channel, addr string, config *nsq.Config, fn func(msg []byte) error) {
	// Instantiate a consumer that will subscribe to the provided channel.
	consumer, err := nsq.NewConsumer(topic, channel, config)
	if err != nil {
		log.Fatal(err)
	}

	// Set the Handler for messages received by this Consumer. Can be called multiple times.
	// See also AddConcurrentHandlers.
	handler := &messageHandler{
		ProcessMessage: fn,
	}
	consumer.AddHandler(handler)

	// Use nsqlookupd to discover nsqd instances.
	// See also ConnectToNSQD, ConnectToNSQDs, ConnectToNSQLookupds.
	err = consumer.ConnectToNSQLookupd(addr)
	if err != nil {
		log.Fatal(err)
	} else {
		consumers = append(consumers, consumer)
	}

	util.Synchronized(func() {
		if !isOpenConsumerListen {
			isOpenConsumerListen = true
			go listenConsumerSignal()
		}
	})

}

func listenConsumerSignal() {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	// wait for signal to exit
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// Gracefully stop the consumer.
	for _, consumer := range consumers {
		consumer.Stop()
	}
	log.Println("nsq consumer stop")
}
