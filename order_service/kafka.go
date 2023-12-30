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
	CONST_HOST       string = "localhost:8097"
	CONST_PROD_TOPIC string = "ordersProducer"
	CONST_CONS_TOPIC string = "orderConfirm"
	CONST_GROUP      string = "order-confirm-group"
)

type KafkaMsg struct {
	Order_id int            `json:"order_id"`
	Pr_info  []ProductsInfo `json:"pr_info"`
}

type ProductsInfo struct {
	Product_id int `json:"product_id"`
	Amount     int `json:"amount"`
}

type OrderConfirmation struct {
	Order_id int    `json:"order_id"`
	Status   string `json:"status"`
}

type Consumer struct {
}

func (*Consumer) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (*Consumer) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (consumer *Consumer) ConsumeClaim(
	sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	for msg := range claim.Messages() {
		fmt.Print("\n\n------------------------------------\n")
		fmt.Printf("Kafka: message consumed at %s.\n\n", CONST_CONS_TOPIC)

		// process message
		ordConfirm := OrderConfirmation{}
		_ = json.Unmarshal(msg.Value, &ordConfirm)
		fmt.Printf("%+v\n", ordConfirm)

		_, err := Db.Exec(`UPDATE orders SET status=? WHERE id=?;`, ordConfirm.Status, ordConfirm.Order_id)
		if err != nil {
			return err

		}

		sess.MarkMessage(msg, "")
	}
	return nil
}

func Produce(str string, limit int) {

	config := sarama.NewConfig()

	config.ClientID = "orders-app"
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
		msg := &sarama.ProducerMessage{Topic: CONST_PROD_TOPIC, Key: nil, Value: sarama.StringEncoder(str)}
		partition, offset, err := producer.SendMessage(msg)
		if err != nil {
			log.Println("SendMessage err: ", err)
			return
		}
		formatValue := strings.ReplaceAll(str, "\n", "")
		formatValue = strings.ReplaceAll(formatValue, " ", "")
		fmt.Print("\n\n------------------------------------\n")
		fmt.Printf("Kafka: message produced at %s.\n\n", CONST_PROD_TOPIC)

		log.Printf("[producer] partition id: %d; offset: %d, value: %s\n", partition, offset, formatValue)
	}

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

	fmt.Println("--> Consumer group initialized.")
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
		err = consumerGroup.Consume(ctx, []string{CONST_CONS_TOPIC}, consumer)
		if err != nil {
			log.Printf("error from consumer: %v\n", err)
		}

		if ctx.Err() != nil {
			return
		}
	}

}
