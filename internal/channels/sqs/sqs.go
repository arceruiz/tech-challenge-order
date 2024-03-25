package sqs

import (
	"context"
	"encoding/json"
	"sync"
	"tech-challenge-order/internal/canonical"
	"tech-challenge-order/internal/config"
	"tech-challenge-order/internal/service"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/rs/zerolog/log"
)

var (
	once     sync.Once
	instance QueueInterface
)

const (
	orderQueue            = "orderQueue"
	paymentPayedQueue     = "paymentPayedQueue"
	paymentCancelledQueue = "paymentCancelledQueue"
)

type QueueInterface interface {
	Start()
}

type queueSQS struct {
	sqsService   *sqs.SQS
	service      service.OrderService
	queueAddress map[string]string
}

func NewSQS() QueueInterface {
	once.Do(func() {
		sess := session.Must(session.NewSessionWithOptions(session.Options{
			Config: aws.Config{
				Region:     aws.String(config.Get().SQS.Region),
				DisableSSL: aws.Bool(true),
			},
		}))

		sqs := &queueSQS{
			sqsService: sqs.New(sess),
			service:    service.NewOrderService(),
			queueAddress: map[string]string{
				orderQueue:            config.Get().SQS.OrderQueue,
				paymentPayedQueue:     config.Get().SQS.PaymentPayedQueue,
				paymentCancelledQueue: config.Get().SQS.PaymentCancelledQueue,
			},
		}

		instance = sqs
	})

	return instance
}

func (q *queueSQS) Start() {
	orderChannel := make(chan *sqs.Message)
	paymentPayedChannel := make(chan *sqs.Message)
	paymentCancelledChannel := make(chan *sqs.Message)

	go q.receiveMessage(q.queueAddress[orderQueue], orderChannel)
	go q.receiveMessage(q.queueAddress[paymentPayedQueue], paymentPayedChannel)
	go q.receiveMessage(q.queueAddress[paymentCancelledQueue], paymentCancelledChannel)

	q.messageProcessor(orderChannel, paymentPayedChannel, paymentCancelledChannel)
}

func (q *queueSQS) messageProcessor(orderChannel, paymentPayedChannel, paymentCancelledChannel chan *sqs.Message) {
	for {
		select {
		case orderMessage := <-orderChannel:

			log.Info().Any("msg_id", orderMessage.MessageId).Msg("msg received from order queue")

			orderId := unmarshalMessageToId(orderMessage)

			_, err := q.service.CheckoutOrder(context.Background(), orderId)
			if err != nil {
				log.Err(err).Any("order_id", orderId).Msg("an error occurred when checkout order")
			}

			q.deleteMessage(orderMessage, q.queueAddress[orderQueue])

		case paymentPayedMessage := <-paymentPayedChannel:
			log.Info().Any("msg_id", paymentPayedMessage.MessageId).Msg("msg received from payment payed queue")

			orderId := unmarshalMessageToId(paymentPayedMessage)

			err := q.service.UpdateStatus(context.Background(), orderId, canonical.ORDER_PAYED)
			if err != nil {
				log.Err(err).Any("order_id", orderId).Msg("an error occurred when update status")
			}

			q.deleteMessage(paymentPayedMessage, q.queueAddress[paymentPayedQueue])

		case paymentCancelledMessage := <-paymentCancelledChannel:
			log.Info().Any("msg_id", paymentCancelledMessage.MessageId).Msg("msg received from payment payed queue")

			orderId := unmarshalMessageToId(paymentCancelledMessage)

			err := q.service.UpdateStatus(context.Background(), orderId, canonical.ORDER_CANCELLED)
			if err != nil {
				log.Err(err).Any("order_id", orderId).Msg("an error occurred when update status")
			}

			q.deleteMessage(paymentCancelledMessage, q.queueAddress[paymentCancelledQueue])
		}
	}
}

func (q *queueSQS) receiveMessage(queueToListen string, ch chan<- *sqs.Message) {
	for {
		paramsOrder := &sqs.ReceiveMessageInput{
			QueueUrl:            &queueToListen,
			MaxNumberOfMessages: aws.Int64(1),
		}

		resp, err := q.sqsService.ReceiveMessage(paramsOrder)
		if err != nil {
			log.Err(err).Msg("an error occurred when receive message from the queue")
			continue
		}

		if len(resp.Messages) > 0 {
			for _, msg := range resp.Messages {
				ch <- msg
			}
		} else {
			log.Info().Msg("no new messages")
			time.Sleep(time.Second * 10)
		}
	}
}

func (q *queueSQS) deleteMessage(msg *sqs.Message, queue string) {
	_, err := q.sqsService.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      &queue,
		ReceiptHandle: msg.ReceiptHandle,
	})
	if err != nil {
		log.Err(err).Any("msg_id", msg.MessageId).Msg("an error occurred when delete message")
	}
}

func unmarshalMessageToId(msg *sqs.Message) string {
	var id string

	err := json.Unmarshal([]byte(*msg.Body), &id)
	if err != nil {
		log.Err(err).Any("msg_id", msg.MessageId).Msg("an error occurred when process message")
	}

	return id
}
