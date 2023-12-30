package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/IBM/sarama"
)

const (
	CONST_HOST  = "localhost:8097"
	CONST_TOPIC = "ordersProducer"
	CONST_GROUP = "orders-group"
)

type Consumer struct {
}

func (*Consumer) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (*Consumer) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (consumer *Consumer) ConsumeClaim(
	sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		// process message
		kafkaMsg := KafkaMsg{}
		_ = json.Unmarshal(msg.Value, &kafkaMsg)
		fmt.Printf("%+v\n", kafkaMsg)

		sess.MarkMessage(msg, "")
	}
	return nil
}

type KafkaMsg struct {
	Order_id int            `json:"order_id"`
	Pr_info  []ProductsInfo `json:"pr_info"`
}

type ProductsInfo struct {
	Product_id int `json:"product_id"`
	Amount     int `json:"amount"`
}

func initializeConsumerGroup() (sarama.ConsumerGroup, error) {

	// Configure Consumer
	config := sarama.NewConfig()

	config.ClientID = "products-app"
	config.Admin.Retry.Backoff = 2000
	config.Admin.Retry.Max = 5

	config.Version = sarama.DefaultVersion

	config.Consumer.Offsets.AutoCommit.Enable = false
	config.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second

	consumerGroup, err := sarama.NewConsumerGroup(
		[]string{CONST_HOST}, CONST_GROUP, config)

	if err != nil {
		return nil, fmt.Errorf("failed to initialize consumer group: %w", err)
	}

	fmt.Println("ConsumerGroup initialised.")
	return consumerGroup, nil

}

func Consume(ctx context.Context) {

	consumerGroup, err := initializeConsumerGroup()
	if err != nil {
		log.Printf("initialization error: %v", err)
		return
	}

	defer consumerGroup.Close()

	consumer := &Consumer{}

	for {
		err = consumerGroup.Consume(ctx, []string{CONST_TOPIC}, consumer)
		if err != nil {
			log.Printf("error from consumer: %v\n", err)
		}

		if ctx.Err() != nil {
			return
		}
	}

}

func Produce(str string, limit int) {

	const (
		CONST_HOST  string = "localhost:8097"
		CONST_TOPIC string = "ordersConfirm"
	)

	config := sarama.NewConfig()

	config.ClientID = "products-app"
	config.Admin.Retry.Backoff = 2000
	config.Admin.Retry.Max = 5
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer([]string{CONST_HOST}, config)
	if err != nil {
		log.Fatal("failed to initialize NewSyncProducer, err:", err)
	}

	defer producer.Close()

	for i := 0; i < limit; i++ {
		msg := &sarama.ProducerMessage{Topic: CONST_TOPIC, Key: nil, Value: sarama.StringEncoder(str)}
		partition, offset, err := producer.SendMessage(msg)
		if err != nil {
			log.Println("SendMessage err: ", err)
			return
		}
		log.Printf("[producer] partition id: %d; offset: %d, value: %s\n", partition, offset, strings.ReplaceAll(str, "\n", ""))
	}

}
