package service

import (
	"context"
	"errors"
	"fmt"
	"tech-challenge-order/internal/canonical"
	"tech-challenge-order/internal/config"
	"tech-challenge-order/internal/integration/product"
	"tech-challenge-order/internal/integration/sqs_publisher"
	"tech-challenge-order/internal/repository"
	"time"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

type OrderService interface {
	GetAll(context.Context) ([]canonical.Order, error)
	Create(context.Context, canonical.Order) error
	Update(context.Context, string, canonical.Order) error
	GetByID(context.Context, string) (*canonical.Order, error)
	GetByStatus(context.Context, canonical.OrderStatus) ([]canonical.Order, error)
	CheckoutOrder(context.Context, string) (*canonical.Order, error)
	UpdateStatus(ctx context.Context, orderId string, status canonical.OrderStatus) error
}

type orderService struct {
	repo                       repository.OrderRepository
	productService             product.ProductService
	publisher                  sqs_publisher.Publisher
	orderQueueAddress          string
	paymentPendingQueueAddress string
}

func NewOrderService() OrderService {
	return &orderService{
		repo:                       repository.NewOrderRepo(),
		productService:             product.NewProduct(),
		publisher:                  sqs_publisher.NewSQS(),
		orderQueueAddress:          config.Get().SQS.OrderQueue,
		paymentPendingQueueAddress: config.Get().SQS.PaymentPendingQueue,
	}
}

func (s *orderService) GetAll(ctx context.Context) ([]canonical.Order, error) {
	return s.repo.GetAll(ctx)
}

func (s *orderService) Create(ctx context.Context, order canonical.Order) error {
	order.ID = canonical.NewUUID()
	order.Status = canonical.ORDER_RECEIVED
	order.CreatedAt = time.Now()

	err := s.productService.GetProducts(ctx, order.OrderItems)
	if err != nil {
		return err
	}

	s.calculateTotal(&order)

	_, err = s.repo.Create(context.Background(), order)
	if err != nil {
		return err
	}

	if err := s.publisher.SendMessage(order.ID, s.orderQueueAddress); err != nil {
		logrus.WithError(err).WithField("order_id", order.ID)
		return fmt.Errorf("an error occurred when creating order")
	}

	return nil
}

func (s *orderService) Update(ctx context.Context, id string, updatedOrder canonical.Order) error {
	s.calculateTotal(&updatedOrder)
	return s.repo.Update(ctx, id, updatedOrder)
}

func (s *orderService) UpdateStatus(ctx context.Context, orderId string, status canonical.OrderStatus) error {
	order, err := s.repo.GetByID(ctx, orderId)
	if err != nil {
		return err
	}

	if order == nil {
		return errors.New("order not found")
	}

	return s.repo.UpdateStatus(ctx, orderId, status)
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

	err = s.publisher.SendMessage(order.ID, s.paymentPendingQueueAddress)
	if err != nil {
		return nil, fmt.Errorf("error checking out order, %w", err)
	}

	order.UpdatedAt = time.Now()

	err = s.repo.UpdateStatus(ctx, orderID, canonical.ORDER_PAYMENT_PENDING)
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
