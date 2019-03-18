package api

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/snappy"
	"time"
)

func Write(msg []byte) error {
	// connect to kafka
	kafkaProducer, err := Configure([]string{"192.168.1.17:32769"}, "my-kafka-client", "test-topic")
	if err != nil {
		return err
	}
	defer kafkaProducer.Close()

	parent := context.Background()
	defer parent.Done()

	err = Push(parent, nil, msg)
	if err != nil {
		fmt.Printf("error: " + err.Error())
	}

	return nil
}

func Push(parent context.Context, key, value []byte) (err error) {
	message := kafka.Message{
		Key:   key,
		Value: value,
		Time:  time.Now(),
	}

	return writer.WriteMessages(parent, message)
}

var writer *kafka.Writer

func Configure(kafkaBrokerUrls []string, clientId string, topic string) (w *kafka.Writer, err error) {
	dialer := &kafka.Dialer{
		Timeout:  10 * time.Second,
		ClientID: clientId,
	}

	config := kafka.WriterConfig{
		Brokers:          kafkaBrokerUrls,
		Topic:            topic,
		Balancer:         &kafka.LeastBytes{},
		Dialer:           dialer,
		WriteTimeout:     10 * time.Second,
		ReadTimeout:      10 * time.Second,
		CompressionCodec: snappy.NewCompressionCodec(),
	}
	w = kafka.NewWriter(config)
	writer = w
	return w, nil
}
