package payment

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"tech-challenge-order/internal/config"
)

const (
	paymentEndpoint = "/api/payment/"
)

var (
	cfg = &config.Cfg
)

type paymentService struct {
	httpClient *http.Client
}

func NewPaymentService(client *http.Client) PaymentService {
	return &paymentService{
		httpClient: client,
	}
}

func (s *paymentService) Create(payment Payment) error {

	jsonPayment, err := json.Marshal(payment)
	if err != nil {
		return err
	}

	request, err := http.NewRequest("POST", cfg.Server.PaymentIntegrationHost+paymentEndpoint, bytes.NewBuffer(jsonPayment))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("payment integration error, code: " + strconv.Itoa(resp.StatusCode))
	}

	return nil
}
