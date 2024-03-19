package service

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

func (m *OrderRepositoryMock) Create(ctx context.Context, order canonical.Order) (*canonical.Order, error) {
	args := m.Called(order)
	return args.Get(0).(*canonical.Order), args.Error(1)
}

func (m *OrderRepositoryMock) Update(ctx context.Context, id string, Order canonical.Order) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *OrderRepositoryMock) UpdateStatus(_ context.Context, id string, status canonical.OrderStatus) error {
	args := m.Called(id)

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

type PublisherMock struct {
	mock.Mock
}

func (m *PublisherMock) SendMessage(inputMsg any, queueURL string) error {
	args := m.Called()

	return args.Error(0)
}

type ProductMock struct {
	mock.Mock
}

func (m *ProductMock) GetProducts(ctx context.Context, ids map[string]*canonical.OrderItem) error {
	args := m.Called(ids)

	return args.Error(0)
}

func (m *ProductMock) MockGetProducts(input string, f func(args mock.Arguments)) {
	m.On("GetProducts", mock.MatchedBy(func(i map[string]*canonical.OrderItem) bool {
		_, ok := i[input]
		return ok
	})).Return(nil).Run(f)
}
