package sqs

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"tech-challenge-order/internal/canonical"
	"tech-challenge-order/internal/config"
	"tech-challenge-order/internal/service"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/sirupsen/logrus"
)

var (
	once     sync.Once
	instance QueueInterface
)

const (
	PAYMENT = "payment"
	ORDER   = "order"
)

type QueueInterface interface {
	ReceiveMessage(string)
}

type queueSQS struct {
	sqsService        *sqs.SQS
	service           service.OrderService
	queuesAddress     map[string]string
	messsageProcessor map[string]func([]byte) error
}

func NewSQS() QueueInterface {
	once.Do(func() {
		sess := session.Must(session.NewSessionWithOptions(session.Options{
			Config: aws.Config{
				Endpoint:   aws.String(config.Get().SQS.Endpoint),
				Region:     aws.String(config.Get().SQS.Region),
				DisableSSL: aws.Bool(true),
			},
		}))

		sqs := &queueSQS{
			sqsService: sqs.New(sess),
			service:    service.NewOrderService(),
			queuesAddress: map[string]string{
				"order":   config.Get().SQS.OrderQueue,
				"payment": config.Get().SQS.PaymentPayedQueue,
			},
			messsageProcessor: make(map[string]func([]byte) error),
		}

		sqs.messsageProcessor[PAYMENT] = sqs.processPaymentMessage
		sqs.messsageProcessor[ORDER] = sqs.processOrderMessage

		instance = sqs
	})

	return instance
}

func (q *queueSQS) ReceiveMessage(queueToListen string) {
	for {
		logrus.Info(fmt.Printf("STARTING LISTENING TO %s", queueToListen))
		queue := q.queuesAddress[queueToListen]

		paramsOrder := &sqs.ReceiveMessageInput{
			QueueUrl:            &queue,
			MaxNumberOfMessages: aws.Int64(1),
		}

		resp, err := q.sqsService.ReceiveMessage(paramsOrder)
		if err != nil {
			continue
		}

		if len(resp.Messages) > 0 {
			for _, msg := range resp.Messages {

				err := q.processMessage(queueToListen, []byte(*msg.Body))
				if err != nil {
					continue
				}

				_, err = q.sqsService.DeleteMessage(&sqs.DeleteMessageInput{
					QueueUrl:      &queue,
					ReceiptHandle: msg.ReceiptHandle,
				})
				if err != nil {
					continue
				}
			}
		} else {
			logrus.Info("there aren't new messages")
			time.Sleep(time.Second * 10)
		}
	}
}

func (q *queueSQS) processMessage(queue string, msg []byte) error {
	f, ok := q.messsageProcessor[queue]
	if !ok {
		return fmt.Errorf("queue processor not found")
	}

	return f(msg)
}

func (q *queueSQS) processPaymentMessage(msg []byte) error {
	var orderId string

	err := json.Unmarshal(msg, &orderId)
	if err != nil {
		return err
	}

	err = q.service.UpdateStatus(context.Background(), orderId, canonical.ORDER_PAYED)
	if err != nil {
		return err
	}

	return nil
}

func (q *queueSQS) processOrderMessage(msg []byte) error {
	var orderId string

	err := json.Unmarshal(msg, &orderId)
	if err != nil {
		return err
	}

	_, err = q.service.CheckoutOrder(context.Background(), orderId)
	if err != nil {
		return err
	}

	return nil
}
