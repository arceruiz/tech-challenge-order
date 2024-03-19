package canonical

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	ErrorNotFound = fmt.Errorf("entity not found")
)

type Product struct {
	ID       string  `bson:"product_id"`
	Name     string  `bson:"name"`
	Price    float64 `bson:"price"`
	Category string  `bson:"category"`
}

type Order struct {
	ID         string                `bson:"_id"`
	CustomerID string                `bson:"customer_id"`
	Status     OrderStatus           `bson:"status"`
	CreatedAt  time.Time             `bson:"created_at"`
	UpdatedAt  time.Time             `bson:"updated_at"`
	Total      float64               `bson:"total"`
	OrderItems map[string]*OrderItem `bson:"order_items"`
}

type OrderItem struct {
	Product
	Quantity int64 `bson:"quantity"`
}

type OrderStatus int

const (
	ORDER_RECEIVED        = 0
	ORDER_PAYMENT_PENDING = 1
	ORDER_PAYED           = 2
	ORDER_PREPARING       = 3
	ORDER_COMPLETED       = 4
	ORDER_CANCELLED       = 5
)

var MapOrderStatus = map[string]OrderStatus{
	"RECEIVED":        ORDER_RECEIVED,
	"PAYMENT_PENDING": ORDER_PAYMENT_PENDING,
	"PAYED":           ORDER_PAYED,
	"PREPARING":       ORDER_PREPARING,
	"COMPLETED":       ORDER_COMPLETED,
	"CANCELLED":       ORDER_CANCELLED,
}

func HandleError(err error) error {
	if errors.Is(err, ErrorNotFound) {
		return err
	}
	return fmt.Errorf("unexpected error occurred %w", err)

}

func NewUUID() string {
	return uuid.New().String()
}
