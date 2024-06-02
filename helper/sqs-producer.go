package helper

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func ProduceMessageSqs(queueURL, messageBody, messageType string) error {
    // Load the AWS configuration with the specified region
    cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-southeast-1"))
    if err != nil {
        return err
    }

    client := sqs.NewFromConfig(cfg)

    _, err = client.SendMessage(context.TODO(), &sqs.SendMessageInput{
        MessageBody: aws.String(messageBody),
        QueueUrl:    &queueURL,
    })

    if err != nil {
        return err
    }

    log.Println("Message sent successfully to SQS queue - ", messageType)
    return nil
}