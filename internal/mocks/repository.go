package mocks

import (
	"context"
	"tech-challenge-order/internal/canonical"

	"github.com/stretchr/testify/mock"
)

type OrderRepositoryMock struct {
	mock.Mock
}

func (m *OrderRepositoryMock) GetAll(ctx context.Context) ([]canonical.Order, error) {
	args := m.Called(ctx)
	return args.Get(0).([]canonical.Order), args.Error(1)
}

func (m *OrderRepositoryMock) Create(ctx context.Context, Order canonical.Order) (*canonical.Order, error) {
	args := m.Called(ctx, Order)
	return args.Get(0).(*canonical.Order), args.Error(1)
}

func (m *OrderRepositoryMock) Update(ctx context.Context, id string, Order canonical.Order) error {
	args := m.Called(ctx, id, Order)
	return args.Error(0)
}

func (m *OrderRepositoryMock) GetByID(ctx context.Context, id string) (*canonical.Order, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*canonical.Order), args.Error(1)
}

func (m *OrderRepositoryMock) GetByStatus(ctx context.Context, status int) ([]canonical.Order, error) {
	args := m.Called(ctx, status)
	return args.Get(0).([]canonical.Order), args.Error(1)
}
