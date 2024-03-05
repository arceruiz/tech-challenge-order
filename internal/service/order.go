package service

import (
	"context"
	"fmt"
	"tech-challenge-order/internal/canonical"
	"tech-challenge-order/internal/integration/payment"
	"tech-challenge-order/internal/repository"
	"time"

	"github.com/shopspring/decimal"
)

type OrderService interface {
	GetAll(context.Context) ([]canonical.Order, error)
	Create(context.Context, canonical.Order) error
	Update(context.Context, string, canonical.Order) error
	GetByID(context.Context, string) (*canonical.Order, error)
	GetByStatus(context.Context, canonical.OrderStatus) ([]canonical.Order, error)
	CheckoutOrder(context.Context, string) (*canonical.Order, error)
}

type orderService struct {
	repo           repository.OrderRepository
	paymentService payment.PaymentService
}

func NewOrderService() OrderService {
	return &orderService{
		repo:           repository.NewOrderRepo(),
		paymentService: payment.NewPaymentService(),
	}
}

func (s *orderService) GetAll(ctx context.Context) ([]canonical.Order, error) {
	return s.repo.GetAll(ctx)
}

func (s *orderService) Create(ctx context.Context, order canonical.Order) error {
	order.ID = canonical.NewUUID()
	order.CreatedAt = time.Now()
	order.Status = canonical.ORDER_RECEIVED
	s.calculateTotal(&order)

	_, err := s.repo.Create(ctx, order)
	if err != nil {
		return fmt.Errorf("error creating order, %w", err)
	}
	return nil
}

func (s *orderService) Update(ctx context.Context, id string, updatedOrder canonical.Order) error {
	s.calculateTotal(&updatedOrder)
	return s.repo.Update(ctx, id, updatedOrder)
}

func (s *orderService) GetByID(ctx context.Context, id string) (*canonical.Order, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *orderService) GetByStatus(ctx context.Context, status canonical.OrderStatus) ([]canonical.Order, error) {
	return s.repo.GetByStatus(ctx, int(status))
}

func (s *orderService) CheckoutOrder(ctx context.Context, orderID string) (*canonical.Order, error) {
	order, err := s.repo.GetByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("payment not criated, error searching order, %w", err)
	}

	err = s.paymentService.Create(payment.Payment{
		PaymentType: 0,
		OrderID:     orderID,
	})
	if err != nil {
		return nil, fmt.Errorf("error checking out order, %w", err)
	}

	order.Status = canonical.ORDER_CHECKED_OUT
	order.UpdatedAt = time.Now()
	err = s.repo.Update(ctx, orderID, *order)
	if err != nil {
		return nil, fmt.Errorf("payment not criated, error updating order, %w", err)
	}

	return order, nil
}

func (s *orderService) calculateTotal(order *canonical.Order) {
	for _, product := range order.OrderItems {
		price := decimal.NewFromFloat(product.Price)
		quantity := decimal.NewFromInt(product.Quantity)
		productTotal, _ := price.Mul(quantity).Float64()

		order.Total += productTotal
	}
}
