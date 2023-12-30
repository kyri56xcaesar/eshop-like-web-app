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
	CONST_HOST       = "localhost:8097"
	CONST_CONS_TOPIC = "ordersProducer"
	CONST_PROD_TOPIC = "orderConfirm"
	CONST_GROUP      = "orders-group"
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
		kafkaMsg := KafkaMsg{}
		_ = json.Unmarshal(msg.Value, &kafkaMsg)
		// fmt.Printf("%+v\n", kafkaMsg)

		// Validate order
		result, _ := validateOrder(kafkaMsg)

		// Produce OrderConfirmation Message
		ord_confirm := OrderConfirmation{Order_id: kafkaMsg.Order_id, Status: result}
		jmsg, _ := json.Marshal(ord_confirm)
		message, _ := formatJSON(jmsg)
		Produce(string(message), 1)

		sess.MarkMessage(msg, "")
	}
	return nil
}

func validateOrder(msg KafkaMsg) (string, error) {
	fmt.Printf("%+v\n", msg)

	valid := true
	quantities := make(map[int]int)

	for _, product := range msg.Pr_info {

		if product.Amount <= 0 || product.Product_id <= 0 {
			log.Printf("Invalid Order")
			return "Reject", nil
		}

		row := Db.QueryRow(`SELECT id, quantity FROM products WHERE id=?;`, product.Product_id)

		pr := Product{}
		if err := row.Scan(&pr.ID, &pr.Quantity); err != nil {

			log.Printf("Scanning error: %v", err)
			return "Failure", err
		}
		fmt.Println("\tRetrieved product.")
		fmt.Printf("\t\t%+v\n", pr)

		quantities[pr.ID] = pr.Quantity

		if pr.Quantity < product.Amount {
			// Invalid order
			// Don't update db, return Reject
			valid = false
			fmt.Println("\t->Failure~!")
		} else {
			fmt.Println("\t->Success~!")
		}

	}

	if valid {
		// Valid order
		// Update db before returning
		for _, product := range msg.Pr_info {

			new_amount := quantities[product.Product_id] - product.Amount
			fmt.Printf("New amount: %d\n", new_amount)
			_, err := Db.Exec(`UPDATE products SET quantity=? WHERE id=?;`, new_amount, product.Product_id)

			if err != nil {
				return "Faulure", err

			}
		}

		return "Success", nil
	} else {
		return "Reject", nil
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

func Produce(str string, limit int) {

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
