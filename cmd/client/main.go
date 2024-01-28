package main

import (
	"net/http"
	"tech-challenge-order/internal/channels/rest"
	"tech-challenge-order/internal/config"
	"tech-challenge-order/internal/integration/payment"
	"tech-challenge-order/internal/repository"
	"tech-challenge-order/internal/service"

	"github.com/sirupsen/logrus"
)

var (
	cfg = &config.Cfg
)

func main() {
	config.ParseFromFlags()
	if err := rest.New(
		rest.NewOrderChannel(
			service.NewOrderService(
				repository.NewOrderRepo(
					repository.NewMongo(),
				),
				payment.NewPaymentService(http.DefaultClient),
			),
		),
	).
		Start(); err != nil {
		logrus.Panic()
	}
}
