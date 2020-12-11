package mq

import (
	"github.com/handsomestWei/go-nsq-connect/util"
	"github.com/nsqio/go-nsq"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var nsqProducers = make(map[string]*nsq.Producer) // 多生产者，按别名区分
var isOpenProducerListen bool

func InitNsqProducer(alia, addr string, config *nsq.Config) *nsq.Producer {
	// Instantiate a producer.
	nsqProducer, err := nsq.NewProducer(addr, config)
	if err != nil {
		log.Fatal(err)
		return nil
	} else {
		nsqProducers[alia] = nsqProducer
	}

	util.Synchronized(func() {
		if !isOpenProducerListen {
			isOpenConsumerListen = true
			go listenProducerSignal()
		}
	})

	return nsqProducer
}

func InitSimpleNsqProducer(alia, addr string) *nsq.Producer {
	return InitNsqProducer(alia, addr, nsq.NewConfig())
}

func GetNsqProducer(alia string) *nsq.Producer {
	return nsqProducers[alia]
}

func listenProducerSignal() {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	// wait for signal to exit
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// Gracefully stop the producer when appropriate (e.g. before shutting down the service)
	for _, producer := range nsqProducers {
		producer.Stop()
	}
	log.Println("nsq producer stop")
}
