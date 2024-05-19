package helper

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/SyamSolution/transaction-service/internal/model"
	_ "github.com/joho/godotenv"
	"os"
)

func ProduceCreateTransactionMessageMail(message model.Message) error {
	producer, err := sarama.NewSyncProducer([]string{os.Getenv("KAFKA_BROKER")}, nil)
	if err != nil {
		return err
	}

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		return err
	}

	produceMessage := &sarama.ProducerMessage{
		Topic: "create-transaction",
		Value: sarama.StringEncoder(jsonMessage),
	}

	if _, _, err := producer.SendMessage(produceMessage); err != nil {
		return err
	}

	return nil
}

func ProduceCompletedTransactionMessageMail(message model.CompleteTransactionMessage) error {
	producer, err := sarama.NewSyncProducer([]string{os.Getenv("KAFKA_BROKER")}, nil)
	if err != nil {
		return err
	}

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		return err
	}

	produceMessage := &sarama.ProducerMessage{
		Topic: "completed-transaction",
		Value: sarama.StringEncoder(jsonMessage),
	}

	if _, _, err := producer.SendMessage(produceMessage); err != nil {
		return err
	}

	return nil
}

func ProduceOrderTicketMessage(message model.MessageOrderTicket) error {
	producer, err := sarama.NewSyncProducer([]string{os.Getenv("KAFKA_BROKER")}, nil)
	if err != nil {
		return err
	}

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		return err
	}

	produceMessage := &sarama.ProducerMessage{
		Topic: "order-ticket",
		Value: sarama.StringEncoder(jsonMessage),
	}

	if _, _, err := producer.SendMessage(produceMessage); err != nil {
		return err
	}

	return nil
}

func ProduceSuccessOrderTicketMessage(message model.MessageOrderTicket) error {
	producer, err := sarama.NewSyncProducer([]string{os.Getenv("KAFKA_BROKER")}, nil)
	if err != nil {
		return err
	}

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		return err
	}

	produceMessage := &sarama.ProducerMessage{
		Topic: "success-order-ticket",
		Value: sarama.StringEncoder(jsonMessage),
	}

	if _, _, err := producer.SendMessage(produceMessage); err != nil {
		return err
	}

	return nil
}

func ProduceFailedOrderTicketMessage(message model.MessageOrderTicket) error {
	producer, err := sarama.NewSyncProducer([]string{os.Getenv("KAFKA_BROKER")}, nil)
	if err != nil {
		return err
	}

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		return err
	}

	produceMessage := &sarama.ProducerMessage{
		Topic: "failed-order-ticket",
		Value: sarama.StringEncoder(jsonMessage),
	}

	if _, _, err := producer.SendMessage(produceMessage); err != nil {
		return err
	}

	return nil
}
