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
	OrderItems []OrderItem `json:"products,omitempty"`
}

type OrderResponse struct {
	ID         string              `json:"id,omitempty"`
	CustomerID string              `json:"customer_id,omitempty"`
	Status     string              `json:"status,omitempty"`
	CreatedAt  time.Time           `json:"created_at,omitempty"`
	UpdatedAt  time.Time           `json:"updated_at,omitempty"`
	Products   []OrderItemResponse `json:"products,omitempty"`
	Total      float64             `json:"total,omitempty"`
}

type OrderItem struct {
	ProductId string `json:"product_id"`
	Quantity  int64  `json:"quantity"`
}

type OrderItemResponse struct {
	ProductItem
	Quantity int64 `json:"quantity"`
}
