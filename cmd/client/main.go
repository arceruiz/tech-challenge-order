package main

import (
	"tech-challenge-order/internal/channels/rest"
	"tech-challenge-order/internal/config"

	"github.com/sirupsen/logrus"
)

func main() {
	config.ParseFromFlags()
	if err := rest.New(
		rest.NewOrderChannel(),
	).
		Start(); err != nil {
		logrus.Panic()
	}
}
