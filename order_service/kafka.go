package main

import (
	"log"
	"strings"
	"time"

	"github.com/IBM/sarama"
)

type KafkaMsg struct {
	Order_id int            `json:"order_id"`
	Pr_info  []ProductsInfo `json:"pr_info"`
}

type ProductsInfo struct {
	Product_id int `json:"product_id"`
	Amount     int `json:"amount"`
}

func Produce(str string, limit int) {

	const (
		CONST_HOST  string = "localhost:8097"
		CONST_TOPIC string = "ordersProducer"
	)

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
		msg := &sarama.ProducerMessage{Topic: CONST_TOPIC, Key: nil, Value: sarama.StringEncoder(str)}
		partition, offset, err := producer.SendMessage(msg)
		if err != nil {
			log.Println("SendMessage err: ", err)
			return
		}
		formatValue := strings.ReplaceAll(str, "\n", "")
		formatValue = strings.ReplaceAll(formatValue, " ", "")
		log.Printf("[producer] partition id: %d; offset: %d, value: %s\n", partition, offset, formatValue)
	}

}

func Consume() {

	const (
		CONST_HOST  = "localhost:8097"
		CONST_TOPIC = "ordersConfirm"
	)

	// Configure Consumer
	config := sarama.NewConfig()

	config.ClientID = "orders-app"
	config.Admin.Retry.Backoff = 2000
	config.Admin.Retry.Max = 5

	config.Version = sarama.DefaultVersion

	// config.Consumer.Group.InstanceId = "products-group"
	config.Consumer.Offsets.AutoCommit.Enable = false
	config.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second

	// Create consumer
	var err error
	Consumer, err := sarama.NewConsumer([]string{CONST_HOST}, config)
	if err != nil {
		log.Fatal("NewConsumer err: ", err)
	}

	defer Consumer.Close()

	partitionList, _ := Consumer.Partitions(CONST_TOPIC)
	//fmt.Println(partitionList)

	messages := make(chan *sarama.ConsumerMessage, 512)
	initialOffset := sarama.OffsetOldest //offset to start reading message from
	for _, partition := range partitionList {
		pc, _ := Consumer.ConsumePartition(CONST_TOPIC, partition, initialOffset)
		go func(pc sarama.PartitionConsumer) {
			for message := range pc.Messages() {
				messages <- message //or call a function that writes to disk
			}
		}(pc)
	}

}
