package mocks

import (
	"tech-challenge-order/internal/integration/payment"

	"github.com/stretchr/testify/mock"
)

type PaymentServiceMock struct {
	mock.Mock
}

func (m *PaymentServiceMock) Create(p payment.Payment) error {
	args := m.Called(p)
	return args.Error(0)
}
