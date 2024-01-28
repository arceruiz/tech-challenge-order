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
	ID         string      `bson:"_id"`
	CustomerID string      `bson:"customer_id"`
	Status     OrderStatus `bson:"status"`
	CreatedAt  time.Time   `bson:"created_at"`
	UpdatedAt  time.Time   `bson:"updated_at"`
	Total      float64     `bson:"total"`
	OrderItems []OrderItem `bson:"order_items"`
}

type OrderItem struct {
	Product
	Quantity int64 `bson:"quantity"`
}

type OrderStatus int

const (
	ORDER_RECEIVED    OrderStatus = iota // Order created, must be checked out
	ORDER_CHECKED_OUT                    // Order checked out, must be payed
	ORDER_PAYED                          // Order payed, must be prepared
	ORDER_PREPARING                      // Order being prepared, must be ready
	ORDER_READY                          // Order ready, must be delivered
	ORDER_DELIEVERED                     // Order delivered, must be concluded
	ORDER_COMPLETED                      // Order concluded
	ORDER_CANCELLED                      // Order cancelled
)

var MapOrderStatus = map[string]OrderStatus{ //ajustar chaves
	"RECEIVED":    ORDER_RECEIVED,
	"CHECKED_OUT": ORDER_CHECKED_OUT,
	"PAYED":       ORDER_PAYED,
	"PREPARING":   ORDER_PREPARING,
	"READY":       ORDER_READY,
	"DELIEVERE":   ORDER_DELIEVERED,
	"COMPLETED":   ORDER_COMPLETED,
	"CANCELLED":   ORDER_CANCELLED,
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
