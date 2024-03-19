package rest

import (
	"context"
	"tech-challenge-order/internal/canonical"

	"github.com/stretchr/testify/mock"
)

type OrderServiceMock struct {
	mock.Mock
}

func (m *OrderServiceMock) GetAll(ctx context.Context) ([]canonical.Order, error) {
	args := m.Called(ctx)
	return args.Get(0).([]canonical.Order), args.Error(1)
}

func (m *OrderServiceMock) Create(ctx context.Context, order canonical.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *OrderServiceMock) Update(ctx context.Context, id string, order canonical.Order) error {
	args := m.Called(id)

	return args.Error(0)
}

func (m *OrderServiceMock) GetByID(ctx context.Context, id string) (*canonical.Order, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*canonical.Order), args.Error(1)
}

func (m *OrderServiceMock) GetByStatus(ctx context.Context, status canonical.OrderStatus) ([]canonical.Order, error) {
	args := m.Called(ctx, status)
	return args.Get(0).([]canonical.Order), args.Error(1)
}

func (m *OrderServiceMock) CheckoutOrder(ctx context.Context, id string) (*canonical.Order, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*canonical.Order), args.Error(1)
}

func (m *OrderServiceMock) UpdateStatus(ctx context.Context, orderId string, status canonical.OrderStatus) error {
	args := m.Called(orderId)

	return args.Error(0)
}
