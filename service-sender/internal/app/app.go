package app

import (
	"fmt"
	"os"
	"os/signal"
	"sender/internal/config"
	"sender/internal/senders/email"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func Run() {

	cfg := config.Get()

	go listen("myGroupEmail", email.SendEmail, cfg)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop

}

func listen(group string, f func(*config.Config, []byte) error, cfg *config.Config) {

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":        cfg.Host,
		"group.id":                 group,
		"auto.offset.reset":        "earliest",
		"enable.auto.offset.store": false,
	})
	if err != nil {
		panic(err)
	}

	err = c.SubscribeTopics([]string{cfg.EmailTopic}, nil)
	if err != nil {
		panic(err)
	}

	run := true
	for run {
		msg, err := c.ReadMessage(time.Second)
		if err == nil {

			err := f(cfg, msg.Value)
			if err != nil {
				fmt.Printf("Consumer error: %v\n", err)
				continue
			}

			_, err = c.StoreMessage(msg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%% Error storing offset after message %s:\n",
					msg.TopicPartition)
			}

		} else if !err.(kafka.Error).IsTimeout() {

			fmt.Printf("Consumer error: %v (%v)\n", err, msg)

		}
	}

	c.Close()

}
