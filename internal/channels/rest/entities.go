package rest

import "time"

type Response struct {
	Message string `json:"message"`
}

type ProductItem struct {
	ID       string  `json:"id,omitempty"`
	Name     string  `json:"name,omitempty"`
	Price    float64 `json:"price,omitempty"`
	Category string  `json:"category,omitempty"`
}

type OrderRequest struct {
	ID         string
	CustomerID string      `json:"customer_id,omitempty"`
	Status     string      `json:"status,omitempty"`
	CreatedAt  time.Time   `json:"created_at,omitempty"`
	UpdatedAt  time.Time   `json:"updated_at,omitempty"`
	Total      float64     `json:"total,omitempty"`
	OrderItems []OrderItem `json:"products,omitempty"`
}

type OrderResponse struct {
	ID         string      `json:"id,omitempty"`
	CustomerID string      `json:"customer_id,omitempty"`
	Status     string      `json:"status,omitempty"`
	CreatedAt  time.Time   `json:"created_at,omitempty"`
	UpdatedAt  time.Time   `json:"updated_at,omitempty"`
	Products   []OrderItem `json:"products,omitempty"`
}

type OrderItem struct {
	ProductItem
	Quantity int64 `json:"quantity"`
}

type PaymentRest struct {
	ID          int       `json:"id"`
	PaymentType int       `json:"payment_type"`
	CreatedAt   time.Time `json:"created_at"`
	Status      int       `json:"status"`
}

type PaymentCallback struct {
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
}
