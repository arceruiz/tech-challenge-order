package main

import (
	"tech-challenge-order/internal/channels/rest"
	"tech-challenge-order/internal/channels/sqs"
	"tech-challenge-order/internal/config"

	"github.com/sirupsen/logrus"
)

const (
	PAYMENT = "payment"
	ORDER   = "order"
)

func main() {
	config.ParseFromFlags()

	sqsInstance := sqs.NewSQS()

	go func() {
		sqsInstance.ReceiveMessage(PAYMENT)
	}()

	go func() {
		sqsInstance.ReceiveMessage(ORDER)
	}()

	if err := rest.New(
		rest.NewOrderChannel(),
	).
		Start(); err != nil {
		logrus.Panic()
	}
}
