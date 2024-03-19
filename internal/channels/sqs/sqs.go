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

type QueueInterface interface {
	Start()
}

type queueSQS struct {
	sqsService     *sqs.SQS
	service        service.OrderService
	queueProcessor map[string]func(queue string, ch chan *sqs.Message)
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
			sqsService:     sqs.New(sess),
			service:        service.NewOrderService(),
			queueProcessor: make(map[string]func(queue string, ch chan *sqs.Message)),
		}

		sqs.queueProcessor[config.Get().SQS.PaymentPayedQueue] = sqs.processPaymentPayedMessage
		sqs.queueProcessor[config.Get().SQS.OrderQueue] = sqs.processOrderMessage
		sqs.queueProcessor[config.Get().SQS.PaymentCancelledQueue] = sqs.processPaymentCancelledMessage

		instance = sqs
	})

	return instance
}

func (q *queueSQS) Start() {
	for queue, processor := range q.queueProcessor {
		channel := make(chan *sqs.Message)
		go q.receiveMessage(queue, channel)
		go processor(queue, channel)
	}
}

func (q *queueSQS) receiveMessage(queueToListen string, ch chan<- *sqs.Message) {
	for {
		logrus.Info(fmt.Printf("STARTING LISTENING TO %s", queueToListen))

		paramsOrder := &sqs.ReceiveMessageInput{
			QueueUrl:            &queueToListen,
			MaxNumberOfMessages: aws.Int64(1),
		}

		resp, err := q.sqsService.ReceiveMessage(paramsOrder)
		if err != nil {
			continue
		}

		if len(resp.Messages) > 0 {
			for _, msg := range resp.Messages {
				ch <- msg
			}
		} else {
			logrus.Info("there aren't new messages")
			time.Sleep(time.Second * 10)
		}
	}
}

func (q *queueSQS) processPaymentPayedMessage(queue string, ch chan *sqs.Message) {
	for {
		msg := <-ch

		var orderId string

		err := json.Unmarshal([]byte(*msg.Body), &orderId)
		if err != nil {
			logrus.WithError(err).WithField("order_id", orderId).Error("an error occurred when parse sqs message")
		}

		err = q.service.UpdateStatus(context.Background(), orderId, canonical.ORDER_PAYED)
		if err != nil {
			logrus.WithError(err).WithField("order_id", orderId).Error("an error occurred when update order status")
		}

		q.deleteMessage(msg, queue)
	}
}

func (q *queueSQS) processPaymentCancelledMessage(queue string, ch chan *sqs.Message) {
	for {
		msg := <-ch

		var orderId string

		err := json.Unmarshal([]byte(*msg.Body), &orderId)
		if err != nil {
			logrus.WithError(err).WithField("order_id", orderId).Error("an error occurred when parse sqs message")
		}

		err = q.service.UpdateStatus(context.Background(), orderId, canonical.ORDER_CANCELLED)
		if err != nil {
			logrus.WithError(err).WithField("order_id", orderId).Error("an error occurred when update order status")
		}

		q.deleteMessage(msg, queue)
	}
}

func (q *queueSQS) processOrderMessage(queue string, ch chan *sqs.Message) {
	for {
		msg := <-ch

		var orderId string

		err := json.Unmarshal([]byte(*msg.Body), &orderId)
		if err != nil {
			logrus.WithError(err).WithField("order_id", orderId).Error("an error occurred when parse sqs message")
		}

		_, err = q.service.CheckoutOrder(context.Background(), orderId)
		if err != nil {
			logrus.WithError(err).WithField("order_id", orderId).Error("an error occurred when checkout order")
		}

		q.deleteMessage(msg, queue)
	}
}

func (q *queueSQS) deleteMessage(msg *sqs.Message, queue string) {
	_, err := q.sqsService.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      &queue,
		ReceiptHandle: msg.ReceiptHandle,
	})
	if err != nil {
		logrus.WithError(err).Error("an error occurred when deleting message")
	}
}
