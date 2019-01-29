package main

import (
	"fmt"
	"log"
	"os"
	"time"
	"github.com/Shopify/sarama"
	"github.com/wvanbergen/kafka/consumergroup"
)

const (
	zookeeperConn = "127.0.0.1:2181"
	cgroup        = "default-consumer-group"
	topic    = "topic-name"
)

func main() {
	// setup sarama log to stdout
	sarama.Logger = log.New(os.Stdout, "", log.Ltime)

	// init consumer
	cg, err := initConsumer()
	if err != nil {
		fmt.Println("Error consumer goup: ", err.Error())
		os.Exit(1)
	}
	defer cg.Close()

	// run consumer
	consume(cg)
}

func initConsumer() (*consumergroup.ConsumerGroup, error) {
	// consumer config
	config := consumergroup.NewConfig()
	config.Offsets.Initial = sarama.OffsetOldest
	config.Offsets.ProcessingTimeout = 10 * time.Second
	// join to consumer group
	cg, err := consumergroup.JoinConsumerGroup(cgroup, []string{topic}, []string{zookeeperConn}, config)
	if err != nil {
		return nil, err
	}
	return cg, err
}

func consume(cg *consumergroup.ConsumerGroup) {
	for {
		select {
		case msg := <-cg.Messages():
			// messages coming through channel
			// only take messages from subscribed topic
			if msg.Topic != topic {
				continue
			}

			fmt.Println("New message from topic: " + msg.Value)

			// commit to zookeeper that message is read
			// this prevent read message multiple times after restart
			err := cg.CommitUpto(msg)
			if err != nil {
				fmt.Println("Error committing to zookeeper: ", err.Error())
			}
		}
	}
}
